package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	dynamock "github.com/gusaul/go-dynamock"
)

var mock *dynamock.DynaMock

var testUser = User{
	UserId:     "1234",
	Username: "test-username",
	Name:       "name",
	Email:      "test-email@domain.com",
}

func TestParseUserFromEvent(t *testing.T) {

	header := events.CognitoEventUserPoolsHeader{
		UserPoolID: "1234",
	}

	request := events.CognitoEventUserPoolsPostConfirmationRequest{
		UserAttributes: map[string]string{
			"username": testUser.Username,
			"name": testUser.Name,
			"email": testUser.Email,
		},
	}

	event := events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: header,
		Request: request,
	}	

	returnUser, _ := parseUserFromEvent(event)

	assert.Equal(t, testUser, *returnUser) 
}

func TestAddRatingToDB(t *testing.T) {

	var dyna DynamoAPI
	dyna.Db, mock = dynamock.New()


	putItem := map[string]*dynamodb.AttributeValue{
		"UserId":     {S: aws.String(testUser.UserId)},
		"Username": {S: aws.String(testUser.Username)},
		"Name":       {S: aws.String(testUser.Name)},
		"Email":      {S: aws.String(testUser.Email)},
	}

	mock.ExpectPutItem().ToTable(TableName).WithItems(putItem)

	_, _ = addUserToDB(dyna, testUser)
}