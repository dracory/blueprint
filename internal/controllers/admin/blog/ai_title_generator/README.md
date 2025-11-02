# AI Title Generator

The AI Title Generator module brings the end-to-end workflow for creating, 
curating, and promoting AI-generated blog post titles inside the admin blog area.
It combines HTMX-powered interactions with background LLM pipelines to keep
the authoring experience responsive.

## Directory Structure

```
ai_title_generator/
├── ai_title_generator_controller.go  # Primary HTTP controller
├── on_generate_titles.go             # HTMX action: generate new titles
├── on_approve_title.go               # HTMX action: approve a title
├── on_reject_title.go                # HTMX action: reject a title
```

## Request Flow

1. **Entry point** – `ai_title_generator_controller.go` exposes `Handler`, which is wired through `internal/controllers/admin/blog/routes.go` under the `ai-title-generator` controller key.
2. **Data preparation** – `prepareData` loads all records of type `blogai.POST_RECORD_TYPE` from the custom store and converts them to `blogai.RecordPost` structures.
3. **Rendering** – `view` builds the admin layout with breadcrumbs, a "Generate New Titles" HTMX button, and the table returned by `tableExistingTitles`.
4. **HTMX dispatch** – When a POST request lands, `Handler` looks at the `action` parameter and routes to the appropriate helper (generate, approve, reject). Every helper returns a SweetAlert response so the front-end can refresh itself.

## HTMX Actions

| Action (`action` param) | Handler / File              | Purpose |
|-------------------------|-----------------------------|---------|
| `generate_titles`       | `onGenerateTitles` (`on_generate_titles.go`) | Calls the LLM pipeline to create fresh title suggestions and persists them. |
| `approve_title`         | `onApproveTitle` (`on_approve_title.go`)     | Marks a pending title as approved so it becomes eligible for the post generator. |
| `reject_title`          | `onRejectTitle` (`on_reject_title.go`)       | Marks a title as rejected, preventing it from being reused. |
| `generate_post`         | Routed to the AI Post Generator controller via `Generate Post` button. |

Each HTMX endpoint returns a SweetAlert modal payload. Success alerts typically include a timed redirect back to the title generator list.

### What Approve / Reject Do

- **Approve** updates the record’s status to `approved` (it is not deleted) so the AI Post Generator can pick it up. Approved titles remain in the table and can still be regenerated into posts later.
- **Reject** keeps the record but flags its status as `rejected`. Rejected titles stay stored and are included in the `existing_titles` list, so future generations skip them automatically.

## LLM Pipeline Overview

`onGenerateTitles` orchestrates a workflow using these steps:

1. **`stepHandlerFindExistingTitles`** – Loads every stored AI title record and current blog post title, merging them into the `existing_titles` slice used for deduplication.
2. **`stepHandlerGenerateTitles`** – Initializes the LLM engine (`shared.LlmEngine`), hydrates a `blogai.TitleGeneratorAgent`, and produces candidate titles while avoiding duplicates.
3. **`stepHandlerSaveTitles`** – Persists each title as a `customstore.Record` with basic metadata (status, timestamps). Newly saved titles begin with status `pending`.

### Deduplication & Output Count

- **Source list** – Before hitting the LLM, `onGenerateTitles` merges every stored AI title plus live blog post titles into an `existing_titles` slice used for deduplication.
- **Prompt constraint** – The title agent’s system prompt instructs the LLM to avoid repeats and the post-processing filter removes any string that still exists in the source list.
- **Batch size** – Each run asks the LLM for 10 suggestions; after filtering, all remaining unique titles are saved as new pending records.

## Status Lifecycle

`getStatusBadgeClass` converts the current status into a Bootstrap badge color. The standard states are:

- `pending` → waiting for review (approve or reject)
- `approved` → ready to be converted into a full blog post
- `rejected` → discarded suggestion
- `draft` / `published` → used for downstream tracking once a post is generated

## Integrations

- **Custom Store** (`app.GetCustomStore`) – Persistent storage for title records.
- **HTMX** – Buttons and forms use `hx-post` and `hx-target="body"` to append modal responses.
- **SweetAlert2** – User feedback for success/failure and timed redirects.
- **AI Post Generator** – Approved titles expose a "Generate Post" button linking to the AI Post Generator controller with the record ID.

## Adding New Actions or Statuses

1. Define a new action constant in `ai_title_generator_controller.go`.
2. Extend the POST switch in `Handler` to route the action to a new helper.
3. Implement the helper in its own `on_<action>.go` file and ensure it returns a valid SweetAlert payload.
4. Update `getStatusBadgeClass` and any table badge copy if you introduce new statuses.

## Troubleshooting

- **No titles visible:** confirm the custom store contains `blogai.POST_RECORD_TYPE` entries. The controller silently skips malformed records but logs warnings.
- **LLM failures:** `onGenerateTitles` returns an error SweetAlert if `shared.LlmEngine` or the `blogai` helper fails. Check environment credentials for the LLM provider.
- **HTMX response missing:** ensure `cdn.Htmx_2_0_0()` is included in the layout options if HTMX interactions stop working.

With this README, future contributors can quickly understand how the AI Title Generator controller hooks into the admin stack and how to extend or debug it.
