package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omer1998/booking_api/config"
	"github.com/omer1998/booking_api/utils"
)

type PostgresDb struct {
	// conn     *pgx.Conn
	connPool *pgxpool.Pool
	cxt      context.Context
}

//	func NewPostgresDb(conn *pgx.Conn) *PostgresDb {
//		return &PostgresDb{
//			conn: conn,
//		}
//	}
func NewPostgresDbPool(connPool *pgxpool.Pool, cxt context.Context) *PostgresDb {
	return &PostgresDb{
		connPool: connPool,
		cxt:      cxt,
	}
}

func Connect(cxt context.Context) *pgx.Conn {
	conn, err := pgx.Connect(cxt, "postgres://postgres:12345678@localhost:5432/booking")
	if err != nil {
		fmt.Println("Unable to connect to database: ", err.Error())
		panic(err)
	}
	fmt.Println("Connected to booking BD at port 5432")
	return conn
}
func ConnectPool(cxt context.Context, config *config.Config) *pgxpool.Pool {
	// here we can set the pool configuration
	// myPoolConfig := pgxpool.Config{
	// 	number
	// }
	// localDbUrl:= "postgres://postgres:12345678@localhost:5432/booking"

	// dockerDbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	connectionDbUrl := config.GetConnectionDbUrl()
	connPool, err := pgxpool.New(cxt, connectionDbUrl)
	if err != nil {
		fmt.Println("Unable to create pool to database: ", err.Error())
		panic(err)
	}
	if connPool == nil {
		panic(errors.New("pool is nil"))
	}
	fmt.Println("connection pool is created to booking BD ")
	pingErr := connPool.Ping(cxt)
	if pingErr != nil {
		fmt.Println("Unable to ping to database: ", pingErr.Error())
		panic(pingErr)

	}
	fmt.Println("ping to booking BD is successful")
	return connPool

}

func (db *PostgresDb) AddDoctor(email, password string) *utils.ApiError {
	_, err := db.connPool.Exec(context.Background(), "INSERT INTO doctors (email, password) VALUES ($1, $2)", email, password)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		// fmt.Println("error message is", pgErr.Message)
		// fmt.Println("error code is", pgErr.Code)
		// fmt.Println("error detail is", pgErr.Detail)
		// fmt.Println("error message is", pgErr.Hint)

		return &utils.ApiError{Error: utils.HandlePgxError(err), Code: http.StatusInternalServerError}
	}
	return nil
}

func (db *PostgresDb) GetDoctors() ([]utils.Doctor, *utils.ApiError) {
	rows, err := db.connPool.Query(context.Background(), "SELECT id, email, password, register_date FROM doctors")
	if err != nil {
		return nil, &utils.ApiError{Error: err.Error(), Code: http.StatusInternalServerError}
	}
	var doctors []utils.Doctor = []utils.Doctor{}
	// var doctor utils.Doctor
	for rows.Next() {
		// err := rows.Scan(&doctor.Id, &doctor.Email, &doctor.Password, &doctor.RegisteredAt)
		doctor, err := utils.ScanIntoDoctor(rows)
		if err != nil {
			return nil, &utils.ApiError{Error: err.Error(), Code: http.StatusInternalServerError}

		}
		doctors = append(doctors, *doctor)
	}
	println("doctors", doctors)

	return doctors, nil
}

func (db *PostgresDb) GetDoctorByEmail(email string) (*utils.Doctor, *utils.ApiError) {

	rows, err := db.connPool.Query(db.cxt, "select id, email,password, register_date from doctors where email=$1", email)
	if err != nil {
		return nil, &utils.ApiError{Error: err.Error(), Code: http.StatusInternalServerError}

	}

	defer rows.Close()
	for rows.Next() {
		doctor, err := utils.ScanIntoDoctor(rows)
		if err != nil {
			return nil, &utils.ApiError{Error: err.Error(), Code: http.StatusInternalServerError}
		}
		if doctor != nil {
			return doctor, nil
		}
	}
	return nil, &utils.ApiError{Error: "Doctor not found", Code: http.StatusNotFound}

	// doctor, err := utils.ScanIntoDoctorSingle(row)
	// if err != nil {
	// 	return nil, &utils.ApiError{Error: err.Error(), Code: http.StatusInternalServerError}
	// }
	// if doctor != nil {
	// 	return doctor, nil
	// }
	// return nil, &utils.ApiError{Error: "Doctor not found", Code: http.StatusNotFound}
}

