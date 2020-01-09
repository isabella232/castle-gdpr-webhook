package main

import (
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var region = aws.String("us-west-2")
var bucket = aws.String("castle-gdpr-user-data")
var keyname = "/hermes/prod/castle/api_secret"

// downloads a url to file
func DownloadFile(filepath, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

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

// reads the HMac secret, a secure string, from the ssm
func getHMacSecret() string {
	sess, err := session.NewSession(&aws.Config{
		Region: region},
	)
	if err != nil {
		log.Printf("getHMacSecret: failed to create session: %s\n", err.Error())
		return ""
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(*region))
	withDecryption := true
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		log.Printf("getHMacSecret: failed to read %s error: %s\n", keyname, err.Error())
		return ""
	}
	//log.Printf("param: %+v\n", param)

	value := *param.Parameter.Value
	return value
}

func HandleAllRequests(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("HandleAllRequests called with path: %s\n", request.Path)

	log.Printf("HandleAllRequests called with body: %s castleSignature: %s\n", request.Body, request.Headers["x-castle-signature"])
	log.Printf("HandleAllRequests request: %+v\n", request.Headers)

	saveRequestBody(request)

	signature := request.Headers["x-castle-signature"] // curl sets the headers this way
	if len(signature) == 0 {
		signature = request.Headers["X-Castle-Signature"] // how it comes from Castle
		if len(signature) == 0 {
			log.Printf("HandleIncomingWebHookData err: no x-castle-signature specified\n")
			return events.APIGatewayProxyResponse{Body: "", StatusCode: 500}, nil
		}
	}

	hmacSecret := getHMacSecret()

	sarDataUrl, userId, err := HandleIncomingWebHookData(request.Body, signature, hmacSecret)
	if err != nil {
		log.Printf("HandleIncomingWebHookData err: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	tmpfile, err := ioutil.TempFile("/tmp", "castlegdpr."+userId+".*.zip")
	if err != nil {
		log.Printf("HandleIncomingWebHookData failed to make tempfile err: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	tmpfile.Close()
	name := tmpfile.Name()

	err = DownloadFile(name, sarDataUrl)
	if err != nil {
		log.Printf("HandleIncomingWebHookData failed to download sar data: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	err = UploadFileToS3(*bucket, userId+".zip", name)
	if err != nil {
		log.Printf("HandleIncomingWebHookData failed to upload sar data to s3: %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleAllRequests)
}
