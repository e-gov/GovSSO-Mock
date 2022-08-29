#!/bin/bash

environment=$1
if [ -z "$environment" ]; then
  echo "\$environment is not assigned"
  exit 1
fi

caName=$2
if [ -z "$caName" ]; then
  echo "\$caName is empty"
  exit 1
fi

caFullName="$caName-ca"

echo "--------------------------- Generating '$caName.$environment' CA certificate"

mkdir -p "$caFullName"

# Generate CA private key
openssl ecparam \
  -genkey \
  -name prime256v1 \
  -out "$caFullName/$caFullName.$environment.key"

# Generate CA certificate
MSYS_NO_PATHCONV=1 \
  openssl req \
  -x509 \
  -new \
  -sha512 \
  -nodes \
  -key "$caFullName/$caFullName.$environment.key" \
  -days 365 \
  -subj "/C=EE/L=Tallinn/O=$caName-$environment/CN=$caFullName.$environment" \
  -out "$caFullName/$caFullName.$environment.crt"
