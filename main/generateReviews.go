package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"gobnb"
	"gobnb/reviews"

	"github.com/gocarina/gocsv"
)

var maxConcurrentReviews = 20

var client gobnb.Client

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	client = gobnb.DefaultClient()
	// Call the function to generate reviews
	GenerateReviewsFromCSV()
}

// getReviews fetches and displays reviews for a specific room ID
func getReviews(roomID int64) {
	fmt.Printf("Fetching reviews for room ID: %d\n", roomID)

	// Fetch reviews for the room
	reviewData, err := reviews.InputData{
		RoomId: roomID,
	}.GetAllReviewsOfRoom(roomID, "USD", nil)
	if err != nil {
		log.Println("Error fetching reviews:", err)
		return
	}

	fmt.Println(len(reviewData))

	// Create folder for the room if it doesn't exist
	folderPath := fmt.Sprintf("./output/%d", roomID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Save all paginated reviews to a file
	if len(reviewData) == 0 {
		return
	}
	allJSON, _ := gocsv.MarshalString(&reviewData)
	filePath := fmt.Sprintf("%s/reviews.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Println("Error saving all reviews:", err)
		return
	}
	fmt.Printf("All paginated reviews saved to %s\n", filePath)
}

// fetchReviewsForRoom fetches reviews for a specific room ID and saves them to a file
func fetchReviewsForRoom(roomID int64, folderPath string) {
	// Check if reviews file already exists
	reviewsFilePath := filepath.Join(folderPath, "reviews.csv")
	if _, err := os.Stat(reviewsFilePath); err == nil {
		// File already exists, skip
		fmt.Printf("Reviews for room %d already exist, skipping\n", roomID)
		return
	}

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
	if err := os.WriteFile(reviewsFilePath, []byte(allJSON), 0644); err != nil {
		log.Printf("Error saving reviews for room %d: %v\n", roomID, err)
		return
	}
	fmt.Printf("Room %d: Reviews saved to %s\n", roomID, reviewsFilePath)
}

// GenerateReviewsFromCSV reads the room details CSV file and concurrently fetches reviews for each room ID
func GenerateReviewsFromCSV() {
	// Path to the CSV file
	csvPath := filepath.Join("output", "rooms_details.csv")

	// Open the CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read CSV header: %v", err)
	}

	// Find the index of the roomID column
	roomIDIndex := -1
	for i, column := range header {
		if column == "RoomID" {
			roomIDIndex = i
			break
		}
	}

	if roomIDIndex == -1 {
		log.Fatalf("Could not find RoomID column in CSV")
	}

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV records: %v", err)
	}

	fmt.Printf("Found %d room records in CSV\n", len(records))

	// Set up concurrency control
	maxConcurrent := maxConcurrentReviews // Limit concurrent requests to avoid rate limiting
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	// We'll skip proxy rotation for this implementation to simplify
	// If needed, proxies can be added later

	// Process each room
	for i, record := range records {
		// Get the room ID
		roomIDStr := record[roomIDIndex]
		roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
		if err != nil {
			log.Printf("Warning: Failed to parse room ID '%s': %v\n", roomIDStr, err)
			continue
		}

		// Add to wait group
		wg.Add(1)

		// Acquire semaphore slot
		semaphore <- struct{}{}

		// Process room asynchronously
		go func(id int64, index int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore slot when done

			// Create output directory if it doesn't exist
			folderPath := filepath.Join("output", "rooms", fmt.Sprintf("%d", id))
			os.MkdirAll(folderPath, 0755)

			// Fetch reviews
			log.Printf("[%d/%d] Processing room ID: %d\n", index+1, len(records), id)
			fetchReviewsForRoom(id, folderPath)

			// Add a random delay between requests to avoid rate limiting
			// Random delay between 500ms and 2000ms
			minDelay := 500
			maxDelay := 2000
			randomDelay := minDelay + rand.Intn(maxDelay-minDelay)
			time.Sleep(time.Duration(randomDelay) * time.Millisecond)
		}(roomID, i)
	}

	// Wait for all goroutines to complete
	fmt.Println("Waiting for all review fetching operations to complete...")
	wg.Wait()
	fmt.Println("All reviews have been fetched successfully!")
}
