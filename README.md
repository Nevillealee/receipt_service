# Receipt Processing Service

This is a Go-based receipt processing service that evaluates receipts and awards points based on various rules. The service provides two endpoints:

- `POST /receipts/process`: Process a receipt and return a unique ID for the receipt.
- `GET /receipts/{id}/points`: Retrieve the points awarded to a specific receipt.

## Prerequisites

Ensure you have Docker installed on your system.

## Running the Application

To build and run the application using Docker, follow these steps:

### Step 1: Build the Docker Image

Use the following command to build the Docker image:

```bash
docker build -t receipt-service .
docker run -p 8080:8080 receipt-service
```

This will start the service on port 8080.

## Accessing the API

You can now access the API through `http://localhost:8080`:
