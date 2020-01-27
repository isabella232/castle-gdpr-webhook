#!/bin/sh

# get the url via `terraform output`
URL=`make output | grep base_url | awk -F '=' '{print $2}' | awk '{$1=$1};1'`
PUBLIC_URL="https://castlewebhook-test.optimizely.com/v1"
JSON="../test.json"
HMAC=`cat $JSON | openssl dgst -binary -sha256 -hmac "$TF_VAR_hmac_secret" | openssl base64`

echo "Test 1: Calling the API endpoint directly. Make sure HMAC_SECRET is set in your environment."
echo "curl \"$URL\" --data-binary @${JSON} -H \"X-Castle-Signature: $HMAC\""

curl "$URL" --data-binary @${JSON} -H "X-Castle-Signature: $HMAC"
if [ $? -eq 0 ]; then
    echo "call succeeded 😀"
else
    echo "call failed 😱"
fi

echo ""
echo "Test 2: Calling the public DNS entry API endpoint. Make sure HMAC_SECRET is set in your environment."
echo "Note: This will fail unless the DNS for castlewebhook-test.optimizely.com --> Target Domain Name"

URL=$PUBLIC_URL
echo ""
echo "To proceed run the following:"
echo "curl \"$URL\" --data-binary @${JSON} -H \"X-Castle-Signature: $HMAC\""
echo ""
