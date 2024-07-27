# URL Shortener

A simple URL shortener service that allows users to shorten URLs, retrieve the original URLs, and track metrics about the most frequently shortened domains. This service provides a REST API for interacting with the URL shortening and redirection functionalities.

## Features

- **Shorten URLs:** Convert long URLs into shortened versions.
- **Redirect URLs:** Redirect shortened URLs to the original long URLs.
- **Track Metrics:** Collect and retrieve metrics on the top domains being shortened.

## Architectural Overview

### User Interaction

- **REST API:** Core interaction point where users can submit URLs for shortening and use shortened URLs for redirection.

### Core Components

1. **URL Shortening Logic:** Generates unique shortened URL identifiers.
2. **Redirection Logic:** Handles redirection from shortened URLs to original URLs.
3. **Data Storage:** Stores URL mappings. Initially in-memory; can be replaced with Redis or another database.
4. **Metrics Collection:** Tracks and provides metrics such as the top domains shortened.

### Database/Storage Layer

- **In-Memory Storage:** For simplicity during initial implementation. Consider replacing with a persistent store like Redis for production use.

### Deployment

- **Docker Container:** Containerize the application for consistent deployment.

## REST API Endpoints

### `POST /shorten`

- **Input:** JSON payload with the original URL.
- **Output:** JSON response containing the shortened URL.
- **Process:**
  - Receives and decodes the URL.
  - Generates a short identifier if the URL isn't already shortened.
  - Stores the original URL and its shortened version.

### `GET /redirect/{shortURL}`

- **Input:** Shortened URL path.
- **Process:**
  - Retrieves the original URL for the given short URL.
  - Redirects the user to the original URL.
- **Output:** HTTP redirection to the original URL.

### `GET /metrics`

- **Output:** JSON response with the top 3 most frequently shortened domains and their counts.
- **Process:**
  - Extracts and tracks domain usage.
  - Retrieves and sorts domain usage to provide metrics.

## Implementation Details

### URL Shortening Logic

- **Unique Identifier Generation:** Creates a unique short string (e.g., 6 characters) to represent the URL.
- **Collision Handling:** Ensures that the generated identifier is unique.

### Storage Mechanism

- **In-Memory Store:** Uses a map for URL mappings and another map for domain usage.
- **Concurrency Considerations:** Uses synchronization mechanisms to handle concurrent operations.

### Metrics Collection

- **Domain Extraction:** Extracts the domain from URLs to track usage.
- **Top Domains Calculation:** Retrieves and sorts the top 3 domains based on the number of shortened URLs.

### Dockerization

- **Dockerfile:** Defines the application environment, including the base image, code copying, and dependency setup.
- **Building and Deployment:** Ensure consistent deployment across different environments.

## Running the Application

1. **Build the Application:**
   ```bash
   go build -o url-shortener
   ```

2. **Run the Application:**
   ```bash
   ./url-shortener
   ```

3. **Access the API:**
   - Shorten URL: `POST http://localhost:8080/shorten`
   - Redirect URL: `GET http://localhost:8080/redirect/{shortURL}`
   - Metrics: `GET http://localhost:8080/metrics`
   

## Dockerization

To build and run the Docker container:

1. **Build Docker Image:**
   ```bash
   docker build -t url-shortener .
   ```

2. **Run Docker Container:**
   ```bash
   docker run -p 8080:8080 url-shortener
   ```