func (db *PostgresDb) AddDoctorInfo(doctorInfo utils.DoctorInfoRequest) *utils.ApiError {

	_, err := db.connPool.Exec(db.cxt, `insert into doctors_info 
	(age, doctor_id, speciality, city, phone, img_url, professional_statement, experience, satisfaction_score, first_name, last_name, governorate) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		doctorInfo.Age, doctorInfo.DoctorId, doctorInfo.Speciality, doctorInfo.City, doctorInfo.Phone, doctorInfo.ImgUrl, doctorInfo.ProfessionalStatement, doctorInfo.Experience, doctorInfo.SatisfactionScore, doctorInfo.FirstName, doctorInfo.LastName, doctorInfo.Governorate)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error adding doctor info: %v", err.Error()), http.StatusInternalServerError)
	}

	return nil
}
func (db *PostgresDb) AddClinicInfo(doctorClinic utils.DoctorClinicInfoRequest) *utils.ApiError {
	rows, err := db.connPool.Query(db.cxt, `insert into doctors_clinic
	(governorate, city, address, latitude, longitude, start_time, end_time, working_hours, holiday, 
	 price_per_appointment, patient_number_per_day, doctor_id, appointment_duration ) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		doctorClinic.Governorate, doctorClinic.City, doctorClinic.Address, doctorClinic.Latitude, doctorClinic.Longitude, doctorClinic.StartTime, doctorClinic.EndTime,
		doctorClinic.WorkingHours, doctorClinic.Holiday, doctorClinic.PricePerAppointment, doctorClinic.PatientNumberPerDay, doctorClinic.DoctorId, doctorClinic.AppointmentDuration)

	if err != nil {
		return utils.NewError(fmt.Sprintf("Error adding clinic info: %v", utils.HandlePgxError(err)), http.StatusInternalServerError)
	}
	if rows.Err() != nil {
		return utils.NewError(fmt.Sprintf("Error adding clinic info: %v", utils.HandlePgxError(rows.Err())), http.StatusInternalServerError)
	}
	defer rows.Close()
	return nil
}

func (db *PostgresDb) GetDoctorInfo(doctorId int) (*utils.DoctorInfo, *utils.ApiError) {
	row := db.connPool.QueryRow(db.cxt, `select * from doctors_info where doctor_id = $1`, doctorId)
	// id, age, doctor_id, speciality, city, phone, img_url, professional_statement, experience, satisfaction_score, first_name, last_name, governorate
	var docInfo = new(utils.DoctorInfo)
	// if err != nil {
	// 	return nil, utils.NewError(fmt.Sprintf("Error getting doctor info: %v", utils.HandlePgxError(err)), http.StatusInternalServerError)
	// }
	// defer rows.Close()
	// doctorDetail := new(utils.DoctorInfo)

	// if rows.Next() {
	// 	scanErr := rows.Scan(&doctorDetail.Age, &doctorDetail.DoctorId, &doctorDetail.Speciality, &doctorDetail.City, &doctorDetail.Phone, &doctorDetail.ImgUrl, &doctorDetail.ProfessionalStatement, &doctorDetail.Experience, &doctorDetail.SatisfactionScore, &doctorDetail.FirstName, &doctorDetail.LastName, &doctorDetail.Governorate)

	// 	if scanErr != nil {
	// 		return nil, utils.NewError(fmt.Sprintf("Error scanning doctor info: %v", utils.HandlePgxError(scanErr)), http.StatusInternalServerError)
	// 	}
	// 	return doctorDetail, nil
	// } else {
	// 	return nil, utils.NewError("Doctor not found", http.StatusNotFound)

	// }

	err := row.Scan(&docInfo.Id, &docInfo.Age, &docInfo.DoctorId, &docInfo.Speciality, &docInfo.City, &docInfo.Phone, &docInfo.ImgUrl, &docInfo.ProfessionalStatement, &docInfo.Experience, &docInfo.SatisfactionScore, &docInfo.FirstName, &docInfo.LastName, &docInfo.Governorate)
	if err != nil && err != pgx.ErrNoRows {
		return nil, utils.NewError(fmt.Sprintf("Error getting doctor info: %v", err.Error()), http.StatusInternalServerError)
	}
	if err == pgx.ErrNoRows {
		return nil, utils.NewError("Doctor not found", http.StatusNotFound)
	}
	return docInfo, nil

}

