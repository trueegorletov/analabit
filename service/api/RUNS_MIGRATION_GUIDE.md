# API Migration Guide: From `iteration` to `run`

This document outlines the changes required to adapt to the new `run`-based API for retrieving applications and results. The legacy `iteration` and `runOffset` parameters have been deprecated and replaced by a single `run` parameter.

## Summary of Changes

- **Single Query Parameter:** The `iteration` and `runOffset` query parameters are replaced by a single `run` parameter on all relevant endpoints.
- **Consistent Naming:** The new `run` parameter provides a unified way to specify the desired dataset snapshot.
- **Updated Response Fields:** The `iteration` field has been removed from all API responses and replaced with `run_id`.

## Endpoint-Specific Changes

### 1. List Applications

`GET /api/applications`

- **DEPRECATED:** `?iteration={number}` and `?runOffset={offset}`
- **NEW:** `?run={run_specifier}`

The `run` parameter accepts:
- A specific numeric run ID (e.g., `?run=123`).
- A relative offset string (e.g., `?run=latest` or `?run=latest-1`).

**Example Response Change:**

**Old Response:**
```json
{
  "id": 123,
  "iteration": 1,
  ...
}
```

**New Response:**
```json
{
  "id": 123,
  "run_id": 1,
  ...
}
```

### 2. Get Calculation & Drained Results

`GET /api/results`

- **DEPRECATED:** `?iteration={number}` and `?runOffset={offset}`
- **NEW:** `?run={run_specifier}`

The `run` parameter on this endpoint behaves identically to the one on `/api/applications`.

**Example Response Change:**

**Old Response:**
```json
{
  "primary": {
    "42": {
      "iteration": 7,
      ...
    }
  },
  "drained": {
    "42": [
      {
        "iteration": 7,
        ...
      }
    ]
  }
}
```

**New Response:**
```json
{
  "primary": {
    "42": {
      "run_id": 7,
      ...
    }
  },
  "drained": {
    "42": [
      {
        "run_id": 7,
        ...
      }
    ]
  }
}
```

## Migration Steps

1. **Update API Calls:**
   - In your client-side code, replace any instances of `?iteration=` or `?runOffset=` with `?run=`.
   - Ensure the value passed to `run` matches the format your application requires (e.g., `latest` for the most recent data).

2. **Update Response Handling:**
   - Adjust your data models and parsing logic to expect a `run_id` field instead of `iteration` in the responses from the `/api/applications` and `/api/results` endpoints.

This migration simplifies API interactions and aligns with the new data versioning model. 