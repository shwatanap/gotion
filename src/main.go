package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/shwatanap/gotion/src/oauth"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Unable to read .env file: %v", err)
	}

	client := oauth.NewClient()
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	events, err := srv.Events.List("primary").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}
