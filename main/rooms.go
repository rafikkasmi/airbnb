package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"gobnb"
	"gobnb/details"
	"gobnb/search"

	// "github.com/johnbalvin/gobnb/utils"
	// "github.com/johnbalvin/gobnb/details"
	"github.com/gocarina/gocsv"
)

var (
	maxConcurrentSearch = 20
	maxConcurrentRooms  = 20
)

var CITIES = []string{
	// "Marrakesh",
	"Gueliz, Marrakesh",
	// "Medina	, Marrakesh",
	// "Sidi Youssef Ben Ali, Marrakesh",
	// "Annakhil, Marrakesh",
	// "Mechouar Kasba, Marrakesh",
	// "Saada, Marrakesh",
	// "Tassoultante, Marrakesh",
	// "Loudaya, Marrakesh",
	// "Alouidane, Marrakesh",
	// "Souihla, Marrakesh",
	// "Oulad Hassoune, Marrakesh",
	// "Harbil, Marrakesh",
	// "Ouled Dlim, Marrakesh",
	// "Ouahat Sidi Brahim, Marrakesh",
	// "Ait Imour, Marrakesh",
	// "M'Nabha, Marrakesh",
	// "Sid Zouine, Marrakesh",
	// "Agafay, Marrakesh",
	// "Bab Ghmat, Marrakesh",
	// //neighboorhoods
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

var (
	client gobnb.Client
	// Global variable to store room details from CSV
	existingRooms []details.Data
	// Map for quick room ID lookups
	existingRoomIDs map[int64]bool
)

// loadRoomDetailsCSV loads the room details from CSV file into memory
func loadRoomDetailsCSV() error {
	// Initialize the map
	existingRoomIDs = make(map[int64]bool)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll("./output", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Try to open the CSV file
	filePath := "./output/rooms_details.csv"
	file, err := os.Open(filePath)
	if err != nil {
		// If file doesn't exist, initialize empty slices and return
		if os.IsNotExist(err) {
			existingRooms = []details.Data{}
			return nil
		}
		return fmt.Errorf("error opening rooms_details.csv: %w", err)
	}
	defer file.Close()

	// Parse the CSV file
	if err := gocsv.UnmarshalFile(file, &existingRooms); err != nil {
		return fmt.Errorf("error parsing rooms_details.csv: %w", err)
	}

	// Build the map of room IDs for quick lookup
	for _, room := range existingRooms {
		existingRoomIDs[room.RoomID] = true
	}

	fmt.Printf("Loaded %d rooms from rooms_details.csv\n", len(existingRooms))
	return nil
}

func main() {
	client = gobnb.DefaultClient()

	// Load room details from CSV at startup
	if err := loadRoomDetailsCSV(); err != nil {
		log.Printf("Warning: Failed to load room details: %v\n", err)
	}

	// Example 1: Search for rooms
	searchForRooms()

	// Example 2: Get room details
	// getRooms()

	// Example 3: Get reviews for a room
	// getReviews(290701)

	// getRoomYearAvailability(290701)

	// Example: Check if a room ID exists and update the CSV file
	// checkAndUpdateRoomDetails(1234567890)
}

// checkRoomExists checks if a room ID exists in the roomDetails.csv file
func checkRoomExists(roomID int64) (bool, error) {
	// Simply check the map for the room ID
	return existingRoomIDs[roomID], nil
}

// updateRoomDetailsCSV updates the roomDetails.csv file with newly added rooms
func updateRoomDetailsCSV(newRooms []details.Data) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll("./output", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	filePath := "./output/rooms_details.csv"

	// Filter out rooms that already exist using our global map
	var uniqueNewRooms []details.Data
	for _, room := range newRooms {
		if !existingRoomIDs[room.RoomID] {
			uniqueNewRooms = append(uniqueNewRooms, room)
			// Update our global map and slice with the new room
			existingRoomIDs[room.RoomID] = true
		}
	}

	// If no new unique rooms, return
	if len(uniqueNewRooms) == 0 {
		fmt.Println("No new rooms to add to rooms_details.csv")
		return nil
	}

	// Add unique new rooms to existing rooms
	existingRooms = append(existingRooms, uniqueNewRooms...)
	fmt.Printf("Adding %d new rooms to existing %d rooms\n", len(uniqueNewRooms), len(existingRooms)-len(uniqueNewRooms))

	// Write all rooms back to the CSV file

	// Create a backup of the original file if it exists
	if _, err := os.Stat(filePath); err == nil {
		backupFile := fmt.Sprintf("./output/rooms_details_backup_%s.csv", time.Now().Format("20060102_150405"))
		if err := copyFile(filePath, backupFile); err != nil {
			log.Printf("Error creating backup of rooms_details.csv: %v", err)
		} else {
			fmt.Printf("Created backup of rooms_details.csv at %s\n", backupFile)
		}
	}

	// Write the updated rooms to the CSV file
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating rooms_details.csv: %w", err)
	}
	defer outFile.Close()

	if err := gocsv.MarshalFile(&existingRooms, outFile); err != nil {
		return fmt.Errorf("error writing to rooms_details.csv: %w", err)
	}

	fmt.Printf("Successfully wrote %d rooms to rooms_details.csv\n", len(newRooms))
	return nil
}

// checkAndUpdateRoomDetails checks if a room exists and updates the CSV file with new room details
func checkAndUpdateRoomDetails(roomID int64) error {
	// Check if room exists
	exists, err := checkRoomExists(roomID)
	if err != nil {
		return fmt.Errorf("error checking if room exists: %w", err)
	}

	if exists {
		fmt.Printf("Room %d already exists in the CSV file\n", roomID)
		return nil
	}

	// Room doesn't exist, fetch details
	roomDetails, err := client.DetailsFromRoomID(roomID)
	if err != nil {
		return fmt.Errorf("error fetching details for room %d: %w", roomID, err)
	}

	// Update the CSV file with the new room
	if err := updateRoomDetailsCSV([]details.Data{roomDetails}); err != nil {
		return fmt.Errorf("error updating CSV file: %w", err)
	}

	fmt.Printf("Successfully added room %d to the CSV file\n", roomID)
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
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
	Cities := CITIES

	fmt.Printf("Starting asynchronous search for rooms in %d cities...\n", len(Cities))

	// Create channels for search results and synchronization
	citiesResultsChan := make(chan citySearchResult)
	var wgCities sync.WaitGroup

	// Set up a worker pool with a maximum of 5 concurrent searches
	maxConcurrentSearches := maxConcurrentSearch
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
	maxConcurrent := maxConcurrentRooms
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

			// Check if the room already exists in the CSV file
			exists, err := checkRoomExists(id)
			if err != nil {
				log.Printf("Error checking if room %d exists: %v", id, err)
			}

			if exists {
				// Room already exists, skip fetching details
				resultsChan <- roomDetailResult{
					roomID: id,
					name:   info.Name,
					err:    fmt.Errorf("room already exists in CSV file"),
				}
				return
			}

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
			// Check if this is our special case for existing rooms
			if result.err.Error() == "room already exists in CSV file" {
				fmt.Println("Skipped (already exists in CSV)")
			} else {
				fmt.Printf("Error: %s\n", result.err)
			}
			continue
		}

		fmt.Println("Success")
		allRoomDetails = append(allRoomDetails, result.details)
	}

	// Save all room details to a JSON file
	fmt.Printf("\nSaving details for %d rooms to rooms_details.csv\n", len(allRoomDetails))

	// Update the CSV file with the new rooms
	if err := updateRoomDetailsCSV(allRoomDetails); err != nil {
		log.Printf("Error updating rooms_details.csv: %v\n", err)
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
