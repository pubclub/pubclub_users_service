package main

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var TableName string = os.Getenv("TABLE_NAME")

type User struct {
	UserId   string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type DynamoAPI struct {
	Db dynamodbiface.DynamoDBAPI
}

func parseUserFromEvent(event events.CognitoEventUserPoolsPostConfirmation) (*User, error) {
	if len(event.Request.UserAttributes) == 0 {
		log.Fatal("Event item empty!")
		return nil, errors.New("Empty event item")
	}
	user := User{
		UserId:   event.CognitoEventUserPoolsHeader.UserPoolID,
		Username: event.Request.UserAttributes["username"],
		Name:     event.Request.UserAttributes["name"],
		Email:    event.Request.UserAttributes["email"],
	}
	return &user, nil
}

func addUserToDB(dyna DynamoAPI, user User) (*dynamodb.PutItemOutput, error) {

	putItem := map[string]*dynamodb.AttributeValue{
		"UserId":   {S: aws.String(user.UserId)},
		"Username": {S: aws.String(user.Username)},
		"Name":     {S: aws.String(user.Name)},
		"Email":    {S: aws.String(user.Email)},
	}

	input := &dynamodb.PutItemInput{
		Item:      putItem,
		TableName: aws.String(TableName),
	}

	output, err := dyna.Db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return nil, err
	}
	return output, nil
}

func HandleRequest(event events.CognitoEventUserPoolsPostConfirmation) (*dynamodb.PutItemOutput, error) {
	user, err := parseUserFromEvent(event)
	if err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	svc := dynamodb.New(sess)
	var dyna DynamoAPI
	dyna.Db = svc

	output, err := addUserToDB(dyna, *user)
	if err != nil {
		panic(err)
	}

	return output, nil
}

func main() {
	lambda.Start(HandleRequest)
}