func (db *PostgresDb) GetClinicInfo(docotId int) (*utils.DoctorClinicInfo, *utils.ApiError) {
	row := db.connPool.QueryRow(db.cxt, `select * from doctors_clinic where doctor_id = $1`, docotId)

	var docInfo = new(utils.DoctorClinicInfo)
	err := row.Scan(&docInfo.Id, &docInfo.Governorate, &docInfo.City, &docInfo.Address, &docInfo.Latitude, &docInfo.Longitude, &docInfo.StartTime, &docInfo.EndTime,
		&docInfo.WorkingHours, &docInfo.Holiday, &docInfo.Location, &docInfo.PricePerAppointment, &docInfo.PatientNumberPerDay, &docInfo.DoctorId, &docInfo.AppointmentDuration)
	if err != nil && err != pgx.ErrNoRows {
		return nil, utils.NewError(fmt.Sprintf("Error getting clinic info: %v", err.Error()), http.StatusInternalServerError)
	}
	if err == pgx.ErrNoRows {
		return nil, utils.NewError("Clinic not found", http.StatusNotFound)

	}

	return docInfo, nil
}

func (db *PostgresDb) AddDoctorDetail(doctorInfo utils.DoctorInfoRequest, doctorClinic utils.DoctorClinicInfoRequest) *utils.ApiError {
	// we can use transaction here
	tx, err := db.connPool.Begin(db.cxt)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error starting transaction: %v", utils.HandlePgxError(err)), http.StatusInternalServerError)
	}
	_, execErr := tx.Exec(db.cxt, `insert into doctors_info 
	(age, doctor_id, speciality, city, phone, img_url, professional_statement, experience, satisfaction_score, first_name, last_name, governorate) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		doctorInfo.Age, doctorInfo.DoctorId, doctorInfo.Speciality, doctorInfo.City, doctorInfo.Phone, doctorInfo.ImgUrl, doctorInfo.ProfessionalStatement, doctorInfo.Experience, doctorInfo.SatisfactionScore, doctorInfo.FirstName, doctorInfo.LastName, doctorInfo.Governorate)
	if execErr != nil {
		return utils.NewError(fmt.Sprintf("[TRANSACTION] Error adding doctor info: %v", utils.HandlePgxError(err)), http.StatusInternalServerError)
	}

	// add clinic info
	_, execErr2 := tx.Exec(db.cxt, `insert into doctors_clinic_info 
	(governorate, city, address, latitude, longitude, start_time, end_time, working_hours, holiday, 
	 price_per_appointment, patient_number_per_day, doctor_id, appointment_duration ) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		doctorClinic.Governorate, doctorClinic.City, doctorClinic.Address, doctorClinic.Latitude, doctorClinic.Longitude, doctorClinic.StartTime, doctorClinic.EndTime,
		doctorClinic.WorkingHours, doctorClinic.Holiday, doctorClinic.PricePerAppointment, doctorClinic.PatientNumberPerDay, doctorClinic.DoctorId, doctorClinic.AppointmentDuration)

	if execErr2 != nil {
		return utils.NewError(fmt.Sprintf("[TRANSACTION] Error adding clinic info: %v", utils.HandlePgxError(err)), http.StatusInternalServerError)
	}

	// commit transaction
	err = tx.Commit(db.cxt)
	if err != nil {

		return utils.NewError(fmt.Sprintf("[TRANSACTION] Error commiting transaction: %v", utils.HandlePgxError(err)), http.StatusInternalServerError)
	}
	// defer tx.Rollback(db.cxt) // print err if rollback fail

	return nil

	// we can use context here
	// we can use defer here
}
func (db *PostgresDb) UpdateDoctorInfo(doctorInfo utils.DoctorInfoRequest) (*utils.DoctorInfo, *utils.ApiError) {
	row := db.connPool.QueryRow(db.cxt, `update doctors_info set age=$1, speciality=$2, city=$3, phone=$4, img_url=$5, professional_statement=$6, experience=$7, satisfaction_score=$8, first_name=$9, last_name=$10, governorate=$11 where doctor_id=$12 returning *`,
		doctorInfo.Age, doctorInfo.Speciality, doctorInfo.City, doctorInfo.Phone, doctorInfo.ImgUrl, doctorInfo.ProfessionalStatement, doctorInfo.Experience, doctorInfo.SatisfactionScore, doctorInfo.FirstName, doctorInfo.LastName, doctorInfo.Governorate, doctorInfo.DoctorId)
	doc := new(utils.DoctorInfo)
	err := row.Scan(&doc.Id, &doc.Age, &doc.DoctorId, &doc.Speciality, &doc.City, &doc.Phone, &doc.ImgUrl, &doc.ProfessionalStatement, &doc.Experience, &doc.SatisfactionScore, &doc.FirstName, &doc.LastName, &doc.Governorate)
	if err != nil && err != pgx.ErrNoRows {
		return nil, &utils.ApiError{Code: http.StatusInternalServerError, Error: err.Error()}
	}
	if err == pgx.ErrNoRows {
		return nil, &utils.ApiError{Code: http.StatusNotFound, Error: "Doctor not found"}
	}
	return doc, nil
}

