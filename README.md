# Castle GDPR Webhook

This Castle GDPR Webhook responds to POST requests from Castle containing GDPR information to be downloaded.

Once called the X-Castle-Signature header will be checked. If the signature is not correct message processing will stop.
If the signature is correct the API version will be checked. If it is correct the user.zip will be downloaded and saved
in a bucket.

For additional information see ./Castle\ GDPR\ Automation.pdf

## Setup

### Create the role

```
aws-okta exec DANGER-dw -- aws iam create-role --role-name lambda-castle-gdpr-webhook \
--assume-role-policy-document file://trust-policy.json
```

### Attach the role

```
aws-okta exec DANGER-dw -- aws iam attach-role-policy --role-name lambda-castle-gdpr-webhook \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

### Verify the role

```
aws-okta exec DANGER-dw -- aws iam get-role --role-name lambda-castle-gdpr-webhook
```

### Create the function

```
$ make create-function
zip function.zip castle-gdpr-webhook
updating: castle-gdpr-webhook (deflated 60%)
aws-okta exec DANGER-dw -- aws lambda create-function \
	       	--function-name castle-gdpr-webhook \
		--runtime go1.x \
  		--zip-file fileb://function.zip \
	       	--handler castle-gdpr-webhook \
  		--role arn:aws:iam::873344020507:role/lambda-castle-gdpr-webhook \
		--region us-west-2
{
    "FunctionName": "castle-gdpr-webhook",
    "FunctionArn": "arn:aws:lambda:us-west-2:873344020507:function:castle-gdpr-webhook",
    "Runtime": "go1.x",
    "Role": "arn:aws:iam::873344020507:role/lambda-castle-gdpr-webhook",
    "Handler": "castle-gdpr-webhook",
    "CodeSize": 8341182,
    "Description": "",
    "Timeout": 3,
    "MemorySize": 128,
    "LastModified": "2020-01-10T17:02:17.128+0000",
    "CodeSha256": "Ycc/b3goc1Bmb0W7ToXA+E5nKk5G0ueY07LFkMu9BJM=",
    "Version": "$LATEST",
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "c9cf5623-9f68-4d43-905e-43fa1a193d6a",
    "State": "Active",
    "LastUpdateStatus": "Successful"
}
```

### Verify that the function can be called

```
$ make test
aws-okta exec DANGER-dw -- aws lambda invoke \
		--function-name castle-gdpr-webhook \
		--invocation-type "RequestResponse" \
		--region us-west-2 \
		response.txt
{
    "StatusCode": 200,
    "ExecutedVersion": "$LATEST"
}
cat response.txt
{"statusCode":500,"headers":null,"multiValueHeaders":null,"body":""}%
```

### Update IAM Role

Update the role in the IAM console. Attach the policy in (./castle-webhook-s3-policy.json).

### Create the s3 bucket

Create the s3 bucket "castle-gdpr-data" if it doesn't exist.

### Build the API Gateway that will invoke the Lambda

Build the new API Gateway that will call the lambda.

In API Gateway create a new "REST API" with the name "CastleGDPRHandler", Endpoint Type "Regional". Don't import
anything from Swagger or the Example API.

In the API Gateway create a new resource "callback".

Next create a method on the resource that allows POST. Configure the integration point to be the Lambda Function
"castle-gdpr-webhook" and select "Use Lambda Proxy integration". Click "OK" and allow the necessary gateway permissions to be created.

Next deploy the API, create a new stage named "production".

The function can now be invoked via the URL, e.g.
"https://71lll8k9g8.execute-api.us-west-2.amazonaws.com/production/callback"

```

curl -v "https://71lll8k9g8.execute-api.us-west-2.amazonaws.com/production/callback" \
    --data-binary @test.json \
    -H 'X-Castle-Signature: DFDUtWGUuoTW8o4uViH78bCVDrSvcdbhsoqC0uYOH0w='
. . .
*   Trying 54.69.46.76...
< HTTP/2 500 
< date: Fri, 10 Jan 2020 18:37:34 GMT
< content-type: application/json
< content-length: 16
< x-amzn-requestid: 086c40aa-a2bf-4003-b06f-2d4a9ab66383
< x-amz-apigw-id: GGOhGEPNPHcFa7Q=
< x-amzn-trace-id: Root=1-5e18c46d-67a6f76e42e0ac3058d5e82b;Sampled=0
< 
* Connection #0 to host 71lll8k9g8.execute-api.us-west-2.amazonaws.com left intact
hmac key invalid%  
```

The post fails because the production hmac key is not available in the local environment and thus the hmac is incorrect;
however the function is being invoked.

### Create the Custom Domain Name

*Note that this has already been done.*

First issue a managed cert for the domain name via the AWS Certificate Manager.

In the API Gateway "Custom Domain Names" create a new HTTP custom domain name for the API Gateway function. Name the
domain name "castlewebhook.optimizely.com". The endpoint will be Regional and select the certificate previously issued.

Then update the DNS to point the DNS entry for castlewebhook.optimizely.com to the target domain name (e.g.
"d-gwgcjup345.execute-api.us-west-2.amazonaws.com").

Add a Base Path Mapping for "v1" to point to the production CastleGDPRHandler environment.

Test the endpoint to make sure it is being invoked.

```
curl -v "https://71lll8k9g8.execute-api.us-west-2.amazonaws.com/production/callback" \
    --data-binary @test.json \
    -H 'X-Castle-Signature: DFDUtWGUuoTW8o4uViH78bCVDrSvcdbhsoqC0uYOH0w='
. . .
*   Trying 54.69.46.76...
< HTTP/2 500
< date: Fri, 10 Jan 2020 18:37:34 GMT
< content-type: application/json
< content-length: 16
< x-amzn-requestid: 086c40aa-a2bf-4003-b06f-2d4a9ab66383
< x-amz-apigw-id: GGOhGEPNPHcFa7Q=
< x-amzn-trace-id: Root=1-5e18c46d-67a6f76e42e0ac3058d5e82b;Sampled=0
<
* Connection #0 to host 71lll8k9g8.execute-api.us-west-2.amazonaws.com left intact
hmac key invalid%
```
