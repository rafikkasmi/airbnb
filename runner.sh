#!/bin/bash

echo "Running Listings Scraper..."
go run main/rooms.go

#sleep for 1 minute
sleep 60


echo "Running Reviews Scraper..."
go run main/generateReviews.go

#sleep for 1 minute
sleep 60


echo "Running Availability Scraper..."
go run main/generateAvailability.go


echo "DONE FOR TODAY GG !"