func (db *PostgresDb) UpdateClinicInfo(doctorInfo utils.DoctorClinicInfoRequest) (*utils.DoctorClinicInfo, *utils.ApiError) {
	row := db.connPool.QueryRow(db.cxt, `update doctors_clinic set governorate=$1, city=$2, address=$3, latitude=$4, longitude=$5, start_time=$6, end_time=$7, working_hours=$8, holiday=$9, 
	 price_per_appointment=$10, patient_number_per_day=$11, doctor_id=$12, appointment_duration=$13 returning *`,
		doctorInfo.Governorate, doctorInfo.City, doctorInfo.Address, doctorInfo.Latitude, doctorInfo.Longitude, doctorInfo.StartTime, doctorInfo.EndTime,
		doctorInfo.WorkingHours, doctorInfo.Holiday, doctorInfo.PricePerAppointment, doctorInfo.PatientNumberPerDay, doctorInfo.DoctorId, doctorInfo.AppointmentDuration)

	clinicInfo := new(utils.DoctorClinicInfo)
	err := row.Scan(&clinicInfo.Id, &clinicInfo.Governorate, &clinicInfo.City, &clinicInfo.Address, &clinicInfo.Latitude, &clinicInfo.Longitude, &clinicInfo.StartTime, &clinicInfo.EndTime,
		&clinicInfo.WorkingHours, &clinicInfo.Holiday, &clinicInfo.Location, &clinicInfo.PricePerAppointment, &clinicInfo.PatientNumberPerDay, &clinicInfo.DoctorId, &clinicInfo.AppointmentDuration)
	if err != nil && err != pgx.ErrNoRows {
		return nil, &utils.ApiError{Code: http.StatusInternalServerError, Error: err.Error()}
	}
	if err == pgx.ErrNoRows {
		return nil, &utils.ApiError{Code: http.StatusNotFound, Error: "Clinic not found"}
	}
	return clinicInfo, nil
}

// creating booking functionalities
func (db *PostgresDb) CreateSchedule(date []time.Time, doctorId int, startTime time.Time, endTime time.Time, appointmentDuration time.Duration) *utils.ApiError {
	if len(date) == 0 {
		return utils.NewError("Date length is zero !!", http.StatusBadRequest)
	}
	tx, err := db.connPool.Begin(db.cxt)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error starting transaction: %v", err.Error()), http.StatusInternalServerError)
	}
	defer tx.Rollback(db.cxt)
	for _, v := range date {
		fmt.Print(v)
		beginTime := startTime
		for beginTime.Before(endTime) {
			_, err = tx.Exec(db.cxt,
				` insert into 
				doctor_schedules (date, start_time, end_time, status, doctor_id) values ($1, $2, $3, $4, $5)  `,
				time.Date(2025, 1, 29, 0, 0, 0, 0, time.Local), beginTime, beginTime.Add(appointmentDuration), utils.Available, doctorId)
			if err != nil {

				return utils.NewError(fmt.Sprintf("Error inserting schedule: %v", err.Error()), http.StatusInternalServerError)
			}
			beginTime = beginTime.Add(appointmentDuration)
		}

	}

	err = tx.Commit(db.cxt)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error committing transaction: %v", err.Error()), http.StatusInternalServerError)
	}
	return nil

}

