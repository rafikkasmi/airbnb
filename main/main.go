package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gobnb/search"
	"gobnb"
	// "github.com/johnbalvin/gobnb/utils"
	// "github.com/johnbalvin/gobnb/details"
)
var client gobnb.Client
func main(){
	client = gobnb.DefaultClient()
	// Uncomment one of these function calls to run different examples
	
	// Example 1: Search for rooms
	searchForRooms()
	
	// Example 2: Get room details
	// getRooms()
	
	// Example 3: Get reviews for a room
	// getReviews()
}

// searchForRooms searches for rooms in a specific location and gets their details
func searchForRooms() {
	fmt.Println("Searching for rooms in New York City...")

	zoomvalue := 1
	checkIn := search.Check{}
	coords := search.CoordinatesInput{}
	
	results, err := search.InputData{
		Coordinates: coords,
		Check:       checkIn,
		ZoomValue:   zoomvalue,
		Query: "New York City, New York, United States",
	}.SearchAll("USD", nil)
	if err != nil {
		log.Println("Search error:", err)
		return
	}
	
	//make array unique by RoomId
	uniqueResults := make(map[int64]search.Data)
	for _, result := range results {
		uniqueResults[result.RoomID] = result
	}
	
	fmt.Println("Found", len(uniqueResults), "unique rooms")
	
	// Save search results
	rawJSON, _ := json.MarshalIndent(uniqueResults, "", "  ")
	os.Remove("./searchResult.json")
	if err := os.WriteFile("./searchResult.json", rawJSON, 0644); err != nil {
		log.Println(err)
		return
	}
	
	fmt.Println("Search results saved to searchResult.json")
	
	// Now get details for each room
	fmt.Println("\nFetching details for each room...")
	allRoomDetails := make([]interface{}, 0, len(uniqueResults))
	
	i := 0
	total := len(uniqueResults)
	for roomID, roomInfo := range uniqueResults {
		i++
		fmt.Printf("[%d/%d] Getting details for room %d: %s\n", i, total, roomID, roomInfo.Name)
		
		// Get room details
		roomDetails, err := client.DetailsFromRoomID(roomID)
		if err != nil {
			log.Printf("Error getting details for room %d: %s\n", roomID, err)
			continue
		}
		
		// Add to our collection
		allRoomDetails = append(allRoomDetails, roomDetails)
		
		// Add a small delay to avoid rate limiting
		// time.Sleep(500 * time.Millisecond)
	}
	
	// Save all room details to a JSON file
	fmt.Printf("\nSaving details for %d rooms to rooms_details.json\n", len(allRoomDetails))
	
	// Marshal the room details to JSON
	detailsJSON, err := json.MarshalIndent(allRoomDetails, "", "  ")
	if err != nil {
		log.Println("Error marshaling room details to JSON:", err)
		return
	}
	
	// Write the JSON to a file
	os.Remove("./rooms_details.json")
	if err := os.WriteFile("./rooms_details.json", detailsJSON, 0644); err != nil {
		log.Println("Error writing room details to file:", err)
		return
	}
	
	fmt.Println("Room details saved to rooms_details.json")
}

func getRoomDetail(roomId int64){
	

}

func getRooms(){
	var roomID int64
	roomID=290701
	// romID:=[]int{roomID}
	data, err := client.DetailsFromRoomID(roomID)
	if err != nil {
		log.Println("test:2 -> err: ", err)
		return
	}
	rawJSON, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s", rawJSON) //in case you don't have write permisions
	if err := os.WriteFile("./details.json", rawJSON, 0644); err != nil {
	   log.Println(err)
	   return
	}
}

// getReviews fetches and displays reviews for a specific room ID
func getReviews() {
	// Use a known valid room ID with reviews
	// Using the room ID from the example URL
	var roomID int64 = 290701
	
	fmt.Printf("Fetching reviews for room ID: %d\n", roomID)
	
	// Create a client
	client := gobnb.DefaultClient()
	
	// Fetch reviews for the room
	reviewData, err := client.ReviewsFromRoomID(roomID)
	if err != nil {
		log.Println("Error fetching reviews:", err)
		return
	}
	
	// Print debug information about the response
	fmt.Printf("Raw review data: %+v\n", reviewData)
	
	// Print summary information
	fmt.Printf("Room %d has %d reviews with an average rating of %.1f\n", 
		roomID, reviewData.TotalReviews, reviewData.Rating)
	
	// Print the first few reviews
	fmt.Println("Sample reviews:")
	for i, review := range reviewData.Reviews {
		if i >= 3 { // Only show first 3 reviews
			break
		}
		fmt.Printf("  Review #%d by %s (%d stars): %s\n", 
			i+1, review.AuthorName, review.Rating, truncateString(review.Comments, 100))
	}
	
	// Save all reviews to a file
	rawJSON, _ := json.MarshalIndent(reviewData, "", "  ")
	if err := os.WriteFile("./reviews.json", rawJSON, 0644); err != nil {
		log.Println("Error saving reviews:", err)
		return
	}
	fmt.Println("All reviews saved to reviews.json")
	
	// Fetch all reviews with pagination if there are more
	if reviewData.HasMoreReviews {
		fmt.Println("Fetching all reviews with pagination...")
		allReviews, err := client.AllReviewsFromRoomID(roomID)
		if err != nil {
			log.Println("Error fetching all reviews:", err)
			return
		}
		
		fmt.Printf("Successfully fetched all %d reviews\n", len(allReviews.Reviews))
		
		// Save all paginated reviews to a file
		allJSON, _ := json.MarshalIndent(allReviews, "", "  ")
		if err := os.WriteFile("./all_reviews.json", allJSON, 0644); err != nil {
			log.Println("Error saving all reviews:", err)
			return
		}
		fmt.Println("All paginated reviews saved to all_reviews.json")
	}
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}