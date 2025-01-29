package booking

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/omer1998/booking_api/fixtures"
	"github.com/stretchr/testify/assert"
)

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

func TestMakeSchedules(t *testing.T) {
	t.Log("TestMakeSchedules")
	path := getEnvPath("..", "..")
	testEnv := fixtures.NewTestEnv(t, path)
	testEnv.SetUpDb()
	defer testEnv.TearDownDb()
	if err := testEnv.Db.AddDoctor("ahmed@ahmed.com", "12345678"); err != nil {
		t.Error("Error adding doctor", err)
	}
	doctor, err := testEnv.Db.GetDoctorByEmail("ahmed@ahmed.com")
	if err != nil {
		t.Error("Error getting doctor", err)
	}
	t.Log("Doctor: ", doctor)

	// setting availability -- creating schedules
	// t.Run("SettingAvailability", func(t *testing.T) {

	// })
	startTime, _ := time.Parse("15:04", "15:00")
	endTime, _ := time.Parse("15:04", "20:00")
	date := time.Date(2025, 1, 29, 0, 0, 0, 0, time.Local)
	t.Log("date", []time.Time{date})
	t.Log("startTime", startTime)
	t.Log("endTime", endTime)
	appDuration := time.Duration(15) * time.Minute

	err = testEnv.Db.CreateSchedule([]time.Time{date}, doctor.Id, startTime, endTime, appDuration)
	if err != nil {
		t.Log("Error creating schedule for doctor: ", err)
	}
	// get total number of scheduls of doctor in a day
	row := testEnv.ConnPool.QueryRow(context.Background(),
		`select count(*) from doctor_schedules where doctor_id = $1 and date = $2`,
		doctor.Id, date,
	)
	var totalNum int
	errS := row.Scan(&totalNum)
	if err != nil {
		t.Error("Error getting total number of schedules", errS)
	}
	t.Log("totalNum", totalNum)
	assert.Equal(t, int(endTime.Sub(startTime)/appDuration), totalNum)
	// assert.NoError(t, errors.New(err.Error))
}
