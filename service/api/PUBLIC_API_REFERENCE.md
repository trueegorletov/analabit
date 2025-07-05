# Analabit Public API

This document describes the current set of REST endpoints exposed by the `service/api` microservice.  All endpoints are prefixed with `/api` as configured in `service/api/main.go`.

> **Status:** _Draft_ — the API surface may evolve.  Please check this document or the OpenAPI specification (to be added) for future changes.

---

## Conventionst

* All endpoints use **HTTPS** (recommended) and the **JSON** media type.
* Successful responses return HTTP status **200 OK** unless specified otherwise.
* Error responses follow Fiber's default JSON error structure and include `code`, `message`, and optional `stack` fields.
* Pagination is implemented with the query parameters `limit` (max-items, default `100`) and `offset` (index of the first item, default `0`).

---

## Endpoints

| Method | Path | Description |
|-------|------|-------------|
| `GET` | `/api/varsities` | List varsities (universities / campuses). |
| `GET` | `/api/headings` | List programme headings (majors) with optional varsity filter. |
| `GET` | `/api/headings/{id}` | Retrieve a single heading by its numeric identifier. |
| `GET` | `/api/applications` | List student applications with optional filtering. |
| `GET` | `/api/students/{id}` | List **all** applications that belong to a particular student. |
| `GET` | `/api/results` | Retrieve calculation (primary) and/or drained results with batching. |

Each endpoint is detailed below.

### 1. List Varsities

`GET /api/varsities`

Retrieve all varsities with optional pagination.

| Query Parameter | Type | Default | Description |
|-----------------|------|---------|-------------|
| `limit` | integer | `100` | Maximum number of results to return. |
| `offset` | integer | `0` | Number of items to skip before starting to collect the result set. |

#### Response `200 OK`

```jsonc
[
  {
    "id": 1,
    "code": "hse",
    "name": "Higher School of Economics"
  }
]
```

#### Possible Errors

| Status | When |
|--------|-------|
| `500 Internal Server Error` | Unexpected database or server failure. |

---

### 2. List Headings

`GET /api/headings`

Retrieve programme headings (study programmes). Pagination and varsity filtering are supported.

| Query Parameter | Type | Default | Description |
|-----------------|------|---------|-------------|
| `limit` | integer | `100` | Maximum number of results. |
| `offset` | integer | `0` | Skip the first *n* items. |
| `varsityCode` | string | — | Filter headings that belong to a varsity with the given `code`. |

#### Response `200 OK`

```jsonc
[
  {
    "id": 42,
    "code": "09.03.03",
    "name": "Applied Mathematics and Computer Science",
    "regular_capacity": 100,
    "target_quota_capacity": 10,
    "dedicated_quota_capacity": 5,
    "special_quota_capacity": 3,
    "varsity": {
      "id": 1,
      "code": "hse",
      "name": "Higher School of Economics"
    }
  }
]
```

#### Possible Errors

| Status | When |
|--------|-------|
| `500 Internal Server Error` | Unexpected server failure. |

---

### 3. Get Heading by ID

`GET /api/headings/{id}`

Retrieve a single heading by its numeric identifier.

| Path Parameter | Type | Description |
|----------------|------|-------------|
| `id` | integer | Unique identifier of the heading. |

#### Response `200 OK`

Returns the same JSON structure as a single element from **List Headings**.

#### Possible Errors

| Status                      | When                                             |
|-----------------------------|--------------------------------------------------|
| `400 Bad Request`           | `id` is not a valid integer.                     |
| `404 Not Found`             | A heading with the provided `id` does not exist. |
| `500 Internal Server Error` | Unexpected server failure.                       |

---

### 4. List Applications

`GET /api/applications`

Retrieve student application records. Supports several filters in addition to pagination.  
Each application object now includes this new set of fields:

* `original_submitted` – `true` if, in the same iteration, the student **submitted original documents** to the varsity this heading belongs to.
* `original_quit` – `true` if `original_submitted` is *false* **and** the student submitted originals to a *different* varsity during the same run.
* `passing_now` – `true` if the student currently passes for exactly this heading in this run.
* `passing_to_more_priority` - `true` if the student currently passes for another heading of the same varsity with a higher priority given to it that to this heading.
* `another_varsities_count` - the number of other varsities where the student has found applications to at the same run.

Note: the `passing_now` and `passing_to_more_priority` fields are never `true` at the same time.

| Query Parameter | Type             | Default  | Description                                                                                                |
|-----------------|------------------|----------|------------------------------------------------------------------------------------------------------------|
| `limit`         | integer          | `100`    | Maximum number of results.                                                                                 |
| `offset`        | integer          | `0`      | Skip the first *n* items.                                                                                  |
| `studentID`     | string           | —        | Return only applications belonging to the student with this identifier.                                    |
| `varsityCode`   | string           | —        | Return only applications whose heading's varsity has this `code`.                                          |
| `headingId`     | integer          | —        | Return only applications for the specified heading.                                                        |
| `run`           | string / integer | `latest` | Specify the run to retrieve applications from. Can be a numeric run ID or a string exactly equal `latest`. |

