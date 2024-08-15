
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
| `expiration_at `      | `time.Time` | **Required** Time in RFC3339 format, ex. 2024-08-10T10:30:20Z|
| `long_name `      | `string` | **Required** Full descriptive product name |
| `short_name `      | `string` | **Required** Short general item name|
| `item_type `      | `int` | **Required** Id corresponding to table "itemtypes"|
| `measurement `      | `int` | **Required** Id corresponding to table "measurements"|
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
| `expiration_at `      | `time.Time` | Time in RFC3339 format, ex. 2024-08-10T10:30:20Z|
| `long_name `      | `string` | Full descriptive product name |
| `short_name `      | `string` | Short general item name|
| `item_type `      | `int` | Id corresponding to table "itemtypes"|
| `measurement `      | `int` | Id corresponding to table "measurements"|
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
