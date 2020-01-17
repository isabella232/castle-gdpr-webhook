#!/bin/sh

# get the url via `terraform output`
URL="https://eq53ky15j0.execute-api.us-west-2.amazonaws.com/test"
JSON="../test.json"
HMAC=`cat $JSON | openssl dgst -binary -sha256 -hmac "$HMAC_SECRET" | openssl base64`

echo "Calling the API endpoint. Make sure HMAC_SECRET is set in your environment."

curl "$URL" --data-binary @${JSON} -H "X-Castle-Signature: $HMAC"
if [ $? -eq 0 ]; then
    echo "call succeeded 😀"
else
    echo "call failed 😱"
fi
