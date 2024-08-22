
# @householdingindex/gorestapi

This repository is meant to serve as the backend to the project @householdingindex.

The @householdingindex project is a tool for cataloging efficiently, such as food or household items.


## API Reference

Although the api is not meant to be used directly via ex. curl commands, as it should communicate via a front end, either mobile or web - it still feels relevant to add endpoints for the sake of documentation.


### The "v1/healthcheck" endpoint

This endpoint displays information about api status, system info, environment and versioning.

#### Display health status

```http
  GET /v1/healthcheck
```











### The "v1/recipies" endpoint

#### Get all recipies

```http
  GET /v1/recipies
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipies:read` | `permission` | **Required**. Account permissions |

#### Post recipe

```http
  POST /v1/recipies
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipies:write` | `permission` | **Required**. Account permissions |
| `name `      | `string` | **Required** Recipe name |
| `description `      | `string` | **Required** Recipe description|
| `cooking_steps `      | `[]string` | **Required** Slice containing recipe instructions |
| `cook_time_minutes `      | `int` | **Required** Expected minutes required to cook recipe|
| `portions `      | `int` | **Required** Expected resulting portions from recipe|
| `tags `      | `[]string` | **Required** Slice containing tags for recipies ex. "dessert", "french" |

#### Get recipe

```http
  GET /v1/recipies/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipies:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch recipe

```http
  PATCH /v1/recipies/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipies:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `name `      | `string` | Recipe name |
| `description `      | `string` | Recipe description|
| `cooking_steps `      | `[]string` | Slice containing recipe instructions |
| `cook_time_minutes `      | `int` | Expected minutes required to cook recipe|
| `portions `      | `int` | Expected resulting portions from recipe|
| `tags `      | `[]string` | Slice containing tags for recipies ex. "dessert", "french" |

#### Delete recipe

```http
  DELETE /v1/recipies/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipies:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |













### The "v1/ingredients" endpoint

#### Get all ingredients

```http
  GET /v1/ingredients
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `ingredients:read` | `permission` | **Required**. Account permissions |

#### Post ingredient

