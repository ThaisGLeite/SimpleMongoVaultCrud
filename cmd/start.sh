#!/bin/bash

# Stop any running containers
docker-compose -f ../docker/docker-compose.yml down

# Build the Vault Docker image and start the Vault Docker container
docker-compose -f ../docker/docker-compose.yml up --build -d vault

# Wait for the Vault server to start
echo "Waiting for Vault server to start..."
while ! docker logs docker-vault-1 2>&1 | grep -q 'Vault server started! Log data will stream in below:'
do
  sleep 1
done

# Initialize the Vault server and store the keys in keys.txt
docker exec docker-vault-1 vault operator init > keys.txt

# Extract unseal keys and root token
KEY1=$(grep 'Unseal Key 1:' keys.txt | awk '{print $NF}')
KEY2=$(grep 'Unseal Key 2:' keys.txt | awk '{print $NF}')
KEY3=$(grep 'Unseal Key 3:' keys.txt | awk '{print $NF}')
ROOT_TOKEN=$(grep 'Initial Root Token:' keys.txt | awk '{print $NF}')

# Delete the keys file
rm keys.txt

# Unseal the Vault
docker exec docker-vault-1 vault operator unseal $KEY1
docker exec docker-vault-1 vault operator unseal $KEY2
docker exec docker-vault-1 vault operator unseal $KEY3

# Print the root token and Export the root token as VAULT_TOKEN
echo "Vault Root Token: $ROOT_TOKEN"
export VAULT_TOKEN=$ROOT_TOKEN

# Store MongoDB credentials in Vault 
docker exec -e VAULT_TOKEN=$ROOT_TOKEN docker-vault-1 vault secrets enable -path=secret kv
docker exec -e VAULT_TOKEN=$ROOT_TOKEN docker-vault-1 vault kv put secret/data/mongodb data='{"username":"thaisdev","password":"DevEnv123"}'

# Start the MongoDB Docker container
docker-compose -f ../docker/docker-compose.yml up -d mongodb

# Wait for MongoDB to start
echo "Waiting for MongoDB to start..."
sleep 3

# Setup MongoDB root user
docker exec docker-mongodb-1 mongo admin --eval 'db.dropUser("thaisdev");'
docker exec docker-mongodb-1 mongo admin --eval 'db.createUser({user: "thaisdev", pwd: "DevEnv123", roles: [ { role: "userAdminAnyDatabase", db: "admin" }, { role: "dbOwner", db: "devenv" } ]});'

# Start the API Docker container
docker-compose -f ../docker/docker-compose.yml up --build api
