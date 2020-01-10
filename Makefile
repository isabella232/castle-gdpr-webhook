SRC := hook.go hook_test.go main.go
EXE := castle-gdpr-webhook
ZIP := function.zip
CWD=$(shell pwd)
AWS_ACCOUNT := DANGER-dw

${EXE} : ${SRC}
	go test
	GOOS=linux go build

deploy: ${EXE}
	zip function.zip ${EXE}
	aws-okta exec ${AWS_ACCOUNT} -- aws lambda update-function-code \
	       	--function-name ${EXE} \
  		--zip-file fileb://${ZIP} \
		--region us-west-2

test:
	aws-okta exec ${AWS_ACCOUNT} -- aws lambda invoke \
		--function-name ${EXE} \
		--invocation-type "RequestResponse" \
		--region us-west-2 \
		response.txt
	cat response.txt

# this only has to be done once
create-function: ${EXE}
	zip function.zip ${EXE}
	aws-okta exec ${AWS_ACCOUNT} -- aws lambda create-function \
	       	--function-name ${EXE} \
		--runtime go1.x \
  		--zip-file fileb://${ZIP} \
	       	--handler ${EXE} \
  		--role arn:aws:iam::873344020507:role/lambda-castle-gdpr-webhook \
		--region us-west-2
clean:
	rm -rf ${EXE} ${ZIP}
