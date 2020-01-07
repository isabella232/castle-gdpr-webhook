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
var secret = "secret"

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

// uploads the specified localfile to the filename in the S3 bucket
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

// save the request body so that it can be inspected, this is mainly for debugging
func saveRequestBody(request events.APIGatewayProxyRequest) {
	tmpfile, err := ioutil.TempFile("/tmp", "request.*.tmp")
	if err != nil {
		log.Printf("saveRequest failed to make tempfile err: %s\n", err.Error())
	}
	tmpfile.Write([]byte(request.Body))
	tmpfile.Close()
	name := tmpfile.Name()

	log.Printf("Wrote saved request to: %s\n", name)

	err = UploadFileToS3(*bucket, "request.tmp", name)
	if err != nil {
		log.Printf("saveRequest failed to request Body to s3: %s\n", err.Error())
	}
}

func HandleAllRequests(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("HandleAllRequests called with path: %s\n", request.Path)

	log.Printf("HandleAllRequests called with body: %s castleSignature: %s\n", request.Body, request.Headers["x-castle-signature"])
	log.Printf("HandleAllRequests request: %+v\n", request.Headers)

	saveRequestBody(request)

	// in golang all headers are lowercase
	sarDataUrl, userId, err := HandleIncomingWebHookData(request.Body, request.Headers["x-castle-signature"], secret)
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
