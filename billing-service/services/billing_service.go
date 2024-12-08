package services

import (
	"billing-service/database"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MakMoinee/go-mith/pkg/email"
)

// CalculateAmount calculates the amount based on time duration (rate: $0.30 per minute)
func CalculateAmount(startTime, endTime time.Time) (float64, error) {
	// Calculate duration in minutes
	duration := endTime.Sub(startTime).Minutes()

	// Rate per minute
	const ratePerMinute = 0.30

	// Calculate amount
	amount := duration * ratePerMinute

	return amount, nil
}

// ApplyPromotion checks the promotion code and applies it to the billing
func ApplyPromotion(reservationID int, promotionCode string) (float64, error) {
	// Query the discount percentage from the database
	discountPercentage, err := database.GetDiscountPercentage(promotionCode)
	if err != nil {
		return 0, err
	}

	// If promotion doesn't exist, return an error
	if discountPercentage == 0 {
		return 0, errors.New("invalid promotion code")
	}

	// Update the billing's discount in the database
	err = database.UpdateBillingDiscount(reservationID, discountPercentage)
	if err != nil {
		return 0, err
	}

	// Return the discount percentage if everything is successful
	return discountPercentage, nil
}

func UpdateReservationAmount(revId int) error {
	_ = SendEmail(revId)
	return database.UpdateReservationAmount(revId)
}

func InsertBilling(revId int, amount float64) error {
	// Check if the billing record already exists
	exists, err := database.CheckBillingExists(revId)
	if err != nil {
		return fmt.Errorf("error checking billing existence: %v", err)
	}

	if exists {
		return fmt.Errorf("billing record already exists for reservation ID: %d", revId)
	}

	// Proceed to insert only if it doesn't exist
	return database.InsertBilling(revId, amount)
}

func SendEmail(revId int) error {
	// Fetch billing data from database
	billingData, err := database.GetReservationForEmail(revId)
	if err != nil {
		log.Printf("Error fetching billing data: %v", err)
		return err
	}

	// email message body
	emailBody := fmt.Sprintf("Hello,\n\nHere are your reservation billing details:\n\n"+"Initial Amount: $%.2f\n"+"Discount: $%.2f\n"+"Final Amount: $%.2f\n\nThank you for your reservation.",
		billingData.InitialAmount, billingData.Discount, billingData.FinalAmount,
	)

	emailService := email.NewEmailService(587, "smtp.gmail.com", "quahminghui@gmail.com", "lhfccuxmswnkegbt")

	isEmailSent, err := emailService.SendEmail("quahmingkid@gmail.com", "Reservation Details", emailBody)
	if err != nil {
		log.Fatalf("Error sending email: %s", err)
	}

	if isEmailSent {
		log.Println("Email Send Successfully")
	} else {
		log.Println("Failed to send email")
	}
	return nil
}
