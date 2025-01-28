package services

import (
	"time"
)

type TimeSlot struct {
	ID        int       `json:"id"`
	DoctorID  int       `json:"doctor_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Date      time.Time `json:"date"`
	IsBooked  bool      `json:"is_booked"`
}

func GenerateTimeSlots(doctorID int, startTime, endTime time.Time, intervalMinutes int) []TimeSlot {
	var timeSlots []TimeSlot
	currentTime := startTime

	for currentTime.Before(endTime) {
		timeSlot := TimeSlot{
			DoctorID:  doctorID,
			StartTime: currentTime,
			EndTime:   currentTime.Add(time.Duration(intervalMinutes) * time.Minute),
			Date:      startTime,
			IsBooked:  false,
		}
		timeSlots = append(timeSlots, timeSlot)
		currentTime = currentTime.Add(time.Duration(intervalMinutes) * time.Minute)
	}

	return timeSlots
}

func MarkTimeSlotAsBooked(timeSlotID int) error {
	// Implement the logic to mark the time slot as booked in the database
	// Example:
	// UPDATE time_slots SET is_booked = true WHERE id = timeSlotID
	return nil
}
