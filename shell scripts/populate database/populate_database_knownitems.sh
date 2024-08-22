#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"serial_number": 123456789, "long_name": "Great Value Whole Vitamin D Milk, Gallon, Plastic, Jug, 128 Fl Oz", "short_name": "Milk", "tags": ["drinkable", "milk"], "item_type": 1, "measurement": 2, "container_size": 3785}' localhost:4000/v1/knownitems
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"serial_number": 543534634, "long_name": "OxiClean Foam-Tastic Foaming Bathroom Cleaner, Removes Soap Scum, Grime & Stains, Lemon Scent, 19 oz", "short_name": "Bathroom Cleaner", "tags": ["cleaning", "bathroom"], "item_type": 2, "measurement": 2, "container_size": 850}' localhost:4000/v1/knownitems
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"serial_number": 645645712, "long_name": "Fresh Banana Fruit, Each", "short_name": "Banana", "tags": ["fruit"], "item_type": 1, "measurement": 1, "container_size": 1}' localhost:4000/v1/knownitems
