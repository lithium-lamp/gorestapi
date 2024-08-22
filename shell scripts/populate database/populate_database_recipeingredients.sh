#!/usr/bin/env bash

source "./shell scripts/populate database/config.sh"

curl -X POST -H "Authorization: Bearer $admin_token" -d '{"recipe_id": 1, "ingredient_id": 1, "amount": 3, "measurement": 1}' localhost:4000/v1/recipeingredients
curl -X POST -H "Authorization: Bearer $admin_token" -d '{"recipe_id": 1, "ingredient_id": 2, "amount": 20, "measurement": 2}' localhost:4000/v1/recipeingredients


# recipe_id int,
# ingredient_id int,
# amount int,
# measurement int,