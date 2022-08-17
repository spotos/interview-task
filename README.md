
# Interview

# Prerequisites

- git
- Golang or other IDE
- Docker and docker-compose
- Postman or other software to send request to local API


# Setup instructions

1. Run ```make init && make start``` to start development server.
1. Modify `/etc/hosts` file to contain: `127.0.0.1 interview.localhost`
1. API should now be accessible through http://interview.localhost:8080/v1/hello
1. API container logs accessible through ```make logs```

# Task description

We need to implement a simple API that allows users to create new freights and returns a list of these freights. 

Freights will be stored in Redis once the rate has been calculated. Rate calculation should be asynchronous and processed via pub/sub queue.

Rules for rate calculation are:
- Base rate for transportation of each freight is 1000 euros;
- If freight is shipped from Lithuania, rate is reduced by 15%;
- If freight is shipped to Germany, rate is increased by 15%;
- If freight is of a high value, rate is increased by 20%.

We will need two endpoints, one for freight creation and one to return a list of freights. When returning a list of freights, you should use a cached response that expires every 15 minutes and should not include freights that do not have rate calculated.

You will find specifications of both endpoints below. Please note that some of the parameters are defined as an enums and should not allow any other values.

## API specification

### Create new freight

`[POST] /v1/freights`

**Request body**

```json
{
  "type": "string(ftl|ltl)",
  "loading_location": {
    "country_code": "string(LTU|DEU|POL|FRA)",
    "postcode": "string"
  },
  "unloading_location": {
    "country_code": "string(LTU|DEU|POL|FRA)",
    "postcode": "string"
  },
  "weight": "integer",
  "high_value": "boolean"
}
```

### Get list of freights

`[GET] /v1/freights`

**Response body**

```json
[
  {
    "type": "string(ftl|ltl)",
    "loading_location": {
      "country_code": "string(LTU|DEU|POL|FRA)",
      "postcode": "string"
    },
    "unloading_location": {
      "country_code": "string(LTU|DEU|POL|FRA)",
      "postcode": "string"
    },
    "weight": "integer",
    "high_value": "boolean",
    "rate": "integer"
  }
]
```
