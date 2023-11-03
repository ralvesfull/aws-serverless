package main

import {
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"	
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
}

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}


func InsertUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Insert user to database

	var user User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body: err.Error(),
		}, nil
	}

	user.ID = uuid.New().String()

	// Insert user to database
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("UsersGo"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(user.ID),
			},
			"name": {
				S: aws.String(user.Name),
			},
			"email": {
				S: aws.String(user.Email),
			},
		},
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body: err.Error(),
		}, nil
	}

	body, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body: err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: string(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	fmt.Println("Start Lambda");
	lambda.Start(InsertUser)
}