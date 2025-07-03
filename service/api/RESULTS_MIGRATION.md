# Migration Guide: /api/results Response Update

**Effective Date:** 2025-07-03 (services branch)

This document outlines the changes introduced to the `/api/results` endpoint and provides steps to migrate existing clients.

---

## What Changed?

1. **Map-based payloads**
   
   * `primary` and `drained` are no longer flat arrays.  
     They are now JSON objects (maps) keyed by **heading ID**.

   ```jsonc
   {
     "primary": {
       "42": { "heading_id": 42, "passing_score": 290, "last_admitted_rating_place": 30 }
     },
     "drained": {
       "42": [ /* DrainedResultDTO[] */ ]
     }
   }
   ```

2. **Explicit heading identifier**
   
   * Each element inside both arrays now exposes a new integer field `heading_id` that duplicates the map key for convenience.

3. **No changes** to query parameters or the `steps` object.

---

## Rationale

* Map-based grouping drastically simplifies client-side look-ups by heading and avoids repeated filtering.
* The additional `heading_id` field brings parity with other endpoints and guards against accidental mis-keying when serialising/deserialising.

---

## Impact on Clients

| Area | Old Behaviour | New Behaviour |
|------|---------------|---------------|
| Access primary results | Iterate over `response.primary` slice. | Access per-heading slice via `response.primary[headingID]`. |
| Access drained results | Iterate over `response.drained` slice. | Access per-heading slice via `response.drained[headingID]`. |
| DTO shape | No `heading_id` field. | Contains `heading_id`, `passing_score`, `last_admitted_rating_place`. |

If you rely on array order, update your code to iterate over `response.primary[hid]` and `response.drained[hid]`.

---

## Example (abbreviated)

```jsonc
{
  "steps": { "42": [0, 25] },
  "primary": {
    "42": { "heading_id": 42, "passing_score": 290, "last_admitted_rating_place": 30 }
  },
  "drained": {
    "42": [
      { "heading_id": 42, "drained_percent": 25, "avg_passing_score": 290 }
    ]
  }
}
```

---

## Versioning & Backwards Compatibility

This is a **breaking change**.  Bump your client's expected API version or pin to an older container tag until migration is complete.

---

For questions, please open an issue or reach out on the project chat. 