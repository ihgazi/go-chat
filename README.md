# Go-Chat

## Description

**Go-Chat** is a real-time chat application built with Go, designed for scalability, security, and ease of use. It leverages WebSockets for instant messaging, supports user authentication, and organizes conversations into rooms with membership management. The project follows a clean architecture, separating concerns into internal modules for users and WebSocket communication, and uses a PostgreSQL database for persistent storage.

Key features include:

- **User Authentication:** Secure registration and login with password hashing.
- **Real-Time Messaging:** WebSocket-based communication for fast, bidirectional chat.
- **Room Management:** Create, join, and manage chat rooms with member tracking.
- **Modular Design:** Organized codebase for maintainability and extensibility.
- **Database Migrations:** Versioned SQL migrations for easy schema updates.

## Local Development

To set up and run ChatSphere Go-Chat locally, follow these steps:

### Prerequisites

- Go (1.20+ recommended)
- PostgreSQL

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-chat.git
cd go-chat
```

### 2. Configure Environment Variables

Copy `.env.example` to `.env` and update the values as needed:

```bash
cp .env.example .env
```

Set your database connection string and other secrets in `.env`.

### 3. Set Up the Database

Create a PostgreSQL database and run the migrations:

```bash
migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go-chat?sslmode=disable" -verbose up
```
**NOTE**: You will require https://github.com/golang-migrate/migrate installed to run the above command, or use a different migration tool of your choice.

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Application

You can use the provided script or run manually:

```bash
go run cmd/main.go
```

The server will start and listen on the port specified in your `.env` file.

### 6. API Endpoints

- **User Registration & Login:** `/api/register`, `/api/login`
- **WebSocket Endpoint:** `/ws`
- **Room Management:** `/api/rooms`, `/api/rooms/:id/members`

Refer to the code and router definitions for more details.

### 7. Testing

You can use tools like [Postman](https://www.postman.com/) or [curl](https://curl.se/) to interact with the REST endpoints and WebSocket clients for real-time messaging.

---

