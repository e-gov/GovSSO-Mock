#!/bin/bash

# Recursively remove all directories from current path
cd "$(command dirname -- "${0}")" || exit
rm -rf ./*/

bash generate-ca-certificate.sh 'localhost' 'govsso'

bash generate-certificate.sh 'localhost' 'govsso-ca' 'client'
bash generate-certificate.sh 'localhost' 'govsso-ca' 'govsso-mock'

bash generate-truststore.sh 'localhost' 'govsso-ca' 'client'
bash generate-truststore.sh 'localhost' 'govsso-ca' 'govsso-mock'
