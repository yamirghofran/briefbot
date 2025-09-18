package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/yamirghofran/briefbot/internal/services"
)

func main() {
	// Load .env file into the environment
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	emailService, err := services.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}

	// Create email request
	emailRequest := services.EmailRequest{
		ToAddresses:    []string{"yamirghofran@gmail.com"},
		Subject:        "What's shaking",
		HTMLBody:       "<h1>What's shaking?</h1>",
		TextBody:       "Hello there",
		ReplyToAddress: "yamirghofran@gmail.com",
	}

	if err := emailService.SendEmail(context.Background(), emailRequest); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	fmt.Println("Email sent successfully!")
}
