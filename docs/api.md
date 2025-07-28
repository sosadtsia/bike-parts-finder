# Bike Parts Finder API

This document outlines the API endpoints for the Bike Parts Finder application.

## Base URL

```
http://localhost:8080/api/v1
```

## Endpoints

### Get All Parts

```
GET /parts
```

Retrieves a paginated list of bike parts.

**Query Parameters:**
- `page` (integer, optional): Page number for pagination. Default: 1
- `limit` (integer, optional): Number of items per page. Default: 20, Maximum: 50

**Response:**
```json
[
  {
    "id": "part-1",
    "brand": "Shimano",
    "model": "XT Brake Set",
    "category": "brakes",
    "sub_category": "hydraulic disc",
    "price": 129.99,
    "msrp": 149.99,
    "currency": "USD",
    "in_stock": true,
    "description": "High performance hydraulic disc brake set with excellent modulation and stopping power.",
    "images": ["https://example.com/image1.jpg"],
    "url": "https://example.com/product/xt-brakes"
  },
  // ...more parts
]
```

### Get Part by ID

```
GET /parts/{id}
```

Retrieves a specific part by ID.

**Path Parameters:**
- `id` (string, required): The unique identifier for the part

**Response:**
```json
{
  "id": "part-1",
  "brand": "Shimano",
  "model": "XT Brake Set",
  "category": "brakes",
  "sub_category": "hydraulic disc",
  "price": 129.99,
  "msrp": 149.99,
  "currency": "USD",
  "in_stock": true,
  "description": "High performance hydraulic disc brake set with excellent modulation and stopping power.",
  "images": ["https://example.com/image1.jpg"],
  "url": "https://example.com/product/xt-brakes"
}
```

### Search Parts

```
GET /parts/search
```

Searches for bike parts based on query parameters.

**Query Parameters:**
- `q` (string, optional): Search query to match against brand, model, or description
- `category` (string, optional): Filter parts by category
- `page` (integer, optional): Page number for pagination. Default: 1
- `limit` (integer, optional): Number of items per page. Default: 20, Maximum: 50

**Response:**
```json
[
  {
    "id": "part-1",
    "brand": "Shimano",
    "model": "XT Brake Set",
    "category": "brakes",
    "sub_category": "hydraulic disc",
    "price": 129.99,
    "msrp": 149.99,
    "currency": "USD",
    "in_stock": true,
    "description": "High performance hydraulic disc brake set with excellent modulation and stopping power.",
    "images": ["https://example.com/image1.jpg"],
    "url": "https://example.com/product/xt-brakes"
  },
  // ...more parts
]
```

## Health Check Endpoints

### Basic Health Check

```
GET /health
```

Returns OK if the service is up and running.

**Response:**
```
OK
```

### Readiness Check

```
GET /health/ready
```

Checks if the service and its dependencies (database, cache) are ready to handle requests.

**Response:**
```
Ready
```

## Error Responses

All endpoints may return the following error responses:

**400 Bad Request**
```json
{
  "error": "bad_request",
  "message": "Description of the error",
  "code": 400
}
```

**404 Not Found**
```json
{
  "error": "not_found",
  "message": "Resource not found",
  "code": 404
}
```

**500 Internal Server Error**
```json
{
  "error": "internal_server_error",
  "message": "An unexpected error occurred",
  "code": 500
}
```
