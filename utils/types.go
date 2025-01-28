package utils

import (
	"fmt"
	"strings"
	"time"
)

// type constraints for generic functions

type Validator interface {
	Validate() error
}

type Doctor struct {
	Id           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type DoctorRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d DoctorRequest) Validate() error {
	if d.Email == "" || d.Password == "" {
		return fmt.Errorf("email and password are required")
	}
	if !strings.Contains(d.Email, "@") {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

func NewDoctorReques(email, password string) *DoctorRequest {
	return &DoctorRequest{Email: email, Password: password}
}

type DoctorInfo struct {
	Id                    int    `json:"id"`
	Age                   int    `json:"age"`
	DoctorId              int    `json:"doctor_id"`
	Speciality            string `json:"speciality"`
	City                  string `json:"city"`
	Phone                 string `json:"phone"`
	ImgUrl                string `json:"img_url"`
	ProfessionalStatement string `json:"professional_statement"`
	Experience            string `json:"experience"`
	SatisfactionScore     int    `json:"satisfaction_score"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Governorate           string `json:"governorate"`
}
type DoctorInfoRequest struct {
	Age                   int    `json:"age"`
	DoctorId              int    `json:"doctor_id"`
	Speciality            string `json:"speciality"`
	City                  string `json:"city"`
	Phone                 string `json:"phone"`
	ImgUrl                string `json:"img_url"`                // *
	ProfessionalStatement string `json:"professional_statement"` // *
	Experience            string `json:"experience"`             // *
	SatisfactionScore     int    `json:"satisfaction_score"`     // *
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Governorate           string `json:"governorate"`
}

func (d DoctorInfoRequest) Validate() error {
	if d.Age == 0 || d.DoctorId == 0 || d.Speciality == "" || d.City == "" || d.Phone == "" || d.FirstName == "" || d.LastName == "" || d.Governorate == "" {
		return fmt.Errorf("age, doctor_id, speciality, city, phone, first_name, last_name, governorate are required")
	}
	return nil

}

type DoctorClinicInfo struct {
	Id                  int       `json:"id"`
	Governorate         string    `json:"governorate"`
	City                string    `json:"city"`
	Address             string    `json:"address"`
	Latitude            string    `json:"latitude"`
	Longitude           string    `json:"longitude"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	WorkingHours        int       `json:"working_hours"`
	Holiday             []int     `json:"holiday"`
	Location            string    `json:"location"`
	PricePerAppointment int       `json:"price_per_appointment"`
	PatientNumberPerDay int       `json:"patient_number_per_day"`
	DoctorId            int       `json:"doctor_id"`
	AppointmentDuration int       `json:"appointment_duration"`
}

type DoctorClinicInfoRequest struct {
	Governorate  string    `json:"governorate"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	Latitude     string    `json:"latitude"`
	Longitude    string    `json:"longitude"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	WorkingHours int       `json:"working_hours"`
	Holiday      []int     `json:"holiday"`

	PricePerAppointment int `json:"price_per_appointment"`  // *
	PatientNumberPerDay int `json:"patient_number_per_day"` // *
	DoctorId            int `json:"doctor_id"`
	AppointmentDuration int `json:"appointment_duration"`
}

func (d DoctorClinicInfoRequest) Validate() error {
	if d.Governorate == "" || d.City == "" || d.Address == "" || d.Latitude == "" || d.Longitude == "" || d.StartTime.IsZero() || d.EndTime.IsZero() || d.WorkingHours == 0 || d.DoctorId == 0 || d.AppointmentDuration == 0 {
		return fmt.Errorf("governorate, city, address, latitude, longitude, start_time, end_time, working_hours, location, doctor_id, appointment_duration are required")
	}
	return nil
}

// status
type ScheduleStatus string

const (
	Booked    ScheduleStatus = "booked"
	Available ScheduleStatus = "available"
)

// this represent a row in doctor_schedules table
// each row represent a time interval on specific date (when doctor is available)
type Schedule struct {
	Id        int            `json:"id"`
	Date      time.Time      `json:"date"` // it will contain date and time when you want to present it to the user you have to format it into a string of this format "YYYY-MM-DD HH:MM:SS"
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	DoctorId  int            `json:"doctor_id"`
	PatientId int            `json:"patient_id"`
	Status    ScheduleStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// the whole point of this struct is when we want to create
// a new schedule on specific date and according to the start and end time
// i'll make this request accept sngle or multiple date in order to populate these date
// with specific time slots and according to the appointment duration
type ScheduleRequest struct {
	Date                []time.Time    `json:"date"` // a slice of date
	StartTime           time.Time      `json:"start_time"`
	EndTime             time.Time      `json:"end_time"`
	DoctorId            int            `json:"doctor_id"`
	AppointmentDuration *time.Duration `json:"appointment_duration"` // in minutes
}

func (s ScheduleRequest) Validate() error {
	if len(s.Date) == 0 || s.StartTime.IsZero() || s.EndTime.IsZero() || s.AppointmentDuration == nil {
		return fmt.Errorf("date, start_time, end_time , appointment_duration and doctor_id are required")
	}
	return nil
}

// to see the available slots for a specific doctor on specific date/date(s)
type GetAvailabilityRequest struct {
	DoctorId *int        `json:"doctor_id"`
	Date     []time.Time `json:"date"`
}

func (g GetAvailabilityRequest) Validate() error {
	if g.DoctorId == nil || len(g.Date) == 0 {
		return fmt.Errorf("doctor_id and date are required")
	}
	return nil
}

// modify availability of a doctor on specific date/date(s)
type ModifiyAvailabilityRequest struct {
	DoctorId          *int          `json:"doctor_id"`
	Date              []time.Time   `json:"date"`
	StartTime         time.Time     `json:"start_time"`
	EndTime           time.Time     `json:"end_time"`
	AppointmentLength time.Duration `json:"appointment_length"`
}

type BookSlotRequest struct {
	PatientId *int `json:"patient_id"`
	SlotId    *int `json:"slot_id"`
}

func (b BookSlotRequest) Validate() error {
	if b.PatientId == nil || b.SlotId == nil {
		return fmt.Errorf("patient_id and slot_id are required")
	}
	return nil
}

type CancelSlotRequest struct {
	SlotId *int `json:"slot_id"`
}

func (c CancelSlotRequest) Validate() error {
	if c.SlotId == nil {
		return fmt.Errorf("slot_id is required")
	}
	return nil
}
