package sms

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/AndroidStudyOpenSource/africastalking-go/sms"
)

var err = godotenv.Load()

// Africa's talking credentials
var (
	username = os.Getenv("AFRICAS_TALKING_USERNAME") //Your Africa's Talking Username
	apiKey   = os.Getenv("AFRICAS_TALKING_API_KEY")  //Production or Sandbox API Key
	env      = os.Getenv("ENVIRONMENT")              // Either sandbox or production
)

// SendSms sends sms using africa's talking API
func SendSms(recipient, message string) {

	//Call the Gateway, and pass the constants here!
	smsService := sms.NewService(username, apiKey, env)

	//Send SMS - REPLACE Recipient and Message with REAL Values
	//Leave ShortCode blank, "", if you don't have one)
	smsResponse, err := smsService.Send("", recipient, message)

	if err != nil {
		log.Printf("SendSms: sms send unsuccessful: %v", err)
	}

	fmt.Println(smsResponse)

}
