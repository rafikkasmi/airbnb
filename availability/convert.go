package availability

import (
	"strconv"
	"strings"
)

// standardize converts the raw availability calendar data to our clean format
func (r root) standardize(roomId int64) AvailabilityData {
	var availabilityData AvailabilityData
	
	// Process calendar months
	for _, month := range r.Data.Merlin.PdpAvailabilityCalendar.CalendarMonths {
		calendarMonth := CalendarMonth{
			Month: month.Month,
			Year:  month.Year,
		}
		
		// Process days in the month
		for _, day := range month.Days {
			formattedPrice := ""
			var amount float32
			var currency string
			
			// Extract price if available
			if priceStr, ok := day.Price.LocalPriceFormatted.(string); ok {
				formattedPrice = priceStr
				
				// Try to parse the price amount and currency
				parts := strings.Fields(priceStr)
				if len(parts) > 0 {
					// Remove any non-numeric characters except decimal point
					numericStr := ""
					for _, char := range parts[0] {
						if (char >= '0' && char <= '9') || char == '.' {
							numericStr += string(char)
						} else if char != ',' {
							// If not a comma (which we ignore), it might be a currency symbol
							currency = string(char)
						}
					}
					
					// Parse the numeric string to float
					if amountFloat, err := strconv.ParseFloat(numericStr, 32); err == nil {
						amount = float32(amountFloat)
					}
				}
			}
			
			calendarDay := CalendarDay{
				Date:                day.CalendarDate,
				Available:           day.Available,
				MaxNights:           day.MaxNights,
				MinNights:           day.MinNights,
				AvailableForCheckin: day.AvailableForCheckin,
				AvailableForCheckout: day.AvailableForCheckout,
				Price: DayPrice{
					Formatted: formattedPrice,
					Amount:    amount,
					Currency:  currency,
				},
			}
			
			calendarMonth.Days = append(calendarMonth.Days, calendarDay)
		}
		
		availabilityData.CalendarMonths = append(availabilityData.CalendarMonths, calendarMonth)
	}
	
	return availabilityData
}
