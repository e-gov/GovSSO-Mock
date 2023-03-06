#!/bin/bash
set -eu

cd tls
bash generate-tls-resources.sh

cd ../id-token
bash generate-id-token-signing-keys.sh

echo "--------------------------- All resources generated"

# Prevents script window to be closed after completion
echo -e "\nPress any key to exit the script."
read -rn1
