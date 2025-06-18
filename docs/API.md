# Post-Quantum Cognitive Decoys API Documentation

## Backend API Endpoints

### Health Check
```
GET /api/health
```

Returns the health status of the API.

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2023-03-27T15:04:05Z"
}
```

### Status
```
GET /api/status
```

Returns detailed status information about the backend service.

**Response**:
```json
{
  "service": "pqcd-backend",
  "status": "operational",
  "version": "0.1.0",
  "timestamp": "2023-03-27T15:04:05Z"
}
```

### Generate Decoys
```
POST /api/decoys/generate
```

Generates cognitive decoys for a target string.

**Request Body**:
```json
{
  "target": "kyber768",
  "complexity": 5,
  "count": 10
}
```

| Parameter | Type | Description |
|-----------|------|-------------|
| target | string | Required. The target string for which to generate decoys. |
| complexity | integer | Optional. How different the decoys should be from the target (1-10). Default: 5. |
| count | integer | Optional. Number of decoys to generate. Default: 10. |

**Response**:
```json
{
  "decoys": [
    "kyber768",
    "kyber467",
    "kyber788",
    "kyber760",
    "kymir768"
  ],
  "timestamp": "2023-03-27T15:04:05Z",
  "metrics": {
    "generation_time_ms": 150,
    "complexity_level": 5
  }
}
```

## AI Service API Endpoints

### Health Check
```
GET /health
```

Returns the health status of the AI service.

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2023-03-27T15:04:05Z"
}
```

### Generate Decoys
```
POST /generate
```

Generates cognitive decoys using the AI model.

**Request Body**:
```json
{
  "target": "kyber768",
  "complexity": 5,
  "count": 10
}
```

| Parameter | Type | Description |
|-----------|------|-------------|
| target | string | Required. The target string for which to generate decoys. |
| complexity | integer | Optional. How different the decoys should be from the target (1-10). Default: 5. |
| count | integer | Optional. Number of decoys to generate. Default: 10. |

**Response**:
```json
{
  "decoys": [
    "kyber768",
    "kyber467",
    "kyber788",
    "kyber760",
    "kymir768"
  ],
  "target": "kyber768",
  "complexity": 5,
  "count": 5,
  "metrics": {
    "generation_time_ms": 42
  }
}
``` 