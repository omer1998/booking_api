package doctor

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/omer1998/booking_api/fixtures"
	"github.com/omer1998/booking_api/utils"
	"gotest.tools/assert"
	// _ "github.com/golang-migrate/migrate/v4/database/postgres"
	// _ "github.com/golang-migrate/migrate/v4/source/github"
	// "github.com/omer1998/booking_api/fixtures"
)

func TestDoctor(t *testing.T) {
	// we need to pass the path of the .env file to the NewTestEnv function in order to uplaod the env variables in config.newwithpath function
	wd, err := os.Getwd()
	if err != nil {
		t.Error("Error getting working directory")
	}
	t.Log("Working directory: ", wd)

	path := filepath.Join(wd, "..", ".env")
	if os.IsNotExist(err) {
		t.Error(".env file does not exist")
	} else if err != nil {
		t.Errorf("Error accessing .env file: %v", err)
	} else {
		t.Log(".env file exists")
		t.Log("Path: ", path)
	}
	testEnv := fixtures.NewTestEnv(t, path)
	testEnv.TestAddDoctor(path)
}

func TestAddDoctorInfo(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error("Error getting working directory")
	}
	t.Log("Working directory: ", wd)

	path := filepath.Join(wd, "..", ".env")
	if os.IsNotExist(err) {
		t.Error(".env file does not exist")
	} else if err != nil {
		t.Errorf("Error accessing .env file: %v", err)
	} else {
		t.Log(".env file exists")
		t.Log("Path: ", path)
	}
	testEnv := fixtures.NewTestEnv(t, path)
	testEnv.SetUpDb()
	t.Cleanup(testEnv.TearDownDb)
	if err := testEnv.Db.AddDoctor("omer@omer.com", "12345678"); err != nil {
		t.Error("Error adding doctor", err)
	}
	doctor, er := testEnv.Db.GetDoctorByEmail("omer@omer.com")
	if er != nil {
		t.Error("Error getting doctor", er)
	}
	t.Log("Doctor: ", doctor)

	experience := "5 years"
	docInfo := utils.DoctorInfoRequest{
		DoctorId:              doctor.Id,
		FirstName:             "Omer",
		LastName:              "Faris",
		Phone:                 "07710210244",
		Governorate:           "Baghdad",
		City:                  "Al Dora",
		Speciality:            "Internal Medicine",
		Age:                   27,
		Experience:            experience,
		SatisfactionScore:     1,
		ProfessionalStatement: "",
		ImgUrl:                "",
	}
	if err := testEnv.Db.AddDoctorInfo(docInfo); err != nil {
		t.Error("Error adding doctor info", err)
	}

	docInfo2, gErr := testEnv.Db.GetDoctorInfo(doctor.Id)
	if gErr != nil {
		t.Error("Error getting doctor info", gErr)
	}
	t.Log("Doctor Info after addition to db and fetching: ", docInfo2)
	assert.Equal(t, docInfo.FirstName, docInfo2.FirstName)
	assert.Equal(t, docInfo.LastName, docInfo2.LastName)
	assert.Equal(t, docInfo.Phone, docInfo2.Phone)

}

func TestUpdateDoctorInfo(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error("Error getting working directory")
	}
	t.Log("Working directory: ", wd)

	path := filepath.Join(wd, "..", ".env")
	if os.IsNotExist(err) {
		t.Error(".env file does not exist")
	} else if err != nil {
		t.Errorf("Error accessing .env file: %v", err)
	} else {
		t.Log(".env file exists")
		t.Log("Path: ", path)
	}
	testEnv := fixtures.NewTestEnv(t, path)
	testEnv.SetUpDb()
	t.Cleanup(testEnv.TearDownDb)
	if err := testEnv.Db.AddDoctor("omer@omer.com", "12345678"); err != nil {
		t.Error("Error adding doctor", err)
	}
	doctor, er := testEnv.Db.GetDoctorByEmail("omer@omer.com")
	if er != nil {
		t.Error("Error getting doctor", er)
	}
	t.Log("Doctor: ", doctor)

	experience := "5 years"
	docInfo := utils.DoctorInfoRequest{
		DoctorId:              doctor.Id,
		FirstName:             "Omer",
		LastName:              "Faris",
		Phone:                 "07710210244",
		Governorate:           "Baghdad",
		City:                  "Al Dora",
		Speciality:            "Internal Medicine",
		Age:                   27,
		Experience:            experience,
		SatisfactionScore:     1,
		ProfessionalStatement: "",
		ImgUrl:                "",
	}
	if err := testEnv.Db.AddDoctorInfo(docInfo); err != nil {
		t.Error("Error adding doctor info", err)
	}

	updateDocInfo := utils.DoctorInfoRequest{
		DoctorId:              doctor.Id,
		FirstName:             "Omer Faris",
		LastName:              "Nawar",
		Phone:                 "07710210244",
		Governorate:           "Baghdad",
		City:                  "Al Dora",
		Speciality:            "Internal Medicine",
		Age:                   27,
		Experience:            "10 years NHS EXPERIENCE",
		SatisfactionScore:     1,
		ProfessionalStatement: "",
		ImgUrl:                "",
	}
	updatedDoc, updateErr := testEnv.Db.UpdateDoctorInfo(updateDocInfo)
	if updateErr != nil {
		t.Error("Error updating doctor info", updateErr)
	}
	assert.Equal(t, updateDocInfo.FirstName, updatedDoc.FirstName)
	t.Log("Updated Doctor Info: ", updatedDoc)

}

