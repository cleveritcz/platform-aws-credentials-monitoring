package main

import (
	"fmt"
	"log"
	"os"
	"time"
        "net/smtp"
	
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

var (
	username   string
	guard      string
)

func init() {
	var ok bool

	if username, ok = os.LookupEnv("HAZARD"); !ok || username == "" {
		log.Fatal("the HAZARD environment variable is not set")
	}

	if guard, ok = os.LookupEnv("GUARD"); !ok || guard == "" {
		log.Fatal("the GUARD environment variable is not set")
	}
}

type RestrictedData struct {
	userName string
	age      int
}

func Notification(message string) {

	// Sender data.
        from := "from@gmail.com"
        password := "<Email Password>"
	
        msg := []byte("From: aws@info.cz\r\n" +
            "To: martin.smola@centrum.cz\r\n" +
            "Subject: guard + "-" + message\r\n\r\n" +
            "Email body\r\n")
	
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
	
        smtpHost := "smtp.gmail.com"
        smtpPort := "587"	
        
}

func getAgeKey(username string) float64 {

	svc := iam.New(session.New())
	input := &iam.ListAccessKeysInput{
		UserName: aws.String(username),
	}
	result, err := svc.ListAccessKeys(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return 0
	}
	duration := time.Since(result.AccessKeyMetadata[0].CreateDate.UTC())
	//Calculate minutes old
	return duration.Minutes()
}

func alert(r RestrictedData) {

	var message string

	switch {
	case r.age < 12:
		message = "is under control"
	case r.age < 25:
		message = "is near to be out of control."
	case r.age >= 25:
		message = "is out of control. We must rotate the key!"
	case r.age >= 92:
		message = "is out of control. We must rotate the key!"
	}

	Notification(fmt.Sprintf("%s %s", r.userName, message))
}

func LambdaHandler() {
	hazardInfo := make([]RestrictedData, 0)

	restrictedData := RestrictedData{
		userName: username,
		age:      int(getAgeKey(username)),
	}

	hazardInfo = append(hazardInfo, restrictedData)
	for _, item := range hazardInfo {
		alert(item)
	}
}

func main() {
	lambda.Start(LambdaHandler)
}
