
# Prompting guidance for small and medium models

## Recommendation 1 — Hardcode filenames

The generic prompt in `prompts/prompt.md` uses placeholder names like
`<deployment-template>.template.md` and `<spec-name>.md`. Larger models
handle this without confusion. Smaller models benefit from having the
actual filenames stated explicitly in the prompt so the context is
unambiguous.

See `tools/pcdp-lint/spec/prompt.md` for an example of this pattern.

## Recommendation 2 — Phase the delivery

Long specifications (pcdp-lint.md is 1000+ lines with 27 examples) push
models toward token and rate limits if asked to produce all deliverables
in a single pass. Splitting the run into explicit phases prevents
mid-generation interruptions:

**Phase 1 — Core implementation**
Ask for only `main.go` and `go.mod`.

**Phase 2 — Build and packaging**
Instruct the model to read existing files first, then produce
`Makefile`, RPM spec, Debian packaging files, and `LICENSE`.

**Phase 3 — Documentation and report (always last)**
Ask for `README.md` and `TRANSLATION_REPORT.md` separately.
`TRANSLATION_REPORT.md` is the most token-intensive output (it must
assess every EXAMPLE) and benefits from having all other files visible
in context before it is written.

The `tools/pcdp-lint/spec/prompt.md` prompt encodes this phase
structure and resume logic directly, so the model knows which files
to skip if it is restarting a partial run.

## Recommendation 3 — Resume awareness

If a run is interrupted (rate limit, context overflow, timeout),
restart with an explicit instruction to read the output directory
before doing anything:

```
Read all files currently in the output directory.
Treat any non-empty file as complete — do not overwrite it.
Produce only the missing files, in delivery order.
```

This avoids regenerating already-complete files and keeps
each resumed session within budget.

## max_tokens guidance

- 16384 is sufficient for most individual file generations
- The TRANSLATION_REPORT.md for a spec with 20+ examples may need
  more; consider 32000 for that phase alone if the model truncates
