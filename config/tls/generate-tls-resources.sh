#!/bin/bash

# Recursively remove all directories from current path
cd "$(command dirname -- "${0}")" || exit
rm -rf ./*/

bash generate-ca-certificate.sh 'test' 'govsso'

bash generate-certificate.sh 'test' 'govsso-ca' 'client'
bash generate-certificate.sh 'test' 'govsso-ca' 'govsso-mock'

bash generate-truststore.sh 'test' 'govsso-ca' 'client'
bash generate-truststore.sh 'test' 'govsso-ca' 'govsso-mock'
