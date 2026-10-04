# Models

Example neat ORM models (see `note.go`).

Blueprint's convention is the dracory stores (`entitystore`, `userstore`,
...) for anything durable — so nothing in this folder is wired up.

Use this folder while you are still exploring a domain and don't yet know
the final shape of a reusable store package: prototype with a plain model
+ neat's query builder (see `dracory/pico`), then convert to a store
package under `pkg/` once the fields and queries settle.
