# Improved Meta-Prompt: Instructions for Generating Plan-Prompts to Add University Support

## Purpose of This Meta-Prompt
This meta-prompt provides detailed instructions for you (the plan-prompt creator agent) to generate a comprehensive plan-prompt for adding support for a new university's competitive lists in the admission analytics system. The plan-prompt will guide the Executor Agent in implementing the necessary code. When using this meta-prompt, you will be provided with specific university details (e.g., name, code, root URL, sample files, notes). Your task is to analyze the data format thoroughly and produce a plan-prompt that includes step-by-step implementation details, based on the requirements below.

Key principles:
- Focus on creating an HTTPHeadingSource as the primary implementation & generator script (see existing generator scripts for details) to simply output ready-to-paste headings definitions.
- Perform comprehensive format analysis yourself using tools like read_file, fetch, grep_search, and sequential-thinking heavily.
- Include detailed analysis results, working extraction examples, and parsing logic in the plan-prompt.
- Be cautious with large files to avoid context overflow: use grep_search for patterns, limit read_file to small ranges (e.g., first 100 lines), use fetch with max_length parameter.
- Use only English variable names in code; avoid mixtures like totalКЦП.
- The final plan-prompt should be self-contained, detailed, and executable by the Executor Agent without unnecessary steps like updating README, docs, doing any extra work or achieving specific test coverage.

<ImportantNote>
IMPORTANT: Whenever you do anything related to investigantion while preparing for plan creation, analyse in detail and with lots of attention, mindfully,
the EXISTING sources implementations for another varsities in the codebase, and try to UNDERSTAND how they are already implemented and what each existing sources
implementations part cares about (for currently reviewed one implementation). Understand & study it deeply before creating any single word of the plan, and prioritize
doing it all by analogy to the existing code in codebase over ANY instructions in metaprompt (if they conflict or especially when you are not sure 100% about something).

Tell the Executor Agent to act in this way too (prioritizing analogy & specific for the WIP-university observations over the plan instructions, when they conflict or are not 100% complete / contain vague moments)
<ImportantNote>

## Input Parameters for Plan-Prompt Generation
When generating the plan-prompt, use these parameters provided in the query:
- UNIV_NAME: Full name of the university (e.g., Университет ИТМО).
- UNIV_CODE: Short Latin code (e.g., itmo) for package and folder names.
- LISTS_ROOT_URL: URL to the page listing all competitive lists or a representative list URL.
- SAMPLE_FILES: Paths to local sample files in sample_data/UNIV_CODE/ (e.g., HTML, PDF, JSON, XLSX).
- NOTES: Any additional hints (e.g., capacity parsing details, URL patterns).

## Detailed Glossary of Terms and Fields
Understand these terms deeply as they are universal across universities. Use them to guide your analysis and ensure the plan-prompt explains how to map university-specific representations to these concepts. Provide examples in the plan-prompt where possible.

- КЦП (Control Numbers of Admission, sometimes ККП): The total number of budget places available for a heading (educational program). It is the sum of all capacities: КЦП = Regular + TargetQuota + DedicatedQuota + SpecialQuota. This is crucial for calculations; if not directly available, compute it from quotas or specify manually.

- Regular (General Competition, ОК): The number of budget places for general applicants, including those with BVI rights. Applicants with CompetitionType = Regular or BVI compete for these places. BVI applicants rank higher than Regular ones in shared lists. core.Capacities.Regular stores this value. Note: Regular includes BVI; do not separate them unless the university does.

- TargetQuota (ЦК, целевой набор/приём): Budget places for target enrollment applicants (CompetitionType = TargetQuota). These are often in separate lists. core.Capacities.TargetQuota.

- DedicatedQuota (ОтК): Budget places for dedicated quota applicants (CompetitionType = DedicatedQuota). core.Capacities.DedicatedQuota.

- SpecialQuota (ОсК, for persons with special rights): Budget places for special quota applicants (CompetitionType = SpecialQuota). core.Capacities.SpecialQuota.

- CompetitionType: Enum for competition categories. Types: Regular (general), BVI (without entrance exams, shares list with Regular but ranks higher), TargetQuota, DedicatedQuota, SpecialQuota. Analyze how the university represents these (e.g., separate files, sections, or JSON fields). BVI often has highlighted rows or 'Да' in a BVI field.

- RatingPlace: Applicant's position in the list (unique per list or competition type). Lists are sorted by this; BVI are at the top, followed by Regular sorted by ScoresSum descending.

- Priority: Integer (1-15+) indicating preference for this program (1 = highest). Parse carefully; default to 1 if missing.

