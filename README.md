# Ascend API

Welcome to the backend server for **Ascend**, a gamified personal development application designed to help users "level up" in real life. This API powers the core functionalities, including user authentication, profile management, and the dynamic generation of daily and weekly quests using Google's Gemini AI.

This repository contains the complete Go source code for the API, built to be scalable, secure, and ready for deployment.

---

## ✨ Features

-   **Secure Authentication**: JWT-based authentication for user registration and login.
-   **Player Profile Management**: Full CRUD operations for player stats (Strength, Agility, Intelligence, etc.).
-   **Dynamic Quest System**: Leverages the Gemini API to generate personalized daily and weekly quests based on user goals.
-   **Gamified Progression**: Users gain XP and level up by completing quests, with stats updated accordingly.
-   **RESTful Architecture**: Clean, organized, and easy-to-understand API endpoints.
-   **Ready for Deployment**: Dockerized for consistent, one-command deployments on platforms like Render.

---

## 🛠️ Tech Stack

-   **Language**: [Go](https://golang.org/) (v1.23+)
-   **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
-   **Database ORM**: [GORM](https://gorm.io/)
-   **Database**: [PostgreSQL](https://www.postgresql.org/)
-   **AI Integration**: [Google Gemini API](https://ai.google.dev/)
-   **Authentication**: [JWT (JSON Web Tokens)](https://jwt.io/)
-   **Containerization**: [Docker](https://www.docker.com/)

---

## 🚀 Getting Started

Follow these instructions to get a local copy of the server up and running for development and testing purposes.

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.23 or higher)
-   [PostgreSQL](https://www.postgresql.org/download/) installed and running
-   [Docker](https://www.docker.com/get-started) (Optional, for containerized development)

### Installation & Setup

1.  **Clone the repository:**
    ```sh
    git clone [https://github.com/rajvirsingh2/ascend-api.git](https://github.com/rajvirsingh2/ascend-api.git)
    cd ascend-api
    ```

2.  **Install Go dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Set up Environment Variables:**
    Create a `.env` file in the root of the project by copying the example file:
    ```sh
    cp .env.example .env
    ```
    Now, open the `.env` file and fill in your specific credentials:

    ```env
    # .env
    
    # Full connection string for your PostgreSQL database
    # Example for local setup: "host=localhost user=postgres password=yourpassword dbname=ascend_db port=5432 sslmode=disable"
    DATABASE_URL="your_database_connection_string"
    
    # A strong, random secret key for signing JWTs
    JWT_SECRET="your_jwt_secret_key"
    
    # Your API key from Google AI Studio for the Gemini API
    GEMINI_API_KEY="your_gemini_api_key"
    
    # The port the server will run on
    API_PORT="8000"
    ```

4.  **Run the server:**
    ```sh
    go run main.go
    ```
    The server should now be running on the port specified in your `.env` file (e.g., `http://localhost:8000`).

---

## API Endpoints

The API is structured into public (authentication) and private (protected) routes.

| Method | Endpoint                        | Protection | Description                                       |
| :----- | :------------------------------ | :--------- | :------------------------------------------------ |
| `POST` | `/auth/register`                | Public     | Registers a new user and creates their profile.   |
| `POST` | `/auth/login`                   | Public     | Logs in a user and returns a JWT token.           |
| `GET`    | `/api/v1/profile`               | Private    | Fetches the authenticated user's profile and stats. |
| `POST`   | `/api/v1/quests/generate`       | Private    | Generates initial quests for a new user.          |
| `GET`    | `/api/v1/quests`                | Private    | Retrieves all active (non-completed) quests.      |
| `POST`   | `/api/v1/quests/:id/complete`   | Private    | Marks a specific quest as complete and updates XP.  |

*Note: Private routes require a valid JWT in the `Authorization: Bearer <token>` header.*

---

## ☁️ Deployment

This application is configured for easy deployment using Docker.

1.  **Build the Docker Image:**
    ```sh
    docker build -t ascend-api .
    ```

2.  **Run the Container:**
    ```sh
    docker run -p 8080:8000 --env-file .env ascend-api
    ```

For production, it is recommended to deploy to a service like **Render**. Simply connect your GitHub repository, select the Docker runtime, and add your environment variables in the Render dashboard. The included `Dockerfile` will be used to build and deploy the service automatically.

---

## 📁 Project Structure

````

/ascend-api
├── ai/                \# AI integration logic (Gemini adapter)
├── config/            \# Database connection and initial setup
├── controller/        \# Gin handlers for API routes (business logic)
├── middleware/        \# Authentication middleware (RequireAuth)
├── models/            \# GORM database models (User, PlayerProfile, Quest)
├── .env.example       \# Example environment file
├── Dockerfile         \# Docker configuration for deployment
├── go.mod             \# Go module dependencies
├── go.sum             \# Go module checksums
└── main.go            \# Main application entry point, route setup