func (db *PostgresDb) GetAvailableTimeSlot(doctorId int, date []time.Time) ([]utils.Schedule, *utils.ApiError) {
	allScheules := make([]utils.Schedule, 0)
	for _, v := range date {
		tx, err := db.connPool.Begin(db.cxt)
		defer tx.Rollback(db.cxt)
		if err != nil {
			return nil, utils.NewError(fmt.Sprintf("Error starting transaction: %v", err.Error()), http.StatusInternalServerError)
		}
		rows, err := tx.Query(db.cxt,
			`select * from doctor_schedules where doctor_id = $1 and date = $2`, doctorId, v)
		if err != nil {
			tx.Rollback(db.cxt)
			return nil, utils.NewError(fmt.Sprintf("Error querying schedules: %v", err.Error()), http.StatusInternalServerError)
		}
		schedule := new(utils.Schedule)
		for rows.Next() {
			err = rows.Scan(&schedule.Id, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Status, &schedule.DoctorId)
			if err != nil {
				return nil, utils.NewError(fmt.Sprintf("Error Scanning Scheules %v", err.Error()), http.StatusInternalServerError)
			}
			allScheules = append(allScheules, *schedule)
		}
	}
	return allScheules, nil
}

func (db *PostgresDb) BookSlot(patientId int, scheduleId int) *utils.ApiError {
	_, err := db.connPool.Exec(db.cxt, `update doctor_schedules status= $1 where id = $2`, utils.Booked, scheduleId)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error booking slot: %v", err.Error()), http.StatusInternalServerError)
	}
	return nil
}

func (db *PostgresDb) CancelSlot(scheduleId int) *utils.ApiError {
	_, err := db.connPool.Exec(db.cxt, `update doctor_schedules status= $1 where id = $2`, utils.Available, scheduleId)
	if err != nil {
		return utils.NewError(fmt.Sprintf("Error canceling slot: %v", err.Error()), http.StatusInternalServerError)
	}
	return nil
}
func (db *PostgresDb) GetDoctorSchedules(doctorId int, startDate time.Time, endDate *time.Time) ([]utils.Schedule, *utils.ApiError) {
	allSchedules := make([]utils.Schedule, 0)
	if endDate == nil {
		rows, err := db.connPool.Query(db.cxt, `select * from doctor_schedules where doctor_id =$1 and date =$2 order by start_time asc`, doctorId, startDate)
		if err != nil && err != pgx.ErrNoRows {
			return nil, utils.NewError(fmt.Sprintf("Error querying schedules: %v", err.Error()), http.StatusInternalServerError)
		}
		if err == pgx.ErrNoRows {
			return allSchedules, nil
		}

		for rows.Next() {
			schedule := new(utils.Schedule)
			err = rows.Scan(&schedule.Id, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Status, &schedule.PatientId, &schedule.DoctorId, &schedule.CreatedAt)
			if err != nil {
				return nil, utils.NewError(fmt.Sprintf("Error Scanning Scheules %v", err.Error()), http.StatusInternalServerError)
			}
			allSchedules = append(allSchedules, *schedule)
		}
	} else {
		rows, err := db.connPool.Query(db.cxt, `select * from doctor_schedules where doctor_id =$1 and date between $2 and $3 order by date asc`, doctorId, startDate, endDate)
		if err != nil && err != pgx.ErrNoRows {
			return nil, utils.NewError(fmt.Sprintf("Error querying schedules: %v", err.Error()), http.StatusInternalServerError)
		}
		if err == pgx.ErrNoRows {
			return allSchedules, nil
		}

		for rows.Next() {
			schedule := new(utils.Schedule)
			err = rows.Scan(&schedule.Id, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Status, &schedule.PatientId, &schedule.DoctorId, &schedule.CreatedAt)
			if err != nil {
				return nil, utils.NewError(fmt.Sprintf("Error Scanning Scheules %v", err.Error()), http.StatusInternalServerError)
			}
			allSchedules = append(allSchedules, *schedule)
		}

	}

	return allSchedules, nil
}