- ScoresSum: Total competitive score (e.g., sum of exams + individual achievements). Ranges like 195-310 for 3 exams. May be zero or hidden for BVI.

- StudentID: Unique identifier (e.g., 7-digit code like 4272036, or SNILS). Critical for matching applicants.

- OriginalSubmitted: Boolean (e.g., 'Да'/'Нет') indicating certificate original/consent submission. IMPORTANT: Sometimes it can be called "Согласие", sometimes it is called "Оригинал". If there are both, prefer the one column/field/etc, which has likely more different values in average for lists. Like, if Foo University format
has both this fields, but "Оригинал" contains always true & 

- Heading: An educational program. Each HeadingSource represents one Heading.

During analysis, map university-specific terms/structures to these (e.g., if a university uses 'Общий конкурс' for Regular + BVI, parse accordingly).

## Architecture and Implementation Requirements
- Package: core/source/UNIV_CODE (e.g., core/source/itmo).
- Primary Type: HTTPHeadingSource implementing HeadingSource interface. It fetches from URLs and parses data into HeadingData and ApplicationData.
- File Structure: common.go (parsing logic), http.go (HTTPHeadingSource), serialize.go (gob registration). Add file.go only if needed for debugging complex formats.
- HeadingCode: Computed via utils.GenerateHeadingCode(PrettyName).
- Capacities: Parse from lists if possible; otherwise, specify in HeadingSource.
- Generator Script: `codegen/UNIV_CODE/main.go` should output Go source code representing a slice of `&UNIV_CODE.HTTPHeadingSource{...}` structs. This output should be ready to be manually pasted into the `[]source.HeadingSource{...}` slice in the university's registry file. Do not serialize to JSON or any other format; the output must be valid Go code.
**Note on Data Fetching**: The `codegen` script should prefer using local `sample_data` files for parsing registries and capacities over fetching them from the web. This ensures stability and reproducibility. Only fetch from the web if local samples are unavailable or explicitly instructed.
- Registry: Add core/registry/UNIV_CODE/UNIV_CODE.go with sourcesList() returning the slice. Update central registry.

**CRITICAL ARCHITECTURAL NOTE:** The `codegen` script is responsible for ALL discovery and registry parsing. The `HTTPHeadingSource` struct MUST be self-contained and hold direct URLs to the final application lists for each competition type (e.g., `RegularBVIListURL`, `TargetQuotaListURLs`). The runtime `LoadTo` method should NEVER parse registries or discover URLs. Its only job is to fetch and parse the final application lists from the pre-resolved URLs provided in its struct fields.

## Instructions for Plan-Prompt Creator (You)
1. **Format Analysis**: Before drafting the plan-prompt, analyze the data format comprehensively. Use sequential-thinking heavily for step-by-step reasoning. Examine SAMPLE_FILES and LISTS_ROOT_URL using tools:
   - read_file with limited lines (e.g., start=1, end=100) to avoid overflow.
   - grep_search for patterns (e.g., query for 'Количество мест' to find capacities).
   - fetch with max_length (e.g., 5000) for web pages.
   - If files are large, NEVER read entirely; search incrementally.
   Finish analysis only when you have COMPLETE understanding: identify structures for each field, CompetitionType separation, edge cases (e.g., BVI handling).

2. **Include in Plan-Prompt**: Detailed analysis results, CSS selectors/regex/examples for extraction, parsing logic pseudocode/brief Golang snippets. Provide working examples (e.g., 'From this HTML snippet, extract StudentID as...').

3. **Handle Large Formats**: Warn Executor to be cautious; recommend tool usage as above if further verification needed.

4. **Generate Plan-Prompt Structure**: Use the outline below, filling with university-specific details.

## Step-by-Step Outline for the Plan-Prompt
The generated plan-prompt should follow this structure:

1. **Executive Summary**: Concise overview of additions (e.g., new package, generator, registry hook).

2. **Input Parameters**: List provided values.

3. **Glossary Recap**: Brief reminder of key terms, tailored with university examples from your analysis.

4. **Data Format Analysis**: Your detailed findings (structures, extraction methods, examples).

5. **Step-by-Step Instructions**:
Ensure that your plan's meaningful part is structured using reasonable separation into steps.

6. **Deliverables**: List files/packages created.

7. **Notes & Pitfalls**: University-specific warnings, large file handling.

## Final Reminders
- Invoke sequential-thinking for complex parts.
- Ensure plan-prompt enables quick, correct implementation.
- Be specific and detailed, don't thrift tokens, tell the Executor Agent to analyse existing (for another varsities) sources implementations before proceeding to the work & also analyse (on his own, not 100% trusting the generated plan) the sample files, links or anything like that proceeding to the wor too.