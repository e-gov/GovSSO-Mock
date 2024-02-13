#!/bin/bash
set -eu

environment=$1
if [ -z "$environment" ]; then
  echo "\$environment is not assigned"
  exit 1
fi

ca=$2
if [ -z "$ca" ]; then
  echo "\$ca is empty"
  exit 1
fi

application=$3
if [ -z "$application" ]; then
  echo "\$application is empty"
  exit 1
fi

basedir=$(dirname "$0")
truststoreFile="${basedir}/$application/$application.$environment.truststore.p12"
caCertificateFile="${basedir}/$ca/$ca.$environment.crt.pem"

[[ -f "${truststoreFile}" ]] || {
  echo "--------------------------- Generating truststore for '$application.$environment'"
  keytool -noprompt \
    -importcert \
    -alias "$ca.$environment" \
    -file "${caCertificateFile}" \
    -storepass changeit \
    -storetype pkcs12 \
    -keystore "${truststoreFile}"
  chmod 644 "${truststoreFile}"
}
