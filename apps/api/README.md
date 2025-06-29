# API

Go API with Gin


## Routes

`GET /api/v1/clusters`

Returns basic Cluster info

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


---

## Catchalls

`404 Not Found`

Returns all other routes

```json
{
  "message": "The route you're looking for was not found. Perhaps you wanted one of these?",
  "routes": [
    {
      "method": "GET",
      "path": "/",
      "description": "Returns basic endpoint info."
    },
    {
      "method": "GET",
      "path": "/api/v1/clusters",
      "description": "Returns the cluster IDs along with connectivity information"
    }
  ]
}
```

`405 Method Not Allowed`
Returns **alternative** routes that match the path


```json
{
  "message": "Method 'POST' not allowed for path '/api/v1/clusters'. Did you mean one of these?",
  "routes": [
    {
      "method": "GET",
      "path": "/api/v1/clusters",
      "description": "Returns the cluster IDs along with connectivity information"
    }
  ]
}
```

