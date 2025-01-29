package database

import (
	"time"

	"github.com/omer1998/booking_api/utils"
)

type Database interface {
	AddDoctor(email, password string) *utils.ApiError
	GetDoctors() ([]utils.Doctor, *utils.ApiError)
	GetDoctorByEmail(email string) (*utils.Doctor, *utils.ApiError)
	DoctorDetail
	ScheduleManagement
}

type DoctorDetail interface {
	AddDoctorInfo(doctorInfo utils.DoctorInfoRequest) *utils.ApiError
	AddClinicInfo(doctorClinic utils.DoctorClinicInfoRequest) *utils.ApiError
	AddDoctorDetail(doctorInfo utils.DoctorInfoRequest, doctorClinic utils.DoctorClinicInfoRequest) *utils.ApiError // this is implemented by transaction
	GetDoctorInfo(id int) (*utils.DoctorInfo, *utils.ApiError)
	GetClinicInfo(docotorId int) (*utils.DoctorClinicInfo, *utils.ApiError)
	UpdateDoctorInfo(doctorInfo utils.DoctorInfoRequest) (*utils.DoctorInfo, *utils.ApiError)
	UpdateClinicInfo(doctorInfo utils.DoctorClinicInfoRequest) (*utils.DoctorClinicInfo, *utils.ApiError)
}

type ScheduleManagement interface {
	CreateSchedule(date []time.Time, doctorId int, startTime time.Time, endTime time.Time, appointmentDuration time.Duration) *utils.ApiError
	//GetSchedule(doctorId int, date time.Time) (*utils.Schedule, *utils.ApiError)
	//GetSchedules(doctorId int, date []time.Time) ([]*utils.Schedule, *utils.ApiError)
	GetAvailableTimeSlot(doctorId int, date []time.Time) ([]utils.Schedule, *utils.ApiError)
	// here we need to complete patient registration
	BookSlot(patientId int, scheduleId int) *utils.ApiError
	CancelSlot(scheduleId int) *utils.ApiError
	GetDoctorSchedules(doctorId int, startDate time.Time, endDate *time.Time) ([]utils.Schedule, *utils.ApiError)
}

// func RunDb(cxt context.Context) *pgx.Conn {
// 	conn, err := pgx.Connect(cxt, "postgres://postgres:postgres@localhost:5432/booking")
// 	if err != nil {
// 		fmt.Println("Unable to connect to database: ", err)
// 		panic(err)
// 	}
// 	fmt.Println("Connected to booking BD at port 5432")
// 	return conn
// }
