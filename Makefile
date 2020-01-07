SRC := hook.go hook_test.go main.go
EXE := castle-gdpr-webhook
CWD=$(shell pwd)

${EXE} : ${SRC}
	GOOS=linux go build

deploy: ${EXE}
	zip function.zip ${EXE}
	aws-okta exec DANGER-security -- aws lambda update-function-code \
	       	--function-name ${EXE} \
  		--zip-file fileb://function.zip \
		--region us-west-2

test:
	aws-okta exec DANGER-security -- aws lambda invoke \
		--function-name ${EXE} \
		--invocation-type "RequestResponse" \
		--region us-west-2 \
		response.txt
	cat response.txt

# this only has to be done once
create-function: ${EXE}
	zip function.zip ${EXE}
	aws-okta exec DANGER-security -- aws lambda create-function \
	       	--function-name ${EXE} \
		--runtime go1.x \
  		--zip-file fileb://function.zip \
	       	--handler ${EXE} \
  		--role arn:aws:iam::987056895854:role/lambda-castle-gdpr-webhook \
		--region us-west-2
clean:
	rm -rf ${EXE}
