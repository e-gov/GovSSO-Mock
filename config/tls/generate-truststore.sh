#!/bin/bash

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

filename=$4
if [ -z "$filename" ]; then
  filename="$application.$environment.truststore.p12"
fi

echo "--------------------------- Generating truststore for '$application.$environment'"

# Remove application existing truststore
rm -f "$application/$filename"

# Generate truststore with CA certificate for application
keytool -noprompt \
  -importcert \
  -alias "$ca.$environment" \
  -file "$ca/$ca.$environment.crt" \
  -storepass changeit \
  -storetype pkcs12 \
  -keystore "$application/$filename"
