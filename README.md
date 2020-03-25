# Castle GDPR Webhook

This Castle GDPR Webhook responds to POST requests from Castle containing GDPR information to be downloaded.

Once called the X-Castle-Signature header will be checked. If the signature is not correct message processing will stop.
If the signature is correct the API version will be checked. If it is correct the user.zip will be downloaded and saved
in a bucket.

For additional information see ./Castle\ GDPR\ Automation.pdf

## Setup

```
$ make create-function
$ make deploy-production
or
$ make deploy-test
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

## Important

The `terraform` directory containing automation is no longer used for production deployment. The production deployment
terrafrom has been integrated to the [dw-infrastructure](https://github.com/optimizely/dw-infrastructure) repository.
