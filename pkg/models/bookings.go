package models

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorInvalidBookingData   = "invalid booking data"
	ErrorBookingAlreadyExists = "booking already exists"
	ErrorBookingDoesNotExist  = "booking does not exist"
)

type Booking struct {
	BookingID string `json:"bookingId"`
	UserID    string `json:"userId"`
	HotelID   string `json:"hotelId"`
	CheckIn   string `json:"checkIn"`
	CheckOut  string `json:"checkOut"`
}

func FetchBookings(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Booking, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	var bookings []Booking
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &bookings)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return &bookings, nil
}

func FetchBooking(bookingID, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Booking, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"bookingId": {
				S: aws.String(bookingID),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Booking)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func CreateBooking(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Booking, error) {
	var b Booking

	if err := json.Unmarshal([]byte(req.Body), &b); err != nil {
		return nil, errors.New(ErrorInvalidBookingData)
	}

	// Check if the booking already exists
	currentBooking, _ := FetchBooking(b.BookingID, tableName, dynaClient)
	if currentBooking != nil && len(currentBooking.BookingID) != 0 {
		return nil, errors.New(ErrorBookingAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(b)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &b, nil
}

func UpdateBooking(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Booking, error) {
	var b Booking
	if err := json.Unmarshal([]byte(req.Body), &b); err != nil {
		return nil, errors.New(ErrorInvalidBookingData)
	}

	// Check if the booking exists
	currentBooking, _ := FetchBooking(b.BookingID, tableName, dynaClient)
	if currentBooking == nil || len(currentBooking.BookingID) == 0 {
		return nil, errors.New(ErrorBookingDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(b)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &b, nil
}

func DeleteBooking(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	bookingID := req.QueryStringParameters["bookingId"]

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"bookingId": {
				S: aws.String(bookingID),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
