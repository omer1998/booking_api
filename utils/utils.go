package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type ApiError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func NewError(message string, code int) *ApiError {
	return &ApiError{Error: message, Code: code}
}

type HandleFunction func(w http.ResponseWriter, r *http.Request) *ApiError

func MakeHandleFunc(fn HandleFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			MakeJsonResponse(w, err.Code, err)

		}

	}
}

func MakeJsonResponse(w http.ResponseWriter, code int, value any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	// err := json.NewEncoder(w).Encode(map[string]any{"message": value})
	json.NewEncoder(w).Encode(value)

}

// another approach to deal with validationa and encoding and decoding of the request and the response
func Encode[T any](v T, status int, w http.ResponseWriter) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to encode response: %v", err)
	}
	return nil
}

func Decode[T Validator](r *http.Request) (*T, error) {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %v", err)
	}
	err := t.Validate()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func MakeJsonErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Add("Content-Type", "application/json")
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type Model interface {
	DoctorRequest
}

func GetModelFromRequest[M Model](w http.ResponseWriter, r *http.Request) *M {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"

			MakeJsonResponse(w, http.StatusUnsupportedMediaType, NewError(msg, http.StatusUnsupportedMediaType))

		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p M
	err := dec.Decode(&p)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			MakeJsonErrorResponse(w, http.StatusBadRequest, msg)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintln("Request body contains badly-formed JSON")
			MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
			// http.Error(w, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
			// http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
			// http.Error(w, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
			// http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			MakeJsonErrorResponse(w, http.StatusRequestEntityTooLarge, msg)
			// http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Print(err.Error())
			MakeJsonErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

	}
	return &p
}

func HandleDecodeError(w http.ResponseWriter, r *http.Request, err error) {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"

			MakeJsonResponse(w, http.StatusUnsupportedMediaType, NewError(msg, http.StatusUnsupportedMediaType))
			return
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	// Catch any syntax errors in the JSON and send an error message
	// which interpolates the location of the problem to make it
	// easier for the client to fix.
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
		return

	// In some circumstances Decode() may also return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
	// is an open issue regarding this at
	// https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintln("Request body contains badly-formed JSON")
		MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
		return
		// http.Error(w, msg, http.StatusBadRequest)

	// Catch any type errors, like trying to assign a string in the
	// JSON request body to a int field in our Person struct. We can
	// interpolate the relevant field name and position into the error
	// message to make it easier for the client to fix.
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
		return
		// http.Error(w, msg, http.StatusBadRequest)

	// Catch the error caused by extra unexpected fields in the request
	// body. We extract the field name from the error message and
	// interpolate it in our custom error message. There is an open
	// issue at https://github.com/golang/go/issues/29035 regarding
	// turning this into a sentinel error.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
		return
		// http.Error(w, msg, http.StatusBadRequest)

	// An io.EOF error is returned by Decode() if the request body is
	// empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		MakeJsonErrorResponse(w, http.StatusBadRequest, msg)
		return
		// http.Error(w, msg, http.StatusBadRequest)

	// Catch the error caused by the request body being too large. Again
	// there is an open issue regarding turning this into a sentinel
	// error at https://github.com/golang/go/issues/30715.
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		MakeJsonErrorResponse(w, http.StatusRequestEntityTooLarge, msg)
		return
		// http.Error(w, msg, http.StatusRequestEntityTooLarge)

	// Otherwise default to logging the error and sending a 500 Internal
	// Server Error response.
	default:
		log.Print(err.Error())
		MakeJsonErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func HandlePgxError(err error) string {
	var pgError *pgconn.PgError
	result := errors.As(err, &pgError)
	if !result {
		panic("error is pgx error handler")
	}
	return pgError.Message
}

// not understandable
// type MR interface {
// 	Scan(dest ...interface{}) error
// }

//	func ScanIntoDoctor[M MR](rows M) (*Doctor, error) {
//		var doctor = new(Doctor)
//		err := rows.Scan(&doctor.Id, &doctor.Email, &doctor.Password, &doctor.RegisteredAt)
//		if err != nil {
//			return nil, err
//		}
//		return doctor, nil
//	}
func ScanIntoDoctor(rows pgx.Rows) (*Doctor, error) {
	var doctor = new(Doctor)
	err := rows.Scan(&doctor.Id, &doctor.Email, &doctor.Password, &doctor.RegisteredAt)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}
func ScanIntoDoctorSingle(rows pgx.Row) (*Doctor, error) {
	var doctor = new(Doctor)
	err := rows.Scan(&doctor.Id, &doctor.Email, &doctor.Password, &doctor.RegisteredAt)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

// this functionality is also known as middleware where the handler here is the next handler in the chain
func CheckAuthState(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtString := strings.TrimSpace(r.Header.Get("Authorization"))
		if jwtString == "" {
			MakeJsonResponse(w, http.StatusUnauthorized, &ApiError{Error: "Unauthorized", Code: http.StatusUnauthorized})
			return
		} else {
			token, err := ParseJWTToken(jwtString)
			if err != nil {
				MakeJsonResponse(w, http.StatusUnauthorized, &ApiError{Error: fmt.Sprintf("Unauthorized: %s", err.Error()), Code: http.StatusUnauthorized})
				return
			}
			if !token.Valid {
				MakeJsonResponse(w, http.StatusUnauthorized, &ApiError{Error: "Unauthorized", Code: http.StatusUnauthorized})
				return
			}
			claim, _ := token.Claims.(MyClaims) // assertion process here
			r.Header.Add("id", claim.Id)
			r.Header.Add("email", claim.Email)
			// you can use these values in the next handler
			// or may be we can create a context and pass these values to the next handler by using context.WithValue
			// ctx := context.WithValue(r.Context(), "id", claim.Id)
			// ctx = context.WithValue(ctx, "email", claim.Email)
			// r = r.WithContext(ctx)

			handler(w, r)
		}

	}

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
