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

Allow reading from the SSM, TODO: replace with actual policy

```
arn:aws:iam::987056895854:role/lambda-castle-gdpr-webhook


{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:DescribeParameters"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameters"
            ],
            "Resource": "arn:aws:ssm:us-west-2:987056895854:parameter/hermes/prod/castle/api_secret"
        }
    ]
}
```
