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
func main(){
	// Uncomment one of these function calls to run different examples
	
	// Example 1: Search for rooms
	searchForRooms()
	
	// Example 2: Get room details
	// getRooms()
	
	// Example 3: Get reviews for a room
	// getReviews()
}

// searchForRooms searches for rooms in a specific location
func searchForRooms() {
	zoomvalue := 1
	checkIn := search.Check{}
	coords := search.CoordinatesInput{}
	
	results, err := search.InputData{
		Coordinates: coords,
		Check:       checkIn,
		ZoomValue:   zoomvalue,
		Query: "Marrakesh",
	}.SearchAll("USD", nil)
	if err != nil {
		log.Println(err)
		return
	}
	
	//make array unique by RoomId
	uniqueResults := make(map[int64]search.Data)
	for _, result := range results {
		uniqueResults[result.RoomID] = result
	}
	
	fmt.Println("Found", len(uniqueResults), "unique rooms")
	
	rawJSON, _ := json.MarshalIndent(uniqueResults, "", "  ")
	os.Remove("./searchResult.json")
	if err := os.WriteFile("./searchResult.json", rawJSON, 0644); err != nil {
		log.Println(err)
		return
	}
	
	fmt.Println("Search results saved to searchResult.json")
}


func getRooms(){
	var roomID int64
	roomID=290701
	// romID:=[]int{roomID}
	client := gobnb.DefaultClient()
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
	// Use a known valid room ID
	// Using a popular Airbnb listing
	var roomID int64 = 563002218
	
	fmt.Printf("Fetching reviews for room ID: %d\n", roomID)
	
	// Create a client
	client := gobnb.DefaultClient()
	
	// Fetch reviews for the room
	reviewData, err := client.ReviewsFromRoomID(roomID)
	if err != nil {
		log.Println("Error fetching reviews:", err)
		return
	}
	
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