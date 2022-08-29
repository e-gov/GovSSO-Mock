#!/bin/bash

cd tls || exit
bash generate-tls-resources.sh

cd ../id-token || exit
bash generate-id-token-signing-keys.sh

echo "--------------------------- All resources generated"

# Prevents script window to be closed after completion
echo -e "\nPress any key to exit the script."
read -rn1
