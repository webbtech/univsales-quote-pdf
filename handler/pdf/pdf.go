package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"

	log "github.com/sirupsen/logrus"
)

// Response data format
type Response struct {
	Code      int         `json:"code"`      // HTTP status code
	Data      interface{} `json:"data"`      // Data payload
	Message   string      `json:"message"`   // Error or status message
	Status    string      `json:"status"`    // Status code (error|fail|success)
	Timestamp int64       `json:"timestamp"` // Machine-readable UTC timestamp in nanoseconds since EPOCH
}

// HandleRequest function
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// var err error

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	t := time.Now()
	log.Errorf("req with log: %+v\n", req.Body)

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return gatewayResponse(Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs), nil
	}

	return gatewayResponse(Response{
		Code:      201,
		Data:      "post data",
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs), nil
}

func main() {
	lambda.Start(thundra.Wrap(HandleRequest))
}

// Your first Thundra API key, use it freely: b41775a3-4be7-4bbc-8106-ef649370ad63

func gatewayResponse(resp Response, hdrs map[string]string) events.APIGatewayProxyResponse {

	body, _ := json.Marshal(&resp)
	if resp.Status == "error" {
		log.Errorf("Error: status: %s, code: %d, message: %s", resp.Status, resp.Code, resp.Message)
	}

	return events.APIGatewayProxyResponse{Body: string(body), Headers: hdrs, StatusCode: resp.Code}
}
