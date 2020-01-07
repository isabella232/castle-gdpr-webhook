package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var region = aws.String("us-west-2")
var bucket = aws.String("castle-gdpr-user-data")

type MyEvent struct {
	Name string `json:"name"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestName string `json:"name"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseName string `json:"name"`
}

// Handler function Using AWS Lambda Proxy Request
// from https://github.com/serverless/examples/blob/master/aws-golang-http-get-post/postFolder/postExample.go
func ServerlessCallback(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		RequestName: "",
	}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseName: bodyRequest.RequestName + " LastName",
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

// Handler function Using AWS Lambda Proxy Request
// https://github.com/serverless/examples/blob/master/aws-golang-http-get-post/getFolder/getExample.go
func ServerlessGetHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//Get the path parameter that was sent
	//name := request.PathParameters["name"]
	name := request.Path

	//Generate message that want to be sent as body
	message := fmt.Sprintf(" { \"Message\" : \"Hello %s \" } ", name)

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

func HandleCallback(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	message := fmt.Sprintf(" { \"Message\" : \"Not Yet Implemented\" } ")
	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

//var userRequest = regex.MustCompile(`/users/[^/]+/`)

// downloads a url to file
func DownloadFile(filepath, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func UploadFileToS3(bucket, filename, localfile string) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: region},
	)
	svc := s3.New(sess)

	file, err := os.Open(localfile)
	if err != nil {
		// handle error
		log.Printf("unable to open localfile: %s\n", localfile)
		return err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	log.Printf("Uploading: %s --> %s/%s size: %d\n", localfile, bucket, filename, size)

	// upload the file
	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(filename),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	_, err = svc.PutObject(params)
	if err != nil {
		log.Printf("error uploading: %s, error: %s\n", localfile, err.Error())
		return err
	}

	return nil
}

func HandleAllRequests(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("HandleAllRequests called with path: %s\n", request.Path)

	log.Printf("HandleAllRequests called with body: %s castleSignature: %s\n", request.Body, request.Headers["X-Castle-Signature"])

	sarDataUrl, userId, err := HandleIncomingWebHookData(request.Body, request.Headers["X-Castle-Signature"], "i'm a secret")
	if err != nil {
		fmt.Printf("HandleIncomingWebHookData err: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	tmpfile, err := ioutil.TempFile("/tmp", "castlegdpr."+userId+".*.zip")
	if err != nil {
		fmt.Printf("HandleIncomingWebHookData failed to make tempfile err: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	tmpfile.Close()
	name := tmpfile.Name()

	err = DownloadFile(name, sarDataUrl)
	if err != nil {
		fmt.Printf("HandleIncomingWebHookData failed to download sar data: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	// s3location

	err = UploadFileToS3(*bucket, userId+".zip", name)
	if err != nil {
		fmt.Printf("HandleIncomingWebHookData failed to upload sar data to s3: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleAllRequests)
}
