# Castle GDPR Webhook

This webhook helps automate Castle GDPR requests. It works as follows.

## Overview

This is a small service that runs in optimizely-hrd. It could reuse the setup & deployment found in
https://github.com/optimizely/token-service and expose 3 APIs.

- DELETE /user/\<email\> - This will cause the user_id to be looked up in the datastore. If no user is found a 200 OK will
    be returned. If a user is found a request to delete the user will be issued to the castle API. Upon success a 200 OK
    status will be returned.

- POST /user/\<email\> - This will cause the user_id to be looked up in the datastore. If a user is found a request to
    access the user records will be issued to the castle API. If a user is found and data is stored in a bucket for the
    user the data (user.zip) will be returned to the caller.

- POST /callback - This is a webhook path that will be registered with Castle.io. It will accept webhooks with the
    following content which was defined by Castle.

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

Once called the X-Castle-Signature header will be checked. If the signature is not correct message processing will
stop. If the signature is correct the API version will be checked. If it is correct the user.zip will be downloaded
and saved in a bucket. The user.zip is what is returned when callers perform a POST to /user/<email>.

The caller is the castle gdpr automation inside https://github.com/optimizely/hermes-airflow. It will have to be updated to
call the service APIs above.

## Configuration

The Castle API secret will need to be available in an environment variable.

The authentication header secret must be available in an environment variable.

## Authentication

The callers will be expected to have set an HTTP secret header named X-Castle-Webhook containing the aforementioned
secret.

The only caller of this server will be the Optimizely GDPR automation in hermes-airflow.

## Additional TODOs

- Open IT ticket to put castle behind Okta.
