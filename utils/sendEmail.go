package utils

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Data struct {
	Name     string
	Sender   string
	Receiver string
	Subject  string
	Payload  interface{}
}

func SendEmail(data Data, templatePath string) (error) {
	log.Println("Sending email to: ", data.Receiver)
	
	var body bytes.Buffer
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	
	log.Println("Working directory: ", wd)
	t, err := template.ParseFiles(wd + templatePath)

	if err != nil {
		return err
	}

	// Execute the template with the provided data
	t.Execute(&body, struct{ Name string }{Name: data.Name})

 	log.Println("Attempting to send email body")
	// Construct the email
	m := gomail.NewMessage()
	m.SetHeader("From", data.Sender)
	m.SetHeader("To", data.Receiver)
	m.SetHeader("Subject", data.Subject)
	// invoiceLink := ""
	// Set the email body as HTML content
	m.SetBody("text/html", body.String())

	if payload, ok := data.Payload.(map[string]interface{}); ok {
		if invoiceLink, ok := payload["invoiceLink"].(string); ok && invoiceLink != "" {
			log.Println("Attaching invoice to email")
			// Example of downloading a file and attaching it
			response, err := http.Get(invoiceLink)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			// Create a temporary file to save the content
			tmpFile, err := os.Create("invoice.pdf")
			if err != nil {
				return err
			}
			
			defer tmpFile.Close()

			// Write the content from the response to the temporary file
			_, err = io.Copy(tmpFile, response.Body)
			if err != nil {
				return err
			}
			

			// Attach the downloaded file to the email
			m.Attach("invoice.pdf")
		}
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}


	log.Println("Attempting to send email")
	// Create the dialer and send the email
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
