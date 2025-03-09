package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gobnb"
	"gobnb/availability"
	"gobnb/details"
	"gobnb/reviews"
	"gobnb/search"

	// "github.com/johnbalvin/gobnb/utils"
	// "github.com/johnbalvin/gobnb/details"
	"github.com/gocarina/gocsv"
)

var client gobnb.Client

func main() {
	client = gobnb.DefaultClient()
	// Uncomment one of these function calls to run different examples

	// Example 1: Search for rooms
	searchForRooms()

	// Example 2: Get room details
	// getRooms()

	// Example 3: Get reviews for a room
	// getReviews(290701)

	// getRoomYearAvailability(290701)
}

// citySearchResult holds the result of a city search operation
type citySearchResult struct {
	results []search.Data
	city    string
	err     error
}

// searchForRooms searches for rooms in multiple locations asynchronously and gets their details
func searchForRooms() {
	// List of cities to search in
	Cities := []string{
		"Menara, Marrakesh",
		"Gueliz, Marrakesh",
		"Medina, Marrakesh",
		"Sidi Youssef Ben Ali, Marrakesh",
		"Annakhil, Marrakesh",
		"Mechouar Kasba, Marrakesh",
		"Saada, Marrakesh",
		"Tassoultante, Marrakesh",
		"Loudaya, Marrakesh",
		"Alouidane, Marrakesh",
		"Souihla, Marrakesh",
		"Oulad Hassoune, Marrakesh",
		"Harbil, Marrakesh",
		"Ouled Dlim, Marrakesh",
		"Ouahat Sidi Brahim, Marrakesh",
		"Ait Imour, Marrakesh",
		"M'Nabha, Marrakesh",
		"Sid Zouine, Marrakesh",
		"Agafay, Marrakesh",
		"Bab Ghmat, Marrakesh",
		//neighboorhoods
		// "Arset El Baraka, Marrakesh",
		// "Arset Moulay Bouaza, Marrakesh",
		// "Djane Ben Chogra, Marrakesh",
		// "Arset El Houta, Marrakesh",
		// "Bab Aylan, Marrakesh",
		// "Arset Sidi Youssef, Marrakesh",
		// "Derb Chtouka, Marrakesh",
		// "Bab Hmar, Marrakesh",
		// "Bab Agnaou, Marrakesh",
		// "Quartier Jnan Laafia, Marrakesh",
		// "Toureg, Marrakesh",
		// "Kasbah, Marrakesh",
		// "Mellah, Marrakesh",
		// "Arset El Maach, Marrakesh",
		// "Arset Moulay Moussa, Marrakesh",
		// "Riad Zitoun Jdid, Marrakesh",
		// "Kennaria, Marrakesh",
		// "Rahba Kedima, Marrakesh",
		// "Kaat Benahid, Marrakesh",
		// "Zaouiat Lahdar, Marrakesh",
		// "El Moukef, Marrakesh",
		// "Riad Laarous, Marrakesh",
		// "Assouel, Marrakesh",
		// "Kechich, Marrakesh",
		// "Douar Fekhara, Marrakesh",
		// "Arset Tihiri, Marrakesh",
		// "Sidi Ben Slimane El Jazouli, Marrakesh",
		// "Diour Jdad, Marrakesh",
		// "Rmila, Marrakesh",
		// "Zaouia Sidi Rhalem, Marrakesh",
		// "Kbour Chou, Marrakesh",
		// "Ain Itti, Marrakesh",
		// "Bab Doukkala, Marrakesh",
		// "El Hara, Marrakesh",
		// "Arset El Bilk, Marrakesh",
	}

	fmt.Printf("Starting asynchronous search for rooms in %d cities...\n", len(Cities))

	// Create channels for search results and synchronization
	citiesResultsChan := make(chan citySearchResult)
	var wgCities sync.WaitGroup

	// Set up a worker pool with a maximum of 5 concurrent searches
	maxConcurrentSearches := 20
	searchSemaphore := make(chan struct{}, maxConcurrentSearches)

	// Launch a goroutine for each city search
	for _, city := range Cities {
		wgCities.Add(1)

		go func(cityName string) {
			defer wgCities.Done()

			// Acquire a semaphore slot
			searchSemaphore <- struct{}{}
			defer func() { <-searchSemaphore }()

			fmt.Printf("Searching for rooms in %s...\n", cityName)

			// Set up search parameters
			zoomvalue := 1
			checkIn := search.Check{}
			coords := search.CoordinatesInput{}

			// Perform the search
			results, err := search.InputData{
				Coordinates: coords,
				Check:       checkIn,
				ZoomValue:   zoomvalue,
				Query:       cityName,
			}.SearchAll("USD", nil)

			// Send the results back through the channel
			citiesResultsChan <- citySearchResult{
				results: results,
				city:    cityName,
				err:     err,
			}
		}(city)
	}

	// Start a goroutine to close the results channel when all searches are done
	go func() {
		wgCities.Wait()
		close(citiesResultsChan)
	}()

	// Collect all search results
	results := []search.Data{}
	completedCity := 0
	totalCities := len(Cities)

	// Process results as they come in
	for result := range citiesResultsChan {
		completedCity++
		fmt.Printf("[%d/%d] Completed search for %s - ", completedCity, totalCities, result.city)

		if result.err != nil {
			fmt.Printf("Error: %s\n", result.err)
			continue
		}

		fmt.Printf("Found %d rooms\n", len(result.results))

		// Add results to our collection
		for _, room := range result.results {
			results = append(results, room)
		}
	}

	//make array unique by RoomId
	uniqueResults := make(map[int64]search.Data)
	for _, result := range results {
		uniqueResults[result.RoomID] = result
	}

	fmt.Println("Found", len(uniqueResults), "unique rooms")

	// Save search results asynchronously
	// go func() {
	rawJSON, _ := json.MarshalIndent(uniqueResults, "", "  ")
	os.Remove("./searchResult.json")
	if err := os.WriteFile("./searchResult.json", rawJSON, 0644); err != nil {
		log.Println("Error saving search results:", err)
		return
	}
	fmt.Println("Search results saved to searchResult.json")
	// }()

	// return;
	// Now get details for each room asynchronously
	fmt.Println("\nFetching details for each room asynchronously...")

	// Create a channel to receive room details
	type roomDetailResult struct {
		details details.Data
		roomID  int64
		name    string
		err     error
	}

	// Set up a worker pool with a maximum of 10 concurrent requests
	maxConcurrent := 50
	semaphore := make(chan struct{}, maxConcurrent)
	resultsChan := make(chan roomDetailResult)

	// Keep track of how many goroutines we're launching
	total := len(uniqueResults)
	active := 0

	// Use a WaitGroup to ensure all goroutines complete
	var wg sync.WaitGroup

	// Launch goroutines to fetch room details
	for roomID, roomInfo := range uniqueResults {
		active++
		wg.Add(1)

		// Launch a goroutine for each room
		go func(id int64, info search.Data) {
			// Ensure the WaitGroup is decremented when the goroutine completes
			defer wg.Done()

			// Acquire a semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Add a small random delay to avoid all requests hitting at once
			time.Sleep(time.Duration(id%100) * time.Millisecond)

			// Create output directory for this room
			folderPath := fmt.Sprintf("./output/%d", id)
			if err := os.MkdirAll(folderPath, 0755); err != nil {
				log.Printf("Error creating directory for room %d: %v", id, err)
			}

			// Get room details
			roomDetails, err := client.DetailsFromRoomID(id)

			// // Asynchronously fetch reviews and availability data
			// var reviewsWg sync.WaitGroup
			// reviewsWg.Add(2) // One for reviews, one for availability

			// // Fetch reviews asynchronously
			// go func() {
			// 	defer reviewsWg.Done()
			// 	fetchReviewsForRoom(id, folderPath)
			// }()

			// // Fetch availability asynchronously
			// go func() {
			// 	defer reviewsWg.Done()
			// 	fetchAvailabilityForRoom(id, folderPath)
			// }()

			// Don't wait for reviews and availability to complete
			// They will finish in the background

			// Send the result back through the channel
			resultsChan <- roomDetailResult{
				details: roomDetails,
				roomID:  id,
				name:    info.Name,
				err:     err,
			}
		}(roomID, roomInfo)
	}

	// Start a goroutine to close the results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	allRoomDetails := make([]details.Data, 0, total)
	completed := 0

	// Progress tracking
	fmt.Printf("Waiting for %d room details to be fetched...\n", total)

	// Collect all the results from the goroutines
	for result := range resultsChan {
		completed++

		fmt.Printf("[%d/%d] Room %d: %s - ", completed, total, result.roomID, result.name)

		if result.err != nil {
			fmt.Printf("Error: %s\n", result.err)
			continue
		}

		fmt.Println("Success")
		allRoomDetails = append(allRoomDetails, result.details)
	}

	// Save all room details to a JSON file
	fmt.Printf("\nSaving details for %d rooms to rooms_details.json\n", len(allRoomDetails))

	// Marshal the room details to JSON

	detailsJSON, err := gocsv.MarshalString(&allRoomDetails)
	if err != nil {
		log.Println("Error marshaling room details to JSON:", err)
		return
	}

	// Write the JSON to a file
	os.Remove("./rooms_details.csv")

	// Create folder for the room if it doesn't exist
	folderPath := fmt.Sprintf("./output/")
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Save all paginated reviews to a file
	filePath := fmt.Sprintf("%s/rooms_details.csv", folderPath)

	if err := os.WriteFile(filePath, []byte(detailsJSON), 0644); err != nil {
		log.Println("Error writing room details to file:", err)
		return
	}

	fmt.Println("Room details saved to rooms_details.csv")
}

