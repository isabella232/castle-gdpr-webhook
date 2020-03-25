SRC := hook.go hook_test.go main.go
EXE := castle-gdpr-webhook
ZIP := function.zip
CWD=$(shell pwd)
TEST_AWS_ACCOUNT ?= DANGER-security
TEST_AWS_ACCOUNT_ID ?= 987056895854
PRODUCTION_AWS_ACCOUNT ?= DANGER-dw
PRODUCTION_AWS_ACCOUNT_ID ?= 873344020507
PRODUCTION_AWS_REGION ?= us-east-1
VERSION=1.0.0
S3_BUCKET="castle-gdpr-releases"
S3_KEY="castle-gdpr-webhook-${VERSION}.zip"

# this must match terrafrom
FUNCTION_NAME=CastleHandler 

default: ${EXE}

${EXE} : ${SRC}
	go test
	GOOS=linux go build
	zip function.zip ${EXE}

.PHONY: upload
upload: ${EXE}
	@echo "Uploading binary to account: ${PRODUCTION_AWS_ACCOUNT} bucket: ${S3_BUCKET} key: ${S3_KEY}"
	aws-okta exec ${PRODUCTION_AWS_ACCOUNT} -- aws s3 cp ${ZIP} s3://${S3_BUCKET}/${S3_KEY}

.PHONY: deploy-production
deploy-production: ${EXE}
	@echo "Deploying to ${PRODUCTION_AWS_ACCOUNT} account id ${PRODUCTION_AWS_ACCOUNT_ID}"
	aws-okta exec ${PRODUCTION_AWS_ACCOUNT} -- aws lambda update-function-code \
	       	--function-name ${EXE} \
		--s3-bucket ${S3_BUCKET} \
		--s3-key ${S3_KEY} \
		--region ${PRODUCTION_AWS_REGION}

.PHONY: deploy-test
deploy-test: ${EXE}
	@echo "Deploying to ${TEST_AWS_ACCOUNT} account id ${TEST_AWS_ACCOUNT_ID}"
	aws-okta exec ${TEST_AWS_ACCOUNT} -- aws lambda update-function-code \
	       	--function-name ${FUNCTION_NAME} \
		--zip-file fileb://${ZIP} \
		--region us-west-2

.PHONY: invoke-lambda 
invoke-lambda:
	@echo "Invoking lambda directly"
	aws-okta exec ${TEST_AWS_ACCOUNT} -- aws lambda invoke \
		--function-name ${EXE} \
		--invocation-type "RequestResponse" \
		--region us-west-2 \
		response.txt
	@echo "response.txt contains"
	@cat response.txt

# this only has to be done once
.PHONY: create-function
create-function: ${EXE}
	@echo "Creating function in ${PRODUCTION_AWS_ACCOUNT} account id ${PRODUCTION_AWS_ACCOUNT_ID}"
	aws-okta exec ${PRODUCTION_AWS_ACCOUNT} -- aws lambda create-function \
	       	--function-name ${EXE} \
		--runtime go1.x \
  		--zip-file fileb://${ZIP} \
	       	--handler ${EXE} \
  		--role arn:aws:iam::${PRODUCTION_AWS_ACCOUNT_ID}:role/lambda-castle-gdpr-webhook \
		--region ${PRODUCTION_AWS_REGION}

.PHONY: clean
clean:
	go clean
	rm -rf ${EXE} ${ZIP}
