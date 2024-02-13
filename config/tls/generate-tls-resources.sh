#!/bin/bash
set -eu

basedir=$(dirname "$0")

bash "${basedir}/generate-ca-certificate.sh" 'localhost' 'govsso'

bash "${basedir}/generate-certificate.sh" 'localhost' 'govsso-ca' 'client'
bash "${basedir}/generate-certificate.sh" 'localhost' 'govsso-ca' 'govsso-mock'

bash "${basedir}/generate-truststore.sh" 'localhost' 'govsso-ca' 'client'
