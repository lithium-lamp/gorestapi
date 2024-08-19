#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "unit"}' localhost:4000/v1/measurements
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "dl"}' localhost:4000/v1/measurements
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "ml"}' localhost:4000/v1/measurements