package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/epsagon/epsagon-go/epsagon"
	"github.com/pulpfree/univsales-quote-pdf/config"
	"github.com/pulpfree/univsales-quote-pdf/model"
	"github.com/pulpfree/univsales-quote-pdf/model/mongo"
	"github.com/pulpfree/univsales-quote-pdf/pdf"
	"github.com/pulpfree/univsales-quote-pdf/pkgerrors"

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

var (
	cfg      *config.Config
	stdError *pkgerrors.StdError
)

func init() {
	cfg = &config.Config{}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleRequest function
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var err error

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	hdrs["Access-Control-Allow-Origin"] = "*"
	hdrs["Access-Control-Allow-Methods"] = "GET,OPTIONS,POST,PUT"
	hdrs["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{Body: string("null"), Headers: hdrs, StatusCode: 200}, nil
	}

	t := time.Now()

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return gatewayResponse(Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs, nil), nil
	}

	// Set and validate request params
	var r *pdf.Request
	json.Unmarshal([]byte(req.Body), &r)

	db, err := mongo.NewDB(cfg.GetMongoConnectURL(), cfg.DBName)
	if err != nil {
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	var q *model.Quote
	q, err = db.FetchQuote(r.QuoteID)
	if err != nil {
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	var p *pdf.PDF
	p = pdf.New(r, q, cfg)
	if r.DocType == "quote" {
		err = p.Quote()
	} else if r.DocType == "invoice" {
		err = p.Invoice()
	}
	if err != nil {
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	location, err := p.SaveToS3()
	if err != nil {
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}
	log.Infof("Successfully created PDF with location: %s", location)

	return gatewayResponse(Response{
		Code:      201,
		Data:      location,
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs, nil), nil
}

func main() {
	log.Println("enter main")
	lambda.Start(epsagon.WrapLambdaHandler(
		&epsagon.Config{ApplicationName: "univsales"},
		HandleRequest))
}

func gatewayResponse(resp Response, hdrs map[string]string, err error) events.APIGatewayProxyResponse {

	if err != nil {
		resp.Code = 500
		resp.Status = "error"
		log.Error(err)
		// send friendly error to client
		if ok := errors.As(err, &stdError); ok {
			resp.Message = stdError.Msg
		} else {
			resp.Message = err.Error()
		}
	}
	body, _ := json.Marshal(&resp)

	return events.APIGatewayProxyResponse{Body: string(body), Headers: hdrs, StatusCode: resp.Code}
}
