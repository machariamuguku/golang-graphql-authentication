package emails

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail : function to send emails using sendgrid
func SendEmail(receiver string, subject string, emailContent string) {

	// get sendgrid api from .env
	sendgridAPI := os.Getenv("SENDGRID_API_KEY")

	// if it returns an empty key
	if sendgridAPI == "" {
		// log for the backend
		log.Println("sendEmail: SENDGRID_API_KEY from .env returned empty")
	}

	// get email address from .env
	emailFrom := os.Getenv("EMAIL_ADDRESS")

	// if it returns an empty key
	if emailFrom == "" {
		// log for the backend
		log.Println("sendEmail: EMAIL_ADDRESS from .env returned empty")
	}

	// sent from
	from := mail.NewEmail("www.muguku.co.ke", emailFrom)
	// sent to
	to := mail.NewEmail("user", receiver)

	// convert content to mail content type
	content := mail.NewContent("text/html", emailContent)

	// compose the email
	message := mail.NewV3MailInit(from, subject, to, content)

	client := sendgrid.NewSendClient(sendgridAPI)
	// try to send the email
	response, err := client.Send(message)
	// if error
	if err != nil {
		// log the error for backend
		log.Printf("sendEmail: email sending error %v", err)
	} else {
		// check response error codes
		// status codes 200 and 201 are success
		if !(response.StatusCode == 200 || response.StatusCode == 202) {

			// log the error for backend
			log.Printf("sendEmail: sendgrid email sending error %v", err)

		}

	}

}
