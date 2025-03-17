package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"gobnb"
	"gobnb/availability"
	"gobnb/utils"

	"github.com/gocarina/gocsv"
)

var maxConcurrentAvailability = 20

var proxyRotator *utils.ProxyRotator

var client gobnb.Client

// loadProxies loads proxy URLs from a file
func loadProxies(filePath string) ([]string, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Proxy file %s does not exist, continuing without proxies", filePath)
		return nil, nil
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading proxy file: %w", err)
	}

	// Split by newlines and filter empty lines
	var proxies []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			proxies = append(proxies, line)
		}
	}

	return proxies, nil
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	proxyFilePath := "./proxies.txt"
	proxyURLs, err := loadProxies(proxyFilePath)
	if err != nil {
		log.Printf("Warning: Failed to load proxies: %v", err)
	}

	// Initialize proxy rotator
	if len(proxyURLs) > 0 {
		proxyRotator, err = utils.NewProxyRotator(proxyURLs)
		if err != nil {
			log.Printf("Warning: Failed to initialize proxy rotator: %v", err)
		} else {
			log.Printf("Successfully loaded %d proxies", proxyRotator.Count())
		}
	} else {
		log.Println("No proxies loaded, requests will use direct connection")
	}

	client = gobnb.DefaultClient()
	// Call the function to generate availability data
	GenerateAvailabilityFromCSV()
}

// getAvailability fetches and displays availability for a specific room ID
func getAvailability(roomID int64) {
	fmt.Printf("Fetching availability for room ID: %d\n", roomID)

	// Fetch availability for the room
	availabilityData, daysData, err := availability.InputData{
		RoomId:     roomID,
		StartMonth: 3,
		StartYear:  2025,
	}.GetAvailabilityCalendar("USD", nil)
	if err != nil {
		log.Println("Error fetching availability:", err)
		return
	}

	fmt.Println("Number of availability days:", len(daysData))
	fmt.Println("Availability data:", availabilityData)

	// Create folder for the room if it doesn't exist
	folderPath := fmt.Sprintf("./output/%d", roomID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Save all availability data to a file
	if len(daysData) == 0 {
		return
	}
	allJSON, _ := gocsv.MarshalString(&daysData)
	filePath := fmt.Sprintf("%s/availability.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Println("Error saving availability data:", err)
		return
	}
	fmt.Printf("Availability data saved to %s\n", filePath)
}

// getRoomYearAvailability fetches availability data for a specific room ID for a year
func getRoomYearAvailability(roomID int64) {
	fmt.Printf("Fetching year availability for room ID: %d\n", roomID)

	// Fetch availability for the room
	availabilityData, daysData, err := availability.InputData{
		RoomId:     roomID,
		StartMonth: 3,
		StartYear:  2025,
	}.GetAvailabilityCalendar("USD", nil)
	if err != nil {
		log.Println("Error fetching availability:", err)
		return
	}

	// Print debug information about the response
	fmt.Println("Number of availability days:", len(daysData))
	fmt.Println("Availability data:", availabilityData)

	// Create folder for the room if it doesn't exist
	folderPath := fmt.Sprintf("./output/%d", roomID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Save availability data to a file
	allJSON, _ := gocsv.MarshalString(&daysData)
	filePath := fmt.Sprintf("%s/availability.csv", folderPath)
	if err := os.WriteFile(filePath, []byte(allJSON), 0644); err != nil {
		log.Println("Error saving availability data:", err)
		return
	}
	fmt.Printf("Availability data saved to %s\n", filePath)
}

// fetchAvailabilityForRoom fetches availability data for a specific room ID and saves it to a file
func fetchAvailabilityForRoom(roomID int64, folderPath string) {
	// Check if availability file already exists
	availabilityFilePath := filepath.Join(folderPath, "availability.csv")
	if _, err := os.Stat(availabilityFilePath); err == nil {
		// File already exists, skip
		fmt.Printf("Availability for room %d already exists, skipping\n", roomID)
		return
	}

	fmt.Printf("Fetching availability for room ID: %d\n", roomID)

	// Get a proxy for this request
	var proxy *url.URL
	if proxyRotator != nil {
		proxy = proxyRotator.GetNextProxy()
		if proxy != nil {
			log.Printf("Using proxy %s for reviews in %s", proxy.String())
		}
	}

	// Fetch availability data for the room
	_, daysData, err := availability.InputData{
		RoomId:     roomID,
		StartMonth: 3,
		StartYear:  2025,
	}.GetAvailabilityCalendar("USD", proxy)
	if err != nil {
		log.Printf("Error fetching availability for room %d: %v\n", roomID, err)
		return
	}

	fmt.Printf("Room %d: Found %d availability days\n", roomID, len(daysData))

	// Save availability data to a file
	allJSON, _ := gocsv.MarshalString(&daysData)
	if err := os.WriteFile(availabilityFilePath, []byte(allJSON), 0644); err != nil {
		log.Printf("Error saving availability for room %d: %v\n", roomID, err)
		return
	}
	fmt.Printf("Room %d: Availability data saved to %s\n", roomID, availabilityFilePath)
}

// GenerateAvailabilityFromCSV reads the room details CSV file and concurrently fetches availability for each room ID
func GenerateAvailabilityFromCSV() {
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
	maxConcurrent := maxConcurrentAvailability // Limit concurrent requests to avoid rate limiting
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

			// Fetch availability
			log.Printf("[%d/%d] Processing room ID: %d\n", index+1, len(records), id)
			fetchAvailabilityForRoom(id, folderPath)

			// Add a random delay between requests to avoid rate limiting
			// Random delay between 500ms and 2000ms
			minDelay := 500
			maxDelay := 2000
			randomDelay := minDelay + rand.Intn(maxDelay-minDelay)
			time.Sleep(time.Duration(randomDelay) * time.Millisecond)
		}(roomID, i)
	}

	// Wait for all goroutines to complete
	fmt.Println("Waiting for all availability fetching operations to complete...")
	wg.Wait()
	fmt.Println("All availability data has been fetched successfully!")
}