func TestAddAndGetClinicInfo(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error("Error getting working directory")
	}
	t.Log("Working directory: ", wd)

	path := filepath.Join(wd, "..", ".env")
	if os.IsNotExist(err) {
		t.Error(".env file does not exist")
	} else if err != nil {
		t.Errorf("Error accessing .env file: %v", err)
	} else {
		t.Log(".env file exists")
		t.Log("Path: ", path)
	}
	testEnv := fixtures.NewTestEnv(t, path)
	testEnv.SetUpDb()
	t.Cleanup(testEnv.TearDownDb)
	if err := testEnv.Db.AddDoctor("omer@omer.com", "12345678"); err != nil {
		t.Error("Error adding doctor", err)
	}
	doctor, er := testEnv.Db.GetDoctorByEmail("omer@omer.com")
	if er != nil {
		t.Error("Error getting doctor", er)
	}
	t.Log("Doctor: ", doctor)

	startTime, _ := time.Parse("15:04", "15:00")
	endTime, _ := time.Parse("15:04", "20:00")
	t.Log("Start Time: ", startTime.Hour(), "End Time: ", endTime.Hour())

	clinicInfo := utils.DoctorClinicInfoRequest{
		DoctorId:     doctor.Id,
		Governorate:  "Baghdad",
		City:         "Al Dora",
		Address:      "Al Dora",
		Latitude:     "33.252515",
		Longitude:    "44.393449",
		StartTime:    startTime,
		EndTime:      endTime,
		WorkingHours: 5,
		Holiday:      []int{1, 2},

		PricePerAppointment: 25,
		PatientNumberPerDay: 15,
		AppointmentDuration: 20,
	}
	if err := testEnv.Db.AddClinicInfo(clinicInfo); err != nil {
		t.Error("Error adding clinic info", err)
	}

	getClinincInfo, getErr := testEnv.Db.GetClinicInfo(doctor.Id)
	if getErr != nil {
		t.Error("Error getting clinic info", getErr)
	}
	t.Log("Clinic Info: ", getClinincInfo)
	t.Log("location :", getClinincInfo.Location)
	t.Log("start time: ", getClinincInfo.StartTime, "end time: ", getClinincInfo.EndTime)

}

func TestUpdateClinicInfo(t *testing.T) {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	t.Error("Error getting working directory")
	// }
	// t.Log("Working directory: ", wd)

	// path := filepath.Join(wd, "..", "..", ".env")
	// if os.IsNotExist(err) {
	// 	t.Error(".env file does not exist")
	// } else if err != nil {
	// 	t.Errorf("Error accessing .env file: %v", err)
	// } else {
	// 	t.Log(".env file exists")
	// 	t.Log("Path: ", path)
	// }
	path := getEnvPath("..", "..")
	testEnv := fixtures.NewTestEnv(t, path)
	testEnv.SetUpDb()
	t.Cleanup(testEnv.TearDownDb)
	if err := testEnv.Db.AddDoctor("omer@omer.com", "12345678"); err != nil {
		t.Error("Error adding doctor", err)
	}
	doctor, er := testEnv.Db.GetDoctorByEmail("omer@omer.com")
	if er != nil {
		t.Error("Error getting doctor", er)
	}
	t.Log("Doctor: ", doctor)
	startTime, _ := time.Parse("15:04", "15:00")
	endTime, _ := time.Parse("15:04", "20:00")
	t.Log("Start Time: ", startTime.Hour(), "End Time: ", endTime.Hour())

	clinicInfo := utils.DoctorClinicInfoRequest{
		DoctorId:     doctor.Id,
		Governorate:  "Baghdad",
		City:         "Al Dora",
		Address:      "Al Dora",
		Latitude:     "33.252515",
		Longitude:    "44.393449",
		StartTime:    startTime,
		EndTime:      endTime,
		WorkingHours: 5,
		Holiday:      []int{1, 2},

		PricePerAppointment: 25,
		PatientNumberPerDay: 15,
		AppointmentDuration: 20,
	}
	if err := testEnv.Db.AddClinicInfo(clinicInfo); err != nil {
		t.Error("Error adding clinic info", err)
	}
	newClinicInfo := utils.DoctorClinicInfoRequest{
		DoctorId:     doctor.Id,
		Governorate:  "Baghdad",
		City:         "Al Dora",
		Address:      "Al Dora",
		Latitude:     "33.252515",
		Longitude:    "44.393449",
		StartTime:    startTime,
		EndTime:      endTime,
		WorkingHours: endTime.Hour() - startTime.Hour(),
		Holiday:      []int{1, 2},

		PricePerAppointment: 25,
		PatientNumberPerDay: 15,
		AppointmentDuration: 20,
	}

	newInfo, updateErr := testEnv.Db.UpdateClinicInfo(newClinicInfo)
	if updateErr != nil {
		t.Error("Error updating clinic info", updateErr)
	}
	t.Log("Updated Clinic Info: ", newInfo)
}

// you need to update the path of the .env file to make the test work
// a function to get the path of the .env file
// where path is the steps (folder name) to get to the .env file
func getEnvPath(paths ...string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic("Error getting working directory")
	}
	for _, path := range paths {
		wd = filepath.Join(wd, path)

	}

	wd = filepath.Join(wd, ".env")
	return wd

}
func TestGetEnvPath(t *testing.T) {
	path := getEnvPath("..", "..")
	actualPathToEnv := `c:\Users\master\Desktop\booking_api\.env`
	t.Log("Path: ", path)
	assert.Equal(t, actualPathToEnv, path)

}
