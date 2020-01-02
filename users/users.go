package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"strings"
)

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

func printMap(m map[string]string) {
	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}

	for k, v := range m {
		log.Println(k + ": " + strings.Repeat(" ", maxLenKey-len(k)) + v)
	}
}

func HandleUserRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Got request for path: %s\n", request.Path, request.PathParameters)
	log.Printf("Path paramters: %#v\n", request.PathParameters)
	log.Printf("Querystring paramters: %#v\n", request.QueryStringParameters)
	log.Printf("Name paramter: %s\n", request.QueryStringParameters["name"])
	// pretty print the fucking map
	printMap(request.QueryStringParameters)

	if request.Path[0] != '/' {
		log.Printf("Invalid request path: %s", request.Path)
		return events.APIGatewayProxyResponse{Body: "invalid request path", StatusCode: 400}, nil
	}

	if S3FileExists(name) == true {
		log.Printf("%s does exists returning data\n", name)

	} else {
		log.Printf("%s does not exist requesting data from castle\n", name)

		file := DownloadFile(name)
		// base64 encode it
		encodedFile := b64.StdEncoding.EncodeToString(file)

		// Delete the file
		DeleteFile(name)

		// return the file
		return encodedFile
	}
	filename := strings.SplitAfter(request.Path, "/")[1] + ".txt"
	log.Printf("Attempting to download filename: %s from s3", filename)

	// download the file and base64 encode it, then set it to the body.
	// https://stackoverflow.com/questions/35804042/aws-api-gateway-and-lambda-to-return-image
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	downloader := s3manager.NewDownloader(sess)
	buff := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String("castle-gdpr-user-data"),
			Key:    aws.String("test.txt"),
		})
	if err != nil {
		log.Printf("Unable to download item %q, %v", "test.txt", err)
		return events.APIGatewayProxyResponse{Body: "unable to download file", StatusCode: 400}, nil
	}

	log.Printf("Downloaded %s bytes\n", numBytes)

	data := buff.Bytes() // now data is my []byte array

	// delete the file

	return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 200}, nil
}

func HandleAllRequests(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//name := request.PathParameters["name"]
	switch path := request.Path; path {
	case "/callback":
		return HandleCallback(request)
	default:
		log.Printf("No such route: %s", path)
		message := fmt.Sprintf("{ \"Message\" : \"No such route %s\" }", path)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, nil
	}
	return events.APIGatewayProxyResponse{Body: "...", StatusCode: 404}, nil
}

func HandleCallbackOrg(ctx context.Context, name MyEvent) (MyResponse, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	log.Print(lc.Identity.CognitoIdentityPoolID)
	log.Print(lc)
	return MyResponse{Message: fmt.Sprintf("Hello %s, context: %+v!", name.Name, lc)}, nil
	//return MyResponse{Message: fmt.Sprintf("Weewaa %s!", name.Name)}, nil

}

/*

func HandleRequest(ctx context.Context, name MyEvent) (MyResponse, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	log.Print(lc.Identity.CognitoIdentityPoolID)
	log.Print(lc)
	//return MyResponse{Message: fmt.Sprintf("Hello %s, %+v!", name.Name, lc)}, nil
	return MyResponse{Message: fmt.Sprintf("Weewaa %s!", name.Name)}, nil}
}
*/

func main() {
	//lambda.Start(HandleRequest)
	//lambda.Start(HandleCallback)
	//lambda.Start(ServerlessGetHandler)
	//lambda.Start(ServerlessCallback)
	//lambda.Start(HandleAllRequests)
	lambda.Start(HandleUserRequest)
}
