# Burgers API

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Rate Limiting and Caching](#rate-limiting-and-caching)

## Introduction
Burgers API allows you to manage a collection of burgers and their ingredients. You can create new burgers, retrieve details about existing burgers, and more.

## Features
- Create a new burger with ingredients
- Retrieve burgers by various criteria
- Rate limiting to prevent abuse
- Caching to improve performance

## Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/Kontentski/burgersDb.git
    ```
2. Navigate to the project directory:
    ```sh
    cd burgersDb
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```
4. Set up the PostgreSQL database and add the url to the .env file:
   ```sh
   export DATABASE_URL="postgres://username:password@localhost:5432/burgers?sslmode=disable"
   ```
6. This project uses goose for migrations, you can start the app and it will migrate automatically :
    ```sh
    go run .
    ```

## Usage
To start the server, run:
```sh
go run .
```
The server will run on port 8080 by default. You can change the port by setting the PORT environment variable.

## Endpoints
Here are the available endpoints:

- GET `/api/burgers`: Retrieve all burgers
- GET `/api/burgers/:id`: Retrieve a burger by its ID
- GET `/api/burgers/n=:name`: Retrieve a burger by its name
- GET `/api/burgers/f=:name`: Retrieve burgers starting with a specific letter
- GET `/api/burgers/random`: Retrieve a random burger
- GET `/api/burgers/randomten`: Retrieve ten random burgers
- GET `/api/burgers/latest`: Retrieve the latest burgers
- GET `/api/ingredients/:name`: Retrieve an ingredient by its name
- GET `/api/burgers/i=:name`: Retrieve burgers by ingredient name
- GET `/api/burgers/ingredients`: Retrieve burgers by multiple ingredients
- GET `/api/burgers/vegan`: Retrieve vegan burgers
- GET `/api/burgers/nonvegan`: Retrieve non-vegan burgers
- POST `/api/burgers/create`: Create a new burger


- Request Body for a new burger:
```json
{
  "burger": {
    "name": "Classic burger with rocks3.0",
    "description": "A delicious classic rockburger",
    "is_vegan": false,
    "image_url": "assets/rockburger.jpg"
  },
  "ingredients": [
    {
      "name": "rock",
      "description": "Juicy rock patty"
    },
    {
      "name": "boulder",
      "description": "Melty boulder slice"
    }
  ],
  "burgerIngredients": [
    {
      "ingredient_name": "rock",
      "measure": "200g"
    },
    {
      "ingredient_name": "boulder",
      "measure": "2 slices"
    }
  ]
}
```
## Rate Limiting and Caching

## Rate Limiting
Rate limiting is implemented to restrict the number of requests a client can make in a specific period. The middleware is set up to allow 5 requests per second per client IP.

## Caching
Caching is used to improve performance by storing the responses of GET requests for 5 minutes. The POST `/api/burgers/create` and GET `/api/burgers/random` routes are excluded from caching.

## Middleware Configuration
Rate limiting and caching are configured in the middleware package and applied globally in the main function.
