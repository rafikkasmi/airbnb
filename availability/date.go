package availability

import (
	"fmt"
	"time"
)

// getStringDate formats a time.Time into a string date in YYYY-MM-DD format
func getStringDate(date time.Time) string {
	day := fmt.Sprintf("%d", date.Day())
	month := fmt.Sprintf("%d", date.Month())
	year := fmt.Sprintf("%d", date.Year())
	if len(day) == 1 {
		day = "0" + day
	}
	if len(month) == 1 {
		month = "0" + month
	}
	if len(year) == 1 {
		year = "0" + year
	}
	dateToUse := fmt.Sprintf("%s-%s-%s", year, month, day)
	return dateToUse
}

// parseCalendarDate parses a date string from the calendar format (YYYY-MM-DD) to time.Time
func parseCalendarDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// isDateAvailable checks if a specific date is available in the availability calendar
func isDateAvailable(availabilityData AvailabilityData, date time.Time) bool {
	dateStr := getStringDate(date)
	
	for _, month := range availabilityData.CalendarMonths {
		for _, day := range month.Days {
			if day.Date == dateStr {
				return day.Available
			}
		}
	}
	
	return false
}

// getAvailableDatesInRange returns all available dates within a specified date range
func getAvailableDatesInRange(availabilityData AvailabilityData, startDate, endDate time.Time) []time.Time {
	var availableDates []time.Time
	
	// Ensure startDate is before endDate
	if startDate.After(endDate) {
		startDate, endDate = endDate, startDate
	}
	
	// Create a map for faster lookups
	availableDaysMap := make(map[string]bool)
	for _, month := range availabilityData.CalendarMonths {
		for _, day := range month.Days {
			if day.Available {
				availableDaysMap[day.Date] = true
			}
		}
	}
	
	// Iterate through the date range and check availability
	current := startDate
	for !current.After(endDate) {
		dateStr := getStringDate(current)
		if availableDaysMap[dateStr] {
			availableDates = append(availableDates, current)
		}
		current = current.AddDate(0, 0, 1) // Add one day
	}
	
	return availableDates
}
