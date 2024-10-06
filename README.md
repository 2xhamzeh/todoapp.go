# ToDo API with Golang

## About:
#### This is a simple api written in go. It helps you manage todo items

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