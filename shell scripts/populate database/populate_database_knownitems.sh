#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"serial_number": 123456789, "long_name": "Great Value Whole Vitamin D Milk, Gallon, Plastic, Jug, 128 Fl Oz", "short_name": "Milk", "tags": ["drinkable", "milk"], "item_type": 1, "measurement": 2, "container_size": 3785}' localhost:4000/v1/knownitems
