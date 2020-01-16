SRC := hook.go hook_test.go main.go
EXE := castle-gdpr-webhook
ZIP := function.zip
CWD=$(shell pwd)
TEST_AWS_ACCOUNT ?= DANGER-security
TEST_AWS_ACCOUNT_ID ?= 987056895854
PRODUCTION_AWS_ACCOUNT ?= DANGER-dw
PRODUCTION_AWS_ACCOUNT_ID ?= 873344020507

# this must match terrafrom
FUNCTION_NAME=CastleHandler 

default: ${EXE}

${EXE} : ${SRC}
	go test
	GOOS=linux go build

.PHONY: deploy-production
deploy-production: ${EXE}
	@echo "Deploying to ${PRODUCTION_AWS_ACCOUNT} account id ${PRODUCTION_AWS_ACCOUNT_ID}"
	zip function.zip ${EXE}
	aws-okta exec ${PRODUCTION_AWS_ACCOUNT} -- aws lambda update-function-code \
	       	--function-name ${EXE} \
  		--zip-file fileb://${ZIP} \
		--region us-west-2

.PHONY: deploy-test
deploy-test: ${EXE}
	@echo "Deploying to ${TEST_AWS_ACCOUNT} account id ${TEST_AWS_ACCOUNT_ID}"
	zip function.zip ${EXE}
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
	zip function.zip ${EXE}
	aws-okta exec ${PRODUCTION_AWS_ACCOUNT} -- aws lambda create-function \
	       	--function-name ${EXE} \
		--runtime go1.x \
  		--zip-file fileb://${ZIP} \
	       	--handler ${EXE} \
  		--role arn:aws:iam::${PRODUCTION_AWS_ACCOUNT_ID}:role/lambda-castle-gdpr-webhook \
		--region us-west-2

.PHONY: clean
clean:
	go clean
	rm -rf ${EXE} ${ZIP}