```http
  POST /v1/ingredients
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `ingredients:write` | `permission` | **Required**. Account permissions |
| `name `      | `string` | **Required** Ingredient name |
| `tags `      | `[]string` | **Required** Slice containing tags for ingredients ex. "cheese", "milk" |

#### Get ingredient

```http
  GET /v1/ingredients/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `ingredients:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch ingredient

```http
  PATCH /v1/ingredients/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `ingredients:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `name `      | `string` | Ingredient name |
| `tags `      | `[]string` | Slice containing tags for ingredients ex. "cheese", "milk" |

#### Delete ingredient

```http
  DELETE /v1/ingredients/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `ingredients:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |




### The "v1/recipeingredients" endpoint

#### Get all recipe ingredients

```http
  GET /v1/recipeingredients
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipeingredients:read` | `permission` | **Required**. Account permissions |

#### Post recipe ingredient

```http
  POST /v1/recipeingredients
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipeingredients:write` | `permission` | **Required**. Account permissions |
| `recipe_id `      | `int` | **Required** Recipe id |
| `ingredient_id `      | `int` | **Required** Ingredient id |
| `amount `      | `int` | **Required** How much of specific measurement is required ex. *15* units |
| `measurement `      | `int` | **Required** Measurement id |

#### Get recipe ingredient

```http
  GET /v1/recipeingredients/${recipe_id}/${ingredient_id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipeingredients:read` | `permission` | **Required**. Account permissions |
| `recipe_id`      | `int` | **Required**. Recipe id of item to fetch |
| `ingredient_id`      | `int` | **Required**. Ingredient id of item to fetch |

Note: Both `recipe_id` and `ingredient_id` must exist in unison, meaning fetching an ingredient "5" that is recipe "7" means the get request should be sent to v1/recipeingredients/7/5.

#### Patch recipe ingredient

```http
  PATCH /v1/recipeingredients/${recipe_id}/${ingredient_id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipeingredients:write` | `permission` | **Required**. Account permissions |
| `recipe_id `      | `int` | **Required** Recipe id of item to fetch|
| `ingredient_id `      | `int` | **Required** Ingredient id of item to fetch |
| `recipe_id `      | `int` | Can be passed into body to alter recipe id |
| `ingredient_id `      | `int` | Can be passed into body to alter ingredient id |
| `amount `      | `int` | How much of specific measurement is required ex. *15* units |
| `measurement `      | `int` | Measurement id |

Note: Both `recipe_id` and `ingredient_id` must exist in unison, meaning fetching an ingredient "5" that is recipe "7" means the get request should be sent to v1/recipeingredients/7/5.

#### Delete recipe ingredient

```http
  DELETE /v1/recipeingredients/${recipe_id}/${ingredient_id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `recipeingredients:write` | `permission` | **Required**. Account permissions |
| `recipe_id`      | `int` | **Required**. Recipe id of item to fetch |
| `ingredient_id`      | `int` | **Required**. Ingredient id of item to fetch |

Note: Both `recipe_id` and `ingredient_id` must exist in unison, meaning fetching an ingredient "5" that is recipe "7" means the get request should be sent to v1/recipeingredients/7/5.




### The "v1/availableitems" endpoint

#### Get all available items

```http
  GET /v1/availableitems
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `availableitems:read` | `permission` | **Required**. Account permissions |

#### Post available item

```http
  POST /v1/availableitems
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `availableitems:write` | `permission` | **Required**. Account permissions |
| `knownitems_id `      | `int` | **Required** Known item id |
| `expiration_at `      | `time.Time` | **Required** Time in RFC3339 format, ex. 2024-08-10T10:30:20Z|
| `container_size `      | `int` | **Required** Relative to unit given in measurement, ex. 3 units ...|

#### Get available item

```http
  GET /v1/availableitems/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `availableitems:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch available item

```http
  PATCH /v1/availableitems/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `availableitems:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `knownitems_id `      | `int` | **Required** Known item id |
| `expiration_at `      | `time.Time` | Time in RFC3339 format, ex. 2024-08-10T10:30:20Z|
| `container_size `      | `int` |  Relative to unit given in measurement, ex. 3 units ...|

#### Delete available item

```http
  DELETE /v1/availableitems/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `availableitems:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |








### The "v1/knownitems" endpoint

#### Get all known items

```http
  GET /v1/knownitems
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `knownitems:read` | `permission` | **Required**. Account permissions |

#### Post known item

```http
  POST /v1/knownitems
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `knownitems:write` | `permission` | **Required**. Account permissions |
| `serial_number `      | `int` | **Required** Serial number of known item |
| `long_name `      | `string` | **Required** Full length name |
| `short_name `      | `string` | **Required** Shorthand name |
| `tags `      | `[]string` | **Required** Tags ex. "grilling", "breakfast", "cheese" |
| `item_type `      | `int` | **Required** Item type id |
| `measurement `      | `int` | **Required** Measurement id |
| `container_size `      | `int` | **Required** Relative to unit given in measurement, ex. 3 units ...|

#### Get known item

```http
  GET /v1/knownitems/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `knownitems:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch known item

```http
  PATCH /v1/knownitems/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `knownitems:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `knownitems_id `      | `int` | **Required** Known item id |
| `serial_number `      | `int` | Serial number of known item |
| `long_name `      | `string` | Full length name |
| `short_name `      | `string` | Shorthand name |
| `tags `      | `[]string` | Tags ex. "grilling", "breakfast", "cheese" |
| `item_type `      | `int` | Item type id |
| `measurement `      | `int` | Measurement id |
| `container_size `      | `int` | Relative to unit given in measurement, ex. 3 units ...|


#### Delete known item

```http
  DELETE /v1/knownitems/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `knownitems:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |




### The "v1/tags" endpoint

#### Get all tags

```http
  GET /v1/tags
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `tags:read` | `permission` | **Required**. Account permissions |

#### Post tag

```http
  POST /v1/tags
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `tags:write` | `permission` | **Required**. Account permissions |
| `itemtype `      | `int` | **Required** Item type id |
| `name `      | `string` | **Required** Name of item type ex. foodstuff |

#### Get tag

```http
  GET /v1/tags/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `tags:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch tag

```http
  PATCH /v1/tags/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `itemtype `      | `int` | Item type id |
| `name `      | `string` | Name of item type ex. foodstuff |

#### Delete tag

```http
  DELETE /v1/tags/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |




### The "v1/itemtypes" endpoint

#### Get all itemtypes

```http
  GET /v1/itemtypes
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:read` | `permission` | **Required**. Account permissions |

#### Post item type

```http
  POST /v1/itemtypes
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:write` | `permission` | **Required**. Account permissions |
| `name `      | `string` | **Required** Name of item type ex. foodstuff|

#### Get item type

```http
  GET /v1/itemtypes/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch item type

```http
  PATCH /v1/itemtypes/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `name `      | `string` | Name of item type ex. foodstuff, note: patch can be used without any changes|

#### Delete item type

```http
  DELETE /v1/itemtypes/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |






### The "v1/measurements" endpoint

#### Get all measurements

```http
  GET /v1/measurements
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `measurements:read` | `permission` | **Required**. Account permissions |

#### Post measurements

```http
  POST /v1/measurements
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `measurements:write` | `permission` | **Required**. Account permissions |
| `name `      | `string` | **Required** Name of measurement ex. units|

#### Get measurements

```http
  GET /v1/measurements/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `measurements:read` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Patch measurements

```http
  PATCH /v1/measurements/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `itemtypes:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |
| `name `      | `string` | Name of measurement ex. units, note: patch can be used without any changes|

#### Delete measurements

```http
  DELETE /v1/measurements/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `measurements:write` | `permission` | **Required**. Account permissions |
| `id`      | `int` | **Required**. Id of item to fetch |



## The "v1/users" endpoint

#### Register user

```http
  POST /v1/users
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `email`      | `string` | **Required**. Valid email string not currently in use, ex. john@domain.com |
| `name `      | `string` | **Required** Firstname and lastname, ex. John Doe |
| `password `      | `string` | **Required** |

#### Activate user

```http
  PUT /v1/users/activated
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `token`      | `string` | **Required**. Valid token, sent to user via email |




## The "v1/tokens" endpoint

#### Get authentication token for specific user

```http
  POST /v1/tokens/authentication
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `email`      | `string` | **Required**. Existing email for registered user |
| `password`      | `string` | **Required**. Password correlating with given email |




## The "debug/vars" endpoint

#### Display runtime stats

```http
  GET /debug/vars
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `bearer token` | `string` | **Required**. A bearer token belonging to an authorized user in the format "Authorization: Bearer XXXXXXXXXXXXXXXX", passed as header |
| `metrics:view` | `permission` | **Required**. Account permissions |
