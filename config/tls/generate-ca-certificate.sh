#!/bin/bash
set -eu

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

basedir=$(dirname "$0")
mkdir -p "${basedir}/$caFullName"
privateKeyFile="${basedir}/$caFullName/$caFullName.$environment.key.pem"
certificateFile="${basedir}/$caFullName/$caFullName.$environment.crt.pem"

[[ -f "${privateKeyFile}" ]] || {
  echo "--------------------------- Generating CA private key"
  openssl ecparam \
    -genkey \
    -name prime256v1 \
    -out "${privateKeyFile}"
  chmod 644 "${privateKeyFile}"
}

[[ -f "${certificateFile}" ]] || {
  echo "--------------------------- Generating CA certificate"
  # MSYS_NO_PATHCOW=1 needed for Git Bash on Windows users - unable to handle "/"-s in -subj parameter.
  MSYS_NO_PATHCONV=1 \
  openssl req \
    -x509 \
    -new \
    -sha512 \
    -nodes \
    -key "${privateKeyFile}" \
    -days 365 \
    -subj "/CN=$caFullName.$environment" \
    -out "${certificateFile}"
  chmod 644 "${certificateFile}"
}
