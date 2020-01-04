# Castle GDPR Webhook

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