#### Response `200 OK`

```jsonc
[
  {
    "id": 123,
    "student_id": "ABC-123456",
    "priority": 1,
    "competition_type": 0,
    "rating_place": 47,
    "score": 286,
    "run_id": 1,
    "updated_at": "2024-06-04T08:12:34Z",
    "heading_id": 42,
    "original_submitted": true,
    "original_quit": false,
    "passing_now": true,
    "passing_to_more_priority": false,
    "another_varsities_count": 0,
  }
]
```

#### Possible Errors

| Status | When |
|--------|-------|
| `500 Internal Server Error` | Unexpected server failure. |

---

### 5. Get Student Applications

`GET /api/students/{id}`

Convenience endpoint that returns **all** applications submitted by the specified student.

| Path Parameter | Type | Description |
|----------------|------|-------------|
| `id` | string | Unique student identifier used in the admission system. |

#### Response `200 OK`

Returns an array of Application objects identical to the **List Applications** response.

#### Possible Errors

| Status | When |
|--------|-------|
| `404 Not Found` | No applications were found for the given student. |
| `500 Internal Server Error` | Unexpected server failure. |

---

### 6. Get Calculation & Drained Results

`GET /api/results`

Aggregated endpoint that retrieves primary calculation results (admitted students) and/or drained statistics. Supports batching by multiple headings or varsity.

| Query Parameter | Type | Default | Description |
|-----------------|------|---------|-------------|
| `headingIds` | string (csv) | — | Comma‐separated list of heading IDs to filter by. |
| `varsityCode` | string | — | Filter by all headings of the given varsity code (ignored if `headingIds` supplied). |
| `primary` | string/int | (absent) | If present, include primary results. Accepts `latest` (default when empty) or a specific iteration number. |
| `drained` | string | (absent) | `all` to include every drained-percent result, or a comma‐separated list of drained-percent steps (e.g. `drained=25,50,100`).  When missing, no drained results are returned. |
| `run` | string/int | `latest` | The run to use for both steps discovery and result retrieval. Can be a numeric run ID or a string like "latest" or "latest-1". |

#### Response `200 OK`

Note: each list inside of the `steps` object is a list of drained-percent steps for the given heading,
and it always contains the value `100`.

```jsonc
{
  "steps": {
    "42": [25, 50, 33, 100]
  },
  "primary": {
    "42": {
      "heading_id": 42,
      "heading_code": "09.03.03",
      "passing_score": 295,
      "last_admitted_rating_place": 50,
      "run_id": 7
    }
  },
  "drained": {
    "42": [
      {
        "heading_id": 42,
        "heading_code": "09.03.03",
        "drained_percent": 42,
        "avg_passing_score": 295,
        "min_passing_score": 260,
        "max_passing_score": 310,
        "med_passing_score": 297,
        "avg_last_admitted_rating_place": 50,
        "min_last_admitted_rating_place": 30,
        "max_last_admitted_rating_place": 70,
        "med_last_admitted_rating_place": 48,
        "run_id": 7
      }
    ]
  }
}
```

#### Possible Errors

| Status | When |
|--------|-------|
| `400 Bad Request` | Invalid iteration or drained step value. |
| `500 Internal Server Error` | Database or server failure. |

---

## Entity Reference

### Varsity

| Field | Type | Notes |
|-------|------|-------|
| `id` | integer | Auto-increment primary key. |
| `code` | string | Unique code identifying the varsity. |
| `name` | string | Human-readable name. |

### Heading

| Field | Type |
|-------|------|
| `id` | integer |
| `code` | string |
| `name` | string |
| `regular_capacity` | integer |
| `target_quota_capacity` | integer |
| `dedicated_quota_capacity` | integer |
| `special_quota_capacity` | integer |
| `varsity` | object | Embedded varsity reference (`id`, `code`, `name`). |

### Application

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Auto-increment primary key. |
| `student_id` | string | External identifier of the entrant. |
| `priority` | integer | Rank of the choice (1 = highest). |
| `competition_type` | integer (enum) | Admission track (e.g., BVI, quota, general). |
| `rating_place` | integer | Entrant's position in published ranking for this heading. |
| `score` | integer | Sum of entrance-exam scores or badge points. |
| `iteration` | integer | Sequential upload round for raw data. |
| `updated_at` | string (RFC 3339) | Timestamp when the row was inserted/updated. |
| `heading_id` | integer | FK to Heading entity. |
| `original_submitted` | bool | Entrant has filed original documents to *this* varsity in this iteration. |
| `original_quit` | bool | Entrant filed originals to a **different** varsity in same iteration. |
| `passing_now` | bool | Entrant is within admitted quota (latest calculation iteration). |

