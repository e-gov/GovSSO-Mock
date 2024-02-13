#!/bin/bash
set -eu

basedir=$(dirname "$0")
privateKeyFile="${basedir}/id-token-sign.key.pem"
publicKeyFile="${basedir}/id-token-sign.pub.pem"

[[ -f "${privateKeyFile}" ]] || {
  echo "--------------------------- Generating RSA public and private key pair for signing identity tokens"
  openssl genrsa \
  -out "${privateKeyFile}" \
  -traditional \
  4096
  chmod 644 "${privateKeyFile}"
}

[[ -f "${publicKeyFile}" ]] || {
  echo "--------------------------- Writing RSA public key for signing identity tokens"
  openssl rsa \
    -in "${privateKeyFile}" \
    -pubout \
    -out "${publicKeyFile}"
  chmod 644 "${publicKeyFile}"
}
