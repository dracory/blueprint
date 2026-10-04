# internal/config

Blueprint-style configuration: env-var loading into
`configImplementation`, exposed through the `config.ConfigInterface`
aggregate (see `config_interfaces.go`).

## Rules

- **Every config value has getter + setter on an interface.** Add new
  settings in a domain file (`*_config.go`), expose them via the matching
  `*ConfigInterface` in `config_interfaces.go`, and implement on
  `configImplementation`.
- **Store toggles use the `Get<Name>StoreUsed()` flags** in
  `stores_config.go`. Migrations, seeders, and background processes all
  gate on these — a flag set wrong silently disables the store. When
  adding a store, also update `store_builders.go`, the app interface,
  migrations, and `testutils.WithXStore`.
- `config_implementation.go` starts with `z_`-style comment — it must
  sort after user-configurable files; keep that convention when adding
  fields.
- Secrets/keys come from env only — never default a real credential in
  `DefaultConf` or constants.
