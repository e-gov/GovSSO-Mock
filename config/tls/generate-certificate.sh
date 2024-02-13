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

applicationName=$3
if [ -z "$applicationName" ]; then
  echo "\$applicationName is empty"
  exit 1
fi

host=$applicationName.$environment # ex. admin.localhost

basedir=$(dirname "$0")
mkdir -p "${basedir}/$applicationName"
privateKeyFile="${basedir}/$applicationName/$host.key.pem"
requestFile="${basedir}/$applicationName/$host.csr.pem"
certificateFile="${basedir}/$applicationName/$host.crt.pem"
caPrivateKeyFile="${basedir}/$ca/$ca.$environment.key.pem"
caCertificateFile="${basedir}/$ca/$ca.$environment.crt.pem"
keystoreFile="${basedir}/$applicationName/$host.keystore.p12"

[[ -f "${privateKeyFile}" ]] || {
  echo "--------------------------- Generating private key for '$host'"
  openssl ecparam \
    -name prime256v1 \
    -genkey \
    -out "${privateKeyFile}"
  chmod 644 "${privateKeyFile}"
}

[[ -f "${requestFile}" ]] || {
  echo "--------------------------- Generating CSR for '$host'"
  # MSYS_NO_PATHCOW=1 needed for Git Bash on Windows users - unable to handle "/"-s in -subj parameter.
  MSYS_NO_PATHCONV=1 \
  openssl req \
    -new \
    -sha512 \
    -nodes \
    -key "${privateKeyFile}" \
    -subj "/CN=$host" \
    -out "${requestFile}"
  chmod 644 "${requestFile}"
}

[[ -f "${certificateFile}" ]] || {
  echo "--------------------------- Generating certificate for '$host'"
  # Configure subject alternate names. Passed to openssl.cnf
  export SAN="DNS:$host"
  openssl x509 \
    -req \
    -sha512 \
    -in "${requestFile}" \
    -CA "${caCertificateFile}" \
    -CAkey "${caPrivateKeyFile}" \
    -CAcreateserial \
    -days 363 \
    -extfile "${basedir}/openssl.cnf" \
    -out "${certificateFile}"
  chmod 644 "${certificateFile}"
}

[[ -f "${keystoreFile}" ]] || {
  echo "--------------------------- Generating keystore for '$host'"
  openssl pkcs12 \
    -export \
    -name "$host" \
    -in "${certificateFile}" \
    -inkey "${privateKeyFile}" \
    -passout pass:changeit \
    -out "${keystoreFile}"
  chmod 644 "${keystoreFile}"
}
