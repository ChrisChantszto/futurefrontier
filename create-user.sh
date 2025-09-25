#!/bin/bash

# Create a user directly using the API
curl -X POST http://localhost:8080/api/auth/init-first-user \
  -H "Content-Type: application/json" \
  -d '{
    "email": "email",
    "password": "pw"
  }'

echo -e "\nUser created! Now you can use this account with OTP login."
