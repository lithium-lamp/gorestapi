#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"knownitems_id": 3, "expiration_at": "2024-08-10T10:30:20Z", "container_size": 3785}' localhost:4000/v1/availableitems
