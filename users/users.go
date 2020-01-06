package main

import (
	"encoding/base64"
	//"encoding/json"
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"net/http"
	"os"
)

var region = aws.String("us-west-2")
var bucket = aws.String("castle-gdpr-user-data")

// sends a request for castle for the GDPR information for the specified uniqueId
func requestGdprInfoFromCastle(uniqueId string) {
	if len(uniqueId) == 0 {
		log.Printf("requestGdprInfoFromCastle called with empty string")
		return
	}

	requestUrl := "https://api.castle.io/v1/privacy/users/"
	req, err := http.NewRequest("POST", requestUrl, nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth("", "secretsauce")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Requesting info for %s from castle\n", uniqueId)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Castle response: %v\n", resp)
	defer resp.Body.Close()

}

// downloads the file from S3. returns the bytes in the file or an error
func download(filename string) ([]byte, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: region},
	)

	downloader := s3manager.NewDownloader(sess)
	buff := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    aws.String(filename),
		})
	if err != nil {
		log.Printf("Unable to download item %q, %v", filename, err)
		return nil, err
	}

	log.Printf("Downloaded %d bytes\n", numBytes)

	return buff.Bytes(), nil
}

func doess3FileExist(filename string) bool {
	rv := true

	sess, _ := session.NewSession(&aws.Config{
		Region: region},
	)

	downloader := s3manager.NewDownloader(sess)
	buff := &aws.WriteAtBuffer{}

	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    aws.String(filename),
		})
	if err != nil {
		// if any errors are reported the file is not able to be downloaded
		rv = false
	}
	return rv
}

// deletes the file from s3
func deleteFile(bucket, filename string) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: region},
	)
	svc := s3.New(sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}

	_, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}
	log.Printf("deleted file: %s\n", filename)
	return err
}

// downloads a url to file
func DownloadFile(filepath, url string) (string, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return "", err
}

// uploads a file s3
func UploadFileToS3(bucket, filename, localfile string) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: region},
	)
	svc := s3.New(sess)

	file, err := os.Open(localfile)
	if err != nil {
		// handle error
		log.Printf("unable to open localfile: %s\n", localfile)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

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
	}

	return err
}

// handles the content posted by castle, typically to /callback
func HandleCallback(webhookContent string) {

	if len(webhookContent) == 0 {
		log.Printf("HandleCallback called with no content, exiting.")
		return
	}

	//	json.Unmarshal

	// download the file
	//filename := path.Base(r.URL.Path)
	//tmpfile, err := ioutil.TempFile("/tmp", "castlegdpr.*.zip")
	//err := DownloadFile(tmpfile, "https://miro.medium.com/max/5472/1*iyIXvY2VZ2ariY9k2z1dzw.jpeg")

	// upload the file
}

// handles the request for gdpr data
func HandleUserRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	uniqueId := request.QueryStringParameters["unique_id"]
	if len(uniqueId) == 0 {
		log.Printf("HandleUserRequest called with no unique_id parameter\n")
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: 400}, nil
	}

	filename := uniqueId + ".zip"

	if doess3FileExist(filename) == true {
		log.Printf("%s does exists returning data\n", filename)

		file, err := download(filename)
		if err != nil {
			log.Printf("Error downloading filename: %s from s3: %v\n", filename, err)
			return events.APIGatewayProxyResponse{Body: "Internal Error", StatusCode: 500}, nil
		}

		// base64 encode it
		encodedFile := base64.StdEncoding.EncodeToString(file)

		deleteFile(*bucket, filename)

		// return the file
		return events.APIGatewayProxyResponse{Body: encodedFile, StatusCode: 200}, nil
	} else {
		log.Printf("%s does not exist requesting data from castle\n", filename)
		requestGdprInfoFromCastle(uniqueId)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
	}
}

func HandleAllRequests(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("HandleAllRequests called with path: %s\n", request.Path)

	// TODO authenticate the requests

	// very complex url routing
	if request.Path == "/users" {
		return HandleUserRequest(request)
	} else if request.Path == "/" {
		HandleCallback(request.Body)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
	} else {
		log.Printf("called with unknown path: %s\n", request.Path)
	}
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleUserRequest)
}
