#!/bin/bash
set -eu

basedir=$(dirname "$0")
bash "${basedir}/id-token/generate-id-token-signing-keys.sh"
bash "${basedir}/tls/generate-tls-resources.sh"

# TODO: Find a better solution than calling chmod everywhere.
