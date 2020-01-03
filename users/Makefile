SRC := users.go
EXE := castle-gdpr-users
NAME := castle-gdpr-users
CWD=$(shell pwd)

${NAME} : ${SRC}
	GOOS=linux go build -o ${NAME}

deploy: ${NAME}
	zip function.zip ${NAME}
	aws-okta exec DANGER-security -- aws lambda update-function-code \
	       	--function-name ${NAME} \
  		--zip-file fileb://function.zip \
		--region us-west-2

test:
	aws-okta exec DANGER-security -- aws lambda invoke \
		--function-name ${NAME} \
		--invocation-type "RequestResponse" \
		--region us-west-2 \
		response.txt
	cat response.txt

copytest:
	openssl rand -base64 16 > test.zip
	aws-okta exec DANGER-security -- aws s3 cp test.zip s3://castle-gdpr-user-data/ --region us-west-2

# this only has to be done once, note that the role must have been created in
# IAM, the role should have AWSLambdaBasicExecutionRole
create-function: ${NAME}
	zip function.zip ${NAME}
	aws-okta exec DANGER-security -- aws lambda create-function \
	       	--function-name ${NAME} \
		--runtime go1.x \
  		--zip-file fileb://function.zip \
	       	--handler ${NAME} \
  		--role arn:aws:iam::987056895854:role/lambda-castle-gdpr-users \
		--region us-west-2
clean:
	rm -rf ${NAME}
