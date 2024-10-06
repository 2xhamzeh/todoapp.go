# ToDo API with Golang

## About:
#### This is a simple api written in go. It helps you manage todo items

---
## Setup:

#### Database:
You can install mongodb community edition and run it locally or pull a docker image and run a container.

Container setup:
1. Pull the Image
```
docker pull mongodb/mongodb-community-server:latest
```
2. Run the image as a container
```
docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
```

### APP

To run the app just run the main file like you would another go app.
```
go run main.go
```


## API Documentation

### Authentication

---

#### **POST /register**
Registers a new user.

**Request:**

- Content-Type: `application/json`
- Body:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

**Response:**

- `201 Created` on success
- `400 Bad Request` if input is invalid
- `500 Internal Server Error` on server error

---

#### **POST /login**
Logs in an existing user and returns a JWT token.

**Request:**

- Content-Type: `application/json`
- Body:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

**Response:**

- `200 OK` on success with Bearer token in authorization header:
- `400 Bad Request` if input is invalid
- `500 Internal Server Error` on server error

---

### To-Do Management

---

#### **GET /todo**
Fetches all the to-do items for the authenticated user.

**Request:**

- Authorization: `Bearer <token>`

**Response:**

- `200 OK` with the list of to-dos:
  ```json
  [
    {
      "id": "objectID",
      "title": "string",
      "text": "string",
      "done": "boolean"
    }
  ]
  ```
- `401 Unauthorized` if token isn't valid
- `500 Internal Server Error` on server error

---

#### **POST /todo**
Creates a new to-do item.

**Request:**

- Authorization: `Bearer <token>`
- Content-Type: `application/json`
- Body:
  ```json
  {
    "title": "string",
    "text": "string"
  }
  ```

**Response:**

- `201 Created` with the created to-do:
  ```json
  {
    "id": "objectID",
    "title": "string",
    "text": "string",
    "done": "boolean"
  }
  ```
- `400 Bad request` if input is invalid
- `401 Unauthorized` if token isn't valid
- `500 Internal Server Error` on server error

---

#### **PUT /todo/{id}**
Updates an existing to-do item.
Fields are optional. This method is meant to be used when you want to mark a todo as done.

**Request:**

- Authorization: `Bearer <token>`
- Content-Type: `application/json`
- Body:
  ```json
  {
    "title": "string",
    "text": "string",
    "done": "boolean"
  }
  ```
  - example:
    ```json
    {
    "done": true
    }
    ```

**Response:**

- `200 OK` with the updated to-do
- `400 Bad request` if input is invalid
- `401 Unauthorized` if token isn't valid
- `403 Forbidden` if user doesn't own the todo
- `404 Not Found` if the to-do doesn't exist
- `500 Internal Server Error` on server error

---

#### **DELETE /todo/{id}**
Deletes a to-do item.

**Request:**

- Authorization: `Bearer <token>`

**Response:**

- `200 OK` on success
- `401 Unauthorized` if token isn't valid
- `403 Forbidden` if user doesn't own the todo
- `404 Not Found` if the to-do doesn't exist
- `500 Internal Server Error` on server error

---

#### **PUT /todo/reorder**
Reorders the to-do list based on the new order of IDs.
the IDs in the array need to be the same as the IDs of the user's todo items
just in a different order

**Request:**

- Authorization: `Bearer <token>`
- Content-Type: `application/json`
- Body:
  ```json
  {
    "order": ["objectID1", "objectID2", "objectID3"]
  }
  ```

**Response:**

- `200 OK` on success
- `400 Bad Request` if the input is invalid
- `401 Unauthorized` if token isn't valid
- `500 Internal Server Error` on server error

---

### Category Management

---

#### **GET /category**
Fetches all categories for the authenticated user.

**Request:**

- Authorization: `Bearer <token>`

**Response:**

- `200 OK` with the list of categories:
  ```json
  [
    {
      "name": "string",
      "todos": ["objectID"]
    }
  ]
  ```
- `500 Internal Server Error` on server error

---

#### **GET /category/{name}**
Fetches all the to-do items under a specific category.

**Request:**

- Authorization: `Bearer <token>`

**Response:**

- `200 OK` with the list of todos for that category:
  ```json
  [
    {
      "id": "objectID",
      "title": "string",
      "text": "string",
      "done": "boolean"
    }
  ]
  ```
