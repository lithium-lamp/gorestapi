#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "food"}' localhost:4000/v1/itemtypes
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "householding"}' localhost:4000/v1/itemtypes