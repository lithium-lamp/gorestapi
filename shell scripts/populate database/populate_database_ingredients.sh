#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "apple", "tags": ["fruit"]}' localhost:4000/v1/ingredients
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "sugar", "tags": ["baking", "sugar"]}' localhost:4000/v1/ingredients


# name string,
# tags []string]
