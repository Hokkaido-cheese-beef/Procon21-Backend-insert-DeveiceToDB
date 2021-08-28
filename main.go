package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type Request struct {
	SensorId  string  `json:"sensorID"`
	TimeStamp int     `json:"timestamp"`
	Co2       int     `json:"co2"`
	Temp      float64 `json:"temp"`
	Hum       float64 `json:"hum"`
}

type InsertDB struct {
	SensorId  string  `dynamodbav:"sensorID" json:"sensorID"`
	TimeStamp int     `dynamodbav:"timestamp" json:"timestamp"`
	Co2       int     `dynamodbav:"co2" json:"co2"`
	Temp      float64 `dynamodbav:"temp" json:"temp"`
	Hum       float64 `dynamodbav:"hum" json:"hum"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// DB接続
	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-northeast-1"))

	log.Println(request.Body)
	item := Request{}
	if err := json.Unmarshal(([]byte)(request.Body), &item); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	res := InsertDB{}
	log.Println(item)
	res = InsertDB{
		SensorId:  item.SensorId,
		TimeStamp: item.TimeStamp,
		Co2:       item.Co2,
		Temp:      item.Temp,
		Hum:       item.Hum,
	}

	insertData, err := dynamodbattribute.MarshalMap(res)
	if err != nil {
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("SensorData"),
		Item:      insertData,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	jsonData, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "origin,Accept,Authorization,Content-Type",
			"Content-Type":                 "application/json",
		},
		Body:       string(jsonData),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(handler)
}
