# AI Interview Backend

## Why we do this

The AI Interview Backend is designed to provide a robust and scalable backend service for an AI-driven interview platform. This platform aims to streamline the interview process by leveraging AI to evaluate candidates, provide feedback, and improve the overall hiring experience.

## How we do this

The backend is built using Go (Golang), a powerful and efficient programming language well-suited for backend development. It follows a microservices architecture to ensure scalability and maintainability. The API layer uses the [chi](https://github.com/go-chi/chi) router for flexible, RESTful routing and path parameter extraction. Handlers use well-defined DTOs, robust JSON parsing, validation, and consistent error handling. Logging middleware is included for observability. Comprehensive tests cover handlers, routing, and middleware.

## OpenAPI Specification

The API documentation for this project is defined in the `docs/openapi.yaml` file. You can view and interact with the API using tools like [Swagger Editor](https://editor.swagger.io/).

To load the OpenAPI specification:

1. Download the `docs/openapi.yaml` file from this repository.
2. Open [Swagger Editor](https://editor.swagger.io/).
3. Import the `docs/openapi.yaml` file to explore the API interactively.

## Detailed Architecture

For a more detailed explanation of the architecture, please refer to the README file located in the `architecture/` folder.

## Folder Structure

The project is organized as follows:

- `api/`: Contains API route definitions (using chi router), handlers, middleware (e.g., logging), DTOs, and comprehensive tests.
- `business/`: Contains business logic and service implementations.
- `data/`: Contains database models and data access logic.
- `ai/`: Contains logic for interacting with AI models.
- `docs/`: Contains documentation, including the `openapi.yaml` file.
- `architecture/`: Contains detailed architecture documentation.

This structure ensures modularity and maintainability.

## Development & Testing

- To run all API tests:

  ```powershell
  go test ./api/...
  ```

- Continuous Integration (CI) runs lint, build, and test workflows on every push and pull request (see `.github/workflows/`).

- The API is fully covered by tests for handlers, routing, and middleware, ensuring robust and maintainable code.

## How to start with your frontend

1. Clone this repository to your local machine:

   ```bash
   git clone <repository-url>
   ```

2. Navigate to the project directory:

   ```bash
   cd AI_Interview_Backend
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. Set up PostgreSQL:

   - Ensure PostgreSQL is installed and running on your system.
   - Create a database for the project:

     ```sql
     CREATE DATABASE ai_interview;
     ```

   - Update the environment variables to include the PostgreSQL connection string. For example:

     ```bash
     export DATABASE_URL=postgres://username:password@localhost:5432/ai_interview
     ```

5. Start the backend server:

   ```bash
   go run main.go
   ```

6. Ensure the backend server is running on the specified port (default: `http://localhost:8080`).

7. Connect your frontend application to the backend by configuring the API base URL in your frontend project settings.
