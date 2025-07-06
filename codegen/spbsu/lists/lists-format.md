# SPbSU Application Lists JSON Format Analysis

## 1. Introduction

This document analyzes the structure of sample application lists from SPbSU (Saint Petersburg State University), specifically the JSON responses from `enrollelists.spbu.ru/lists?id=...&without_control=true`, and details how to extract the information required for the `ApplicationData` struct used in the project. The analysis is based on two sample files: `list-1407.json` and `list-630.json`.

## 2. JSON Structure Overview

Both files are top-level JSON objects with at least the following keys:
- `list`: an array of application entries (each is a dictionary/object)
- `trials`: (present in some lists) an array of arrays of trial names (subjects)
- `have_places`: (present in some lists) a number (not relevant for ApplicationData)

For the purpose of populating `ApplicationData`, only the `list` array is relevant. Each element in `list` represents a single application.

## 3. ApplicationData Struct

The struct is defined as follows:

```go
// ApplicationData mirrors the data needed for core.VarsityCalculator.AddApplication
// and adds additional fields which could be useful for further processing.
type ApplicationData struct {
    HeadingCode       string
    StudentID         string
    ScoresSum         int
    RatingPlace       int
    Priority          int
    CompetitionType   core.Competition
    OriginalSubmitted bool
}
```

Each field is intended to capture a specific aspect of an application, such as the program applied to, the applicant's unique ID, their total score, their rank, their application priority, the type of competition/quota, and whether the original document was submitted.

## 4. Field-by-Field Mapping

| ApplicationData field | JSON field (in `list` entry) | Mapping/Logic |
|-----------------------|-------------------------------|--------------|
| StudentID             | `user_code`                   | Use as string |
| ScoresSum             | `score_overall`               | Use as int |
| RatingPlace           | `order_number`                | Use as int (1-based rank) |
| Priority              | `priority_number`             | Use as int |
| OriginalSubmitted     | `original_document`           | Use as bool |

### CompetitionType Logic
- If `target_organization` is not null: **Target quota**
- Else if `without_trials` is true: **Special competition** (e.g., olympiad winner)
- Else: **General competition**

This logic should be mapped to the appropriate values of your `core.Competition` enum.
## 6. Recommendations and Future-Proofing

- **Modularity**: Structure the extraction logic so that it is easy to extend if new fields are needed (e.g., exam scores, achievements, or status fields).
- **Documentation**: Clearly document any assumptions, such as the uniqueness of `user_code` or the presence of required fields.
- **Validation**: Consider adding validation steps to check for missing or malformed data before populating `ApplicationData`.
- **Extensibility**: If future requirements demand more detailed analytics (e.g., per-subject scores), the logic should be easy to adapt.

## 7. Summary Table

| JSON field           | ApplicationData field | Type/Logic                                      |
|----------------------|----------------------|-------------------------------------------------|
| user_code            | StudentID            | string                                          |
| score_overall        | ScoresSum            | int                                             |
| order_number         | RatingPlace          | int                                             |
| priority_number      | Priority             | int                                             |
| original_document    | OriginalSubmitted    | bool                                            |



## 8. Conclusion

The mapping between the SPbSU application lists' JSON structure and the `ApplicationData` struct is clear and robust, provided that the extraction logic accounts for possible missing or null values and prioritizes competition type logic correctly. The approach outlined here should serve as a reliable foundation for further development and data processing.
