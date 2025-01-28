package services

import (
	"fmt"
	"net/http"

	"github.com/omer1998/booking_api/database"
	"github.com/omer1998/booking_api/utils"
)

type AuthenticationService struct {
	db database.Database
}

func NewAuthenticationService(db database.Database) *AuthenticationService {
	return &AuthenticationService{db: db}
}

func (auth *AuthenticationService) LoginDoctor(email, password string) (*utils.Doctor, *utils.ApiError) {
	// .GetConnPool().Query(auth.db.cxt, "select id, email,password, register_date from doctors where email=$1", email)
	doc, err := auth.db.GetDoctorByEmail(email)
	if err != nil {
		return nil, err
	}
	if doc != nil {
		verifyPassword := utils.VerifyPassword(password, doc.Password)
		if verifyPassword {
			return doc, nil
		}
		return nil, utils.NewError("Invalid password", http.StatusUnauthorized)
	}
	return nil, &utils.ApiError{Error: "Not found, Register first", Code: http.StatusNotFound}

}

func (auth *AuthenticationService) RegisterDoctor(email, password string) *utils.ApiError {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error hashing password: %v", err.Error()), http.StatusInternalServerError)
	}
	dbErr := auth.db.AddDoctor(email, hashedPassword)
	if dbErr != nil {
		// var pgErr *pgconn.PgError
		// errors.As(err, &pgErr)
		// fmt.Println("error message is", pgErr.Message)
		// fmt.Println("error code is", pgErr.Code)
		// fmt.Println("error detail is", pgErr.Detail)
		// fmt.Println("error message is", pgErr.Hint)

		return dbErr
	}
	return nil
}
