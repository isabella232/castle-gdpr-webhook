#!/bin/sh

#curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
#    -d @valid.json \
#    -H 'Content-Type: application/json' -H 'X-Castle-Signature: DSgTD4h47B0IexVVMCcbyG1T80LDnFmgiYXWqEmu/gI='

curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
    -d @test.json \
    -H 'Content-Type: application/json' -H 'X-Castle-Signature: DSgTD4h47B0IexVVMCcbyG1T80LDnFmgiYXWqEmu/gI='
