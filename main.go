package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Movie Entity
type Movie struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func findAll(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	size, err := strconv.Atoi(request.Headers["Count"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Count Header should be a number",
		}, nil
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		// TableName: aws.String(os.Getenv("TABLE_NAME")),
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Limit:     aws.Int64(int64(size)),
	}

	result, err := db.Scan(params)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning DynamoDB",
		}, nil
	}

	var movies []Movie

	dynamodbattribute.UnmarshalListOfMaps(result.Items, &movies)

	response, err := json.Marshal(&movies)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil

}

func main() {
	lambda.Start(findAll)
}
