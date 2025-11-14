# Web Note Work - Backend

This repository contains the backend service for the "Web Note Work" application. It is a RESTful API built with Go and the Gin framework, designed to manage tasks. The service connects to a MongoDB database for data persistence.

## Core Technologies

*   **Go**: The primary programming language.
*   **Gin**: A high-performance HTTP web framework.
*   **MongoDB**: The NoSQL database used to store task data.
*   **Official Go Driver for MongoDB**: For database connectivity and operations.

## Features

*   **CRUD Operations**: Full Create, Read, Update, and Delete functionality for tasks.
*   **Task Filtering**: Retrieve tasks based on timeframes: `today`, `week`, `month`, or `all`.
*   **Task Analytics**: Provides counts of `active` and `complete` tasks alongside the task list.
*   **CORS Configuration**: Pre-configured to allow requests from specific frontend origins.

## Getting Started

Follow these instructions to get the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   Go (version 1.25 or newer)
*   A running MongoDB instance

### Installation and Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/dangdinh2405/web-note-work-backend.git
    ```

2.  **Navigate to the project directory:**
    ```sh
    cd web-note-work-backend
    ```

3.  **Create a `.env` file** in the root of the project and add the following environment variables. The `main.go` file requires these to connect to the database and run the server.

    ```env
    # Your MongoDB connection string
    MONGO_CONNECTION="mongodb://user:password@host:port"

    # The name of the database to use
    MONGO_DB_NAME="web-note-work"

    # The port for the server to listen on (defaults to 10000)
    PORT="10000"
    ```

4.  **Run the application:**
    ```sh
    go run cmd/api/main.go
    ```

The server will start and be available at `http://localhost:10000`.

## API Endpoints

The API provides the following endpoints under the `/tasks` group.

### Get All Tasks

Retrieves a list of tasks with filtering options. Also returns counts of active and completed tasks.

*   **Endpoint**: `GET /tasks`
*   **Query Parameter**:
    *   `filter` (optional): `today` (default), `week`, `month`, `all`.
*   **Success Response (200 OK)**:
    ```json
    {
        "tasks": [
            {
                "id": "67b9d7e5f9d1e4c3b2a1b0c9",
                "title": "My first task",
                "status": "active",
                "completedAt": null,
                "createdAt": "2024-08-15T10:00:00Z",
                "updatedAt": "2024-08-15T10:00:00Z"
            }
        ],
        "activeCount": 1,
        "completeCount": 0
    }
    ```

### Create a Task

Creates a new task.

*   **Endpoint**: `POST /tasks`
*   **Request Body**:
    ```json
    {
        "title": "A new task to complete"
    }
    ```
*   **Success Response (201 Created)**:
    ```json
    {
        "_id": "67b9d8e5f9d1e4c3b2a1b0d0",
        "title": "A new task to complete",
        "status": "active",
        "completedAt": null,
        "createdAt": "2024-08-15T11:00:00Z",
        "updatedAt": "2024-08-15T11:00:00Z"
    }
    ```

### Update a Task

Updates the title, status, or completion date of an existing task.

*   **Endpoint**: `PUT /tasks/:id`
*   **Path Parameter**:
    *   `:id`: The ObjectID of the task to update.
*   **Request Body** (include only the fields to be updated):
    ```json
    {
        "title": "Updated task title",
        "status": "complete",
        "completedAt": "2024-08-15T12:00:00Z"
    }
    ```
*   **Success Response (200 OK)**:
    ```json
    {
        "id": "67b9d8e5f9d1e4c3b2a1b0d0",
        "title": "Updated task title",
        "status": "complete",
        "completedAt": "2024-08-15T12:00:00Z",
        "createdAt": "2024-08-15T11:00:00Z",
        "updatedAt": "2024-08-15T12:00:00Z"
    }
    ```

### Delete a Task

Deletes a task by its ID.

*   **Endpoint**: `DELETE /tasks/:id`
*   **Path Parameter**:
    *   `:id`: The ObjectID of the task to delete.
*   **Success Response (200 OK)**: Returns the task object that was deleted.
    ```json
    {
        "id": "67b9d8e5f9d1e4c3b2a1b0d0",
        "title": "Updated task title",
        "status": "complete",
        "completedAt": "2024-08-15T12:00:00Z",
        "createdAt": "2024-08-15T11:00:00Z",
        "updatedAt": "2024-08-15T12:00:00Z"
    }
    ```

## Project Structure

```
.
├── cmd/api/
│   └── main.go         # Application entry point, server setup, and CORS config.
├── internal/
│   ├── data/
│   │   └── db.go       # MongoDB connection logic.
│   ├── handler/
│   │   └── tasksHandler.go # Request handlers for all task-related endpoints.
│   ├── http/
│   │   └── router.go   # API route definitions.
│   └── models/
│       └── Task.go     # Data structure for the Task model.
├── go.mod              # Go module definitions.
├── go.sum              # Go module checksums.
└── .env.example        # Example environment variables file.
