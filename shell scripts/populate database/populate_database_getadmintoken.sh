#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -d '{"email": "admin@admin.com", "password": "pa55word"}' localhost:4000/v1/tokens/authentication
