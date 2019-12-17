# Castle GDPR Webhook

This webhook helps automate Castle GDPR requests.

## Overview

The castle GDPR webhook that works as follows and could be done inside an AWS lambda. It will have two API endpoints.

POST /callback - This is a webhook path that will be registered with Castle.io. It will accept webhooks with the following content which was defined by Castle.

```
{
    "api_version": "v1",
    "app_id": "382395555537961",
    "type": "$gdpr.subject_access_request.completed",
    "created_at": "2019-12-01T19:38:28.483Z",
    "data": {
    "id": "test",
    "download_url": "https://url/user.zip"
    "download_url_expires_at": "2020-12-12T00:00.00Z",
    "user_id": "2",
    "user_traits": {
        "id": "2",
            "email": "email@example.com"
        }
    }
}
```

For additional information see ./Castle\ GDPR\ Automation.pdf

Once called the X-Castle-Signature header will be checked. If the signature is not correct message processing will stop. If the signature is correct the API version will be checked. If it is correct the user.zip will be downloaded and saved in a bucket.

GET /user/unique_id

This will return the user.zip saved from when it was provided by castle calling the webhook. The caller must be authenticated. This could be done with a secret header. After the zip is downloaded it can be removed.

The intended caller is https://github.com/optimizely/hermes-airflow.
