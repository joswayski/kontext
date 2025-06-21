## API

Go API to extract data from Kafka clusters

### Endpoints

#### GET /api/v1/clusters

Returns information about configured Kafka clusters.

**Response:**
```json
[
  {
    "id": "cluster1",
    "bootstrap_servers": "localhost:9092",
    "status": "connected",
    "error?": "optional error message"
  }
]
```

**Response Fields:**
- `id` (string): Unique identifier for the cluster
- `bootstrap_servers` (string): Comma-separated list of bootstrap servers
- `status` (string): Connection status of the cluster
- `error` (string, optional): Error message if connection failed (omitted if empty)

#### GET /health & /api/v1/health 

Health check endpoint to verify the API is running.

#### GET / & /api/v1

Root endpoint that returns basic API information.
