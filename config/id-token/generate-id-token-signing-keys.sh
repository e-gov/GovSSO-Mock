#!/bin/bash

echo "--------------------------- Generating RSA public and private key pair for signing identity tokens"

openssl genrsa \
  -out idTokenSign.key \
  -traditional \
  4096

openssl rsa \
  -in idTokenSign.key \
  -pubout > idTokenSign.pub
