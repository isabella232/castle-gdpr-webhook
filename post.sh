#!/bin/sh

#curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
#    -d @valid.json \
#    -H 'Content-Type: application/json' -H 'X-Castle-Signature: DSgTD4h47B0IexVVMCcbyG1T80LDnFmgiYXWqEmu/gI='

#curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
#    -d @onelinetest.json \
#    -H 'Content-Type: application/json' -H 'X-Castle-Signature: VYDK6AH2SbaY7GBEqSeG9bgpwtfwrPHLHHellr8cn+Y='

#curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
#    -d @test.json \
#    -H 'Content-Type: application/json' -H 'X-Castle-Signature: DFDUtWGUuoTW8o4uViH78bCVDrSvcdbhsoqC0uYOH0w='

# -d strips CR and LF but this works
#curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
#    -d @request.txt \
#    -H 'Content-Type: application/json' -H 'X-Castle-Signature: ze6/FiQ/bes2eimxiOK4/aLaM2FclJaAYOT6OE6DU5o='

# instead do binary
#curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
#    --data-binary @request.txt \
#    -H 'Content-Type: application/json' -H 'X-Castle-Signature: DFDUtWGUuoTW8o4uViH78bCVDrSvcdbhsoqC0uYOH0w='

curl -v "https://uj05s8dbk7.execute-api.us-west-2.amazonaws.com/test/callback" \
    --data-binary @test.json \
    -H 'X-Castle-Signature: DFDUtWGUuoTW8o4uViH78bCVDrSvcdbhsoqC0uYOH0w='
