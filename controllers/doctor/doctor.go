package doctor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/omer1998/booking_api/database"
	"github.com/omer1998/booking_api/services"
	"github.com/omer1998/booking_api/utils"
)

type DoctorApi struct {
	Db   database.Database
	Auth *services.AuthenticationService
}

func NewDoctorApi(db database.Database, auth *services.AuthenticationService) *DoctorApi {
	return &DoctorApi{
		Db:   db,
		Auth: auth,
	}
}

// func (d *DoctorApi) HandleGetDoctorByEmail(w http.ResponseWriter, r *http.Request) *utils.ApiError {
// 	email := mux.Vars(r)["email"]

// 	doctor := utils.NewDoctor(email, "12345678")
// 	err := utils.MakeJsonResponse(w, http.StatusOK, doctor)
// 	if err != nil {
// 		fmt.Println(err)
// 		return &utils.ApiError{Error: "Internal Server Error", Code: http.StatusInternalServerError}
// 	}
// 	return nil

// }

func (d *DoctorApi) HandleAddDoctor() http.HandlerFunc {
	return utils.MakeHandleFunc(func(w http.ResponseWriter, r *http.Request) *utils.ApiError {
		req := new(utils.DoctorRequest)
		err := json.NewDecoder(r.Body).Decode(req)
		fmt.Println("doctorReqData", req)

		if err != nil {
			fmt.Println("error in decoding the request body", err)
			return &utils.ApiError{Error: err.Error(), Code: http.StatusBadRequest}
		}
		// utils.GetModelFromRequest[utils.DoctorRequest](w, r,)
		doctorRequest := utils.NewDoctorReques(req.Email, req.Password)
		println("doctorRequest", doctorRequest.Email, doctorRequest.Password)
		if d.Db == nil {
			fmt.Println("database connection is nil")
			return &utils.ApiError{Error: "database connection is nil", Code: http.StatusInternalServerError}
		}
		if err := d.Db.AddDoctor(doctorRequest.Email, doctorRequest.Password); err != nil {
			return err
		}
		token, err := utils.CreateJWTToken(doctorRequest.Email, "10", 24)
		if err != nil {
			return &utils.ApiError{Error: err.Error(), Code: http.StatusInternalServerError}
		}
		utils.MakeJsonResponse(w, http.StatusOK, map[string]string{"message": "Doctor added successfully", "token": *token})
		return nil
	})
	// utils.HandleRequestError[utils.DoctorRequest](w, r)

	// doctorReqData := utils.GetModelFromRequest[utils.DoctorRequest](w, r)

}

func (d *DoctorApi) HandleGetAllDoctors(w http.ResponseWriter, r *http.Request) *utils.ApiError {
	doctors, err := d.Db.GetDoctors()
	if err != nil {
		return err
	}
	if len(doctors) == 0 {
		utils.MakeJsonResponse(w, http.StatusOK, map[string]string{"message": "No doctors found"})
		return nil
	}
	utils.MakeJsonResponse(w, http.StatusOK, doctors)
	return nil
}

func (d *DoctorApi) HandleGetDoctorByEmail(w http.ResponseWriter, r *http.Request) *utils.ApiError {
	email := strings.TrimSpace(mux.Vars(r)["email"])
	if email == "" {
		return &utils.ApiError{Error: "email is required", Code: http.StatusBadRequest}
	}
	if !strings.Contains(email, "@") {
		return &utils.ApiError{Error: "email is invalid", Code: http.StatusBadRequest}
	}
	// fmt.Println("email", email)
	doctor, err := d.Db.GetDoctorByEmail(email)
	if err != nil {
		return err
	}
	utils.MakeJsonResponse(w, http.StatusOK, doctor)

	return nil
}

func (d *DoctorApi) HandleDoctorLogin(w http.ResponseWriter, r *http.Request) *utils.ApiError {
	doctor, err := utils.Decode[utils.DoctorRequest](r)
	if err != nil {
		return utils.NewError(err.Error(), http.StatusBadRequest)

	}
	doc, logErr := d.Auth.LoginDoctor(doctor.Email, doctor.Password)
	if logErr != nil {
		return logErr
	}
	token, tokenErr := utils.CreateJWTToken(doc.Email, strconv.Itoa(doc.Id), 40)
	if tokenErr != nil {
		return utils.NewError(tokenErr.Error(), http.StatusInternalServerError)
	}
	response := map[string]any{
		"doctor": *doc,
		"token":  *token,
	}
	utils.MakeJsonResponse(w, http.StatusOK, response)
	return nil

}

func (d *DoctorApi) HandleDoctorRegister() http.HandlerFunc {
	return utils.MakeHandleFunc(func(w http.ResponseWriter, r *http.Request) *utils.ApiError {

		doctor, err := utils.Decode[utils.DoctorRequest](r)
		if err != nil {
			return utils.NewError(err.Error(), http.StatusBadRequest)
		}
		regErr := d.Auth.RegisterDoctor(doctor.Email, doctor.Password)
		if regErr != nil {
			return regErr
		}

		utils.MakeJsonResponse(w, http.StatusOK, map[string]string{"message": "Doctor registered successfully", "email": doctor.Email})
		return nil
	})

}

// TODO:: Implement the following functions
func (d *DoctorApi) HandleDoctorInfo() http.HandlerFunc {
	return utils.MakeHandleFunc(func(w http.ResponseWriter, r *http.Request) *utils.ApiError {
		if r.Method == "GET" {
			doctorId := r.URL.Query().Get("doctor_id")
			slog.Info("doctorId", doctorId, doctorId)
		} else if r.Method == "PUT" || r.Method == "POST" {
			docInfo, err := utils.Decode[utils.DoctorInfoRequest](r)
			if err != nil {
				return utils.NewError(err.Error(), http.StatusBadRequest)
			}
			if err := d.Db.AddDoctorInfo(*docInfo); err != nil {
				return err
			}
			return nil
		} else {
			return utils.NewError("Method not allowed", http.StatusMethodNotAllowed)
		}

		return nil

	})

}
func (d *DoctorApi) HandleClinicInfo(w http.ResponseWriter, r *http.Request) *utils.ApiError {
	if r.Method == "GET" {
		return nil
	} else if r.Method == "PUT" || r.Method == "POST" {
		// Handle PUT request
		return nil
	} else {
		return utils.NewError("Method not allowed", http.StatusMethodNotAllowed)
	}

}
