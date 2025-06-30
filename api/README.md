# API

Go API with Gin

## Endpoints

### GET /api/v1/clusters
Returns cluster connectivity information.

**Response:**
```json
{
  "analytics": {
    "status": "connected",
    "message": "Saul Goodman"
  },
  "production": {
    "status": "connected",
    "message": "Saul Goodman"
  }
}
```

### GET /
Returns basic endpoint info.

## Error Handling

**404 Not Found** - Returns available routes when a path doesn't exist
**405 Method Not Allowed** - Returns alternative routes for the requested path

Both error responses include a `routes` array with available endpoints and their descriptions.
