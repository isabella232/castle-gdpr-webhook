# Castle GDPR Webhook

This Castle GDPR Webhook responds to POST requests from Castle containing GDPR information to be downloaded.

Once called the X-Castle-Signature header will be checked. If the signature is not correct message processing will stop.
If the signature is correct the API version will be checked. If it is correct the user.zip will be downloaded and saved
in a bucket.

For additional information see ./Castle\ GDPR\ Automation.pdf

## Setup

Create the role

```
aws-okta exec DANGER-security -- aws iam create-role --role-name lambda-castle-gdpr-webhook \
--assume-role-policy-document file://trust-policy.json
```

Attach the role

```
aws-okta exec DANGER-security -- aws iam attach-role-policy --role-name lambda-castle-gdpr-webhook \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

Verify the role

```
aws-okta exec DANGER-security -- aws iam get-role --role-name lambda-castle-gdpr-webhook
```

TODO describe granting access to the S3 bucket.
