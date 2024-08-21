#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"name": "Apple pie", "description": "A delicious apple pie", "cooking_steps": ["First slice the apples", "Then do something else"], "cook_time_minutes": 60, "portions": 6, "tags": ["dessert", "pie"]}' localhost:4000/v1/recipies

# name text,
# description text,
# cooking_steps text [] DEFAULT '{}',
# cook_time_minutes int,
# portions int,
# tags text [] DEFAULT '{}',
