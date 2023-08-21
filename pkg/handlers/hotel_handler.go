package handlers

import (
	"awslambda/pkg/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func GetHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	hotelID := req.QueryStringParameters["hotelId"]
	if len(hotelID) > 0 {
		result, err := models.FetchHotel(hotelID, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{
				aws.String(err.Error()),
			})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := models.FetchHotels(tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

func CreateHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := models.CreateHotel(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

func UpdateHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := models.UpdateHotel(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteHotel(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := models.DeleteHotel(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, nil)
}

func GetAvailableHotels(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	// Parse request parameters
	selectedDateStr := req.QueryStringParameters["date"]
	city := req.QueryStringParameters["city"]

	// Convert selected date to time.Time
	_, err := time.Parse("2006-01-02", selectedDateStr)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMsg: aws.String("invalid date format"),
		})
	}

	// Fetch all hotels
	hotels, err := models.FetchHotels(tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusInternalServerError, ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	// Filter available hotels based on selected date and city
	availableHotels := make([]models.Hotel, 0)
	for _, hotel := range *hotels {
		if hotel.City == city && hotel.AvailableRooms > 0 {
			availableHotels = append(availableHotels, hotel)
		}
	}

	// Return the list of available hotels
	responseBody, err := json.Marshal(availableHotels)
	if err != nil {
		return apiResponse(http.StatusInternalServerError, ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, responseBody)
}