func getRoomDetail(roomId int64) {

}

func getRooms() {
	var roomID int64
	roomID = 290701
	// romID:=[]int{roomID}
	data, err := client.DetailsFromRoomID(roomID)
	if err != nil {
		log.Println("test:2 -> err: ", err)
		return
	}
	rawJSON, _ := gocsv.MarshalString(&data)
	fmt.Printf("%s", rawJSON) //in case you don't have write permisions
	if err := os.WriteFile("./details.csv", []byte(rawJSON), 0644); err != nil {
		log.Println(err)
		return
	}
}

// getReviews fetches and displays reviews for a specific room ID
func getReviews(roomID int64) {
	// Use a known valid room ID with reviews
	// Using the room ID from the example URL
	// var roomID int64 = 290701

	fmt.Printf("Fetching reviews for room ID: %d\n", roomID)

	// Create a client
	// client := gobnb.DefaultClient()

	// Fetch reviews for the room
	reviewData, err := reviews.InputData{
		RoomId: roomID,
	}.GetAllReviewsOfRoom(roomID, "USD", nil)
	if err != nil {
		log.Println("Error fetching reviews:", err)
		return
	}

	// Print debug information about the response
	// fmt.Printf("Raw review data: %+v\n", reviewData)
	fmt.Println(len(reviewData))

	// // Print summary information
	// fmt.Printf("Room %d has %d reviews with an average rating of %.1f\n",
	// 	roomID, reviewData.TotalReviews, reviewData.Rating)

	// // Print the first few reviews
	// fmt.Println("Sample reviews:")
	// for i, review := range reviewData.Reviews {
	// 	if i >= 3 { // Only show first 3 reviews
	// 		break
	// 	}
	// 	fmt.Printf("  Review #%d by %s (%d stars): %s\n",
	// 		i+1, review.AuthorName, review.Rating, truncateString(review.Comments, 100))
	// }

	// // Save all reviews to a file
	// rawJSON, _ := json.MarshalIndent(reviewData, "", "  ")
	// if err := os.WriteFile("./reviews.json", rawJSON, 0644); err != nil {
	// 	log.Println("Error saving reviews:", err)
	// 	return
	// }
	// fmt.Println("All reviews saved to reviews.json")

	// // Fetch all reviews with pagination if there are more
	// if reviewData.HasMoreReviews {
	// 	fmt.Println("Fetching all reviews with pagination...")
	// 	allReviews, err := client.AllReviewsFromRoomID(roomID)
	// 	if err != nil {
	// 		log.Println("Error fetching all reviews:", err)
	// 		return
	// 	}

	// 	fmt.Printf("Successfully fetched all %d reviews\n", len(allReviews.Reviews))

	// Save all paginated reviews to a file
	// Create folder for the room if it doesn't exist
	folderPath := fmt.Sprintf("./output/%d", roomID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Save all paginated reviews to a file
	allJSON, _ := gocsv.MarshalString(&reviewData)
	filePath := fmt.Sprintf("%s/reviews.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Println("Error saving all reviews:", err)
		return
	}
	fmt.Printf("All paginated reviews saved to %s\n", filePath)
	// }
}

func getRoomYearAvailability(roomID int64) {

	fmt.Printf("Fetching availability for room ID: %d\n", roomID)

	// Create a client
	// client := gobnb.DefaultClient()

	// Fetch reviews for the room
	availabilityData, daysData, err := availability.InputData{
		RoomId:     roomID,
		StartMonth: 3,
		StartYear:  2025,
	}.GetAvailabilityCalendar("USD", nil)
	if err != nil {
		log.Println("Error fetching availbility:", err)
		return
	}

	// Print debug information about the response
	// fmt.Printf("Raw review data: %+v\n", reviewData)
	fmt.Println(len(daysData))

	fmt.Println(availabilityData)

	// // Print summary information
	// fmt.Printf("Room %d has %d reviews with an average rating of %.1f\n",
	// 	roomID, reviewData.TotalReviews, reviewData.Rating)

	// // Print the first few reviews
	// fmt.Println("Sample reviews:")
	// for i, review := range reviewData.Reviews {
	// 	if i >= 3 { // Only show first 3 reviews
	// 		break
	// 	}
	// 	fmt.Printf("  Review #%d by %s (%d stars): %s\n",
	// 		i+1, review.AuthorName, review.Rating, truncateString(review.Comments, 100))
	// }

	// // Save all reviews to a file
	// rawJSON, _ := json.MarshalIndent(reviewData, "", "  ")
	// if err := os.WriteFile("./reviews.json", rawJSON, 0644); err != nil {
	// 	log.Println("Error saving reviews:", err)
	// 	return
	// }
	// fmt.Println("All reviews saved to reviews.json")

	// // Fetch all reviews with pagination if there are more
	// if reviewData.HasMoreReviews {
	// 	fmt.Println("Fetching all reviews with pagination...")
	// 	allReviews, err := client.AllReviewsFromRoomID(roomID)
	// 	if err != nil {
	// 		log.Println("Error fetching all reviews:", err)
	// 		return
	// 	}

	// 	fmt.Printf("Successfully fetched all %d reviews\n", len(allReviews.Reviews))

	// Create folder for the room if it doesn't exist
	folderPath := fmt.Sprintf("./output/%d", roomID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Save all paginated reviews to a file
	allJSON, _ := gocsv.MarshalString(&daysData)
	filePath := fmt.Sprintf("%s/availability.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Println("Error saving availability data:", err)
		return
	}
	fmt.Printf("Availability data saved to %s\n", filePath)
	// }
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// fetchReviewsForRoom fetches reviews for a specific room ID and saves them to a file
func fetchReviewsForRoom(roomID int64, folderPath string) {
	fmt.Printf("Fetching reviews for room ID: %d\n", roomID)

	// Fetch reviews for the room
	reviewData, err := reviews.InputData{
		RoomId: roomID,
	}.GetAllReviewsOfRoom(roomID, "USD", nil)
	if err != nil {
		log.Printf("Error fetching reviews for room %d: %v\n", roomID, err)
		return
	}

	fmt.Printf("Room %d: Found %d reviews\n", roomID, len(reviewData))

	// Save reviews to a file
	allJSON, _ := gocsv.MarshalString(&reviewData)
	filePath := fmt.Sprintf("%s/reviews.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Printf("Error saving reviews for room %d: %v\n", roomID, err)
		return
	}
	fmt.Printf("Room %d: Reviews saved to %s\n", roomID, filePath)
}

// fetchAvailabilityForRoom fetches availability data for a specific room ID and saves it to a file
func fetchAvailabilityForRoom(roomID int64, folderPath string) {
	fmt.Printf("Fetching availability for room ID: %d\n", roomID)

	// Fetch availability data for the room
	_, daysData, err := availability.InputData{
		RoomId:     roomID,
		StartMonth: 3,
		StartYear:  2025,
	}.GetAvailabilityCalendar("USD", nil)
	if err != nil {
		log.Printf("Error fetching availability for room %d: %v\n", roomID, err)
		return
	}

	fmt.Printf("Room %d: Found %d availability days\n", roomID, len(daysData))

	// Save availability data to a file
	allJSON, _ := gocsv.MarshalString(&daysData)
	filePath := fmt.Sprintf("%s/availability.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Printf("Error saving availability for room %d: %v\n", roomID, err)
		return
	}
	fmt.Printf("Room %d: Availability data saved to %s\n", roomID, filePath)
}
