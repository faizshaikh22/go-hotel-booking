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
	ErrorInvalidHotelData   = "invalid hotel data"
	ErrorHotelAlreadyExists = "hotel already exists"
	ErrorHotelDoesNotExist  = "hotel does not exist"
)

type Hotel struct {
	HotelID        string `json:"hotelId"`
	Name           string `json:"name"`
	City           string `json:"city"`
	TotalRooms     int    `json:"totalRooms"`
	AvailableRooms int    `json:"availableRooms"`
}

func FetchHotels(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Hotel, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	var hotels []Hotel
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &hotels)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return &hotels, nil
}

func FetchHotel(hotelID, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Hotel, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"hotelId": {
				S: aws.String(hotelID),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Hotel)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func CreateHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Hotel, error) {
	var h Hotel

	if err := json.Unmarshal([]byte(req.Body), &h); err != nil {
		return nil, errors.New(ErrorInvalidHotelData)
	}

	// Check if the hotel already exists
	currentHotel, _ := FetchHotel(h.HotelID, tableName, dynaClient)
	if currentHotel != nil && len(currentHotel.HotelID) != 0 {
		return nil, errors.New(ErrorHotelAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(h)
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
	return &h, nil
}

func UpdateHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Hotel, error) {
	var h Hotel
	if err := json.Unmarshal([]byte(req.Body), &h); err != nil {
		return nil, errors.New(ErrorInvalidHotelData)
	}

	// Check if the hotel exists
	currentHotel, _ := FetchHotel(h.HotelID, tableName, dynaClient)
	if currentHotel == nil || len(currentHotel.HotelID) == 0 {
		return nil, errors.New(ErrorHotelDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(h)
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
	return &h, nil
}

func DeleteHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	hotelID := req.QueryStringParameters["hotelId"]

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"hotelId": {
				S: aws.String(hotelID),
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
