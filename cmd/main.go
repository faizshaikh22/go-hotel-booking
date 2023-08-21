//GOOS=linux GOARCH=amd64 go build -o main
//use the above command to build the executable for aws lambda in gitbash

package main

import (
	"awslambda/pkg/handlers"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const (
	userTableName    = "user"
	hotelTableName   = "hotel"
	bookingTableName = "booking"
)

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.Resource {
	case "/user":
		return handleUserRequests(req)
	case "/hotel":
		return handleHotelRequests(req)
	case "/booking":
		return handleBookingRequests(req)
	default:
		return handlers.UnhandledMethod()
	}
}

func handleUserRequests(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, userTableName, dynaClient)
	case "POST":
		return handlers.CreateUser(req, userTableName, dynaClient)
	case "PUT":
		return handlers.UpdateUser(req, userTableName, dynaClient)
	case "DELETE":
		return handlers.DeleteUser(req, userTableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}

func handleHotelRequests(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		if req.Path == "/hotel/available" {
			return handlers.GetAvailableHotels(req, hotelTableName, dynaClient)
		}
		return handlers.GetHotel(req, hotelTableName, dynaClient)
	case "POST":
		return handlers.CreateHotel(req, hotelTableName, dynaClient)
	case "PUT":
		return handlers.UpdateHotel(req, hotelTableName, dynaClient)
	case "DELETE":
		return handlers.DeleteHotel(req, hotelTableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}

func handleBookingRequests(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetBooking(req, bookingTableName, dynaClient)
	case "POST":
		return handlers.CreateBooking(req, bookingTableName, dynaClient)
	case "PUT":
		return handlers.UpdateBooking(req, bookingTableName, dynaClient)
	case "DELETE":
		return handlers.DeleteBooking(req, bookingTableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}