- `401 Unauthorized` if token isn't valid
- `404 Not Found` if category doesn't exist
- `500 Internal Server Error` on server error
---

#### **POST /category/{name}**
Creates a new category for the user.

**Request:**

- Authorization: `Bearer <token>`
- Content-Type: `application/json`
- Body:
  ```json
  {
    "name": "string"
  }
  ```

**Response:**

- `201 Created` on success
- `401 Unauthorized` if token isn't valid
- `409 Conflict` if the category already exists
- `500 Internal Server Error` on server error

---

#### **DELETE /category/{name}**
Deletes a category from the user's list.

**Request:**

- Authorization: `Bearer <token>`

**Response:**

- `200 OK` on success
- `401 Unauthorized` if token isn't valid
- `404 Not Found` if the category doesn't exist
- `500 Internal Server Error` on server error

---

#### **POST /category/{name}/{id}**
Adds a to-do item to the specified category.

**Request:**

- Authorization: `Bearer <token>`
- Path Parameters:
    - `name`: The name of the category
    - `id`: The ID of the to-do item to be added

**Response:**

- `200 OK` on success (also if it is already in the category)
- `400 Bad request` if id of todo isn't valid
- `401 Unauthorized` if token isn't valid
- `403 Forbidden` if user doesn't own todo
- `404 Not Found` if the category doesn't exist
- `500 Internal Server Error` on server error

---

#### **DELETE /category/{name}/{id}**
Removes a to-do item from the specified category.

**Request:**

- Authorization: `Bearer <token>`
- Path Parameters:
    - `name`: The name of the category
    - `id`: The ID of the to-do item to be removed

**Response:**

- `200 OK` on success
- `400 Bad request` if id of todo isn't valid
- `401 Unauthorized` if token isn't valid
- `403 Forbidden` if user doesn't own todo (also if todo doesn't exist)
- `404 Not Found` if the category doesn't exist
- `500 Internal Server Error` on server error

---

### Share Management

---

#### **POST /share/{username}**
Shares the authenticated user's todos with another user.

**Request:**

- Authorization: `Bearer <token>`
- Path Parameters:
  - `username`: The username of the user to share todos with.

**Response:**

- `200 OK` on success
- `404 Not Found` if the username doesn't exist
- `409 Conflict` if todos have already been shared with this user
- `500 Internal Server Error` on server error

---

#### **DELETE /share/{username}**
Unshares the authenticated user's todos from another user.

**Request:**

- Authorization: `Bearer <token>`
- Path Parameters:
  - `username`: The username of the user to unshare todos with.

**Response:**

- `200 OK` on success
- `404 Not Found` if the username doesn't exist or if todos were not shared
- `500 Internal Server Error` on server error

---

#### **GET /share/{username}**
Fetches todos shared by a specific user with the authenticated user.

**Request:**

- Authorization: `Bearer <token>`
- Path Parameters:
  - `username`: The username of the user who shared their todos.

**Response:**

- `200 OK` with the list of shared todos:
  ```json
  [
    {
      "id": "objectID",
      "title": "string",
      "text": "string",
      "done": "boolean"
    }
  ]
  ```
- `403 Forbidden` if the user hasn't shared their todos
- `404 Not Found` if the username doesn't exist
- `500 Internal Server Error` on server error

---

#### **GET /share**
Fetches a list of users who have shared their todos with the authenticated user.

**Request:**

- Authorization: `Bearer <token>`

**Response:**

- `200 OK` with the list of usernames:
  ```json
  [
    "username1",
    "username2"
  ]
  ```
- `500 Internal Server Error` on server error

---

#### **PUT /share/{id}**
Updates a specific todo that was shared with the authenticated user.

**Request:**

- Authorization: `Bearer <token>`
- Path Parameters:
  - `id`: The ID of the todo to update.
- Content-Type: `application/json`
- Body:
  ```json
  {
    "title": "string",
    "text": "string",
    "done": "boolean"
  }
  ```

**Response:**

- `200 OK` with the updated todo:
  ```json
  {
    "id": "objectID",
    "title": "string",
    "text": "string",
    "done": "boolean"
  }
  ```
- `400 Bad Request` if the body or `id` are invalid
- `403 Forbidden` if the todo wasn't shared with the authenticated user
- `500 Internal Server Error` on server error

---