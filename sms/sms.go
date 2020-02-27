package sms

import (
	"fmt"
	"log"
	"os"

	"github.com/AndroidStudyOpenSource/africastalking-go/sms"
	// "github.com/machariamuguku/golang-graphql-authentication/africastalking-go/sms"
)

// SendSms sends sms using africa's talking API
func SendSms(recipient, message string) {
	// Africa's talking credential
	username := os.Getenv("AFRICAS_TALKING_USERNAME") //Your Africa's Talking Username
	apiKey := os.Getenv("AFRICAS_TALKING_API_KEY")    //Production or Sandbox API Key
	env := "sandbox"                                  // Choose either Sandbox or Production

	//Call the Gateway, and pass the constants here!
	smsService := sms.NewService(username, apiKey, env)

	//Send SMS - REPLACE Recipient and Message with REAL Values
	//Leave ShortCode blank, "", if you don't have one)
	smsResponse, err := smsService.Send("", recipient, message)
	if err != nil {
		// log for the backend
		log.Printf("SendSms: error generating phone string: %v", err)
	}

	fmt.Println(smsResponse)
}
