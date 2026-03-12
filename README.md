# metabase-cli

A CLI for [Metabase](https://metabase.com) — query cards, dashboards, and collections from the command line. Single binary, no runtime required.

## Install

**Homebrew (macOS/Linux):**
```bash
brew install danielbrodie/tap/metabase
```

**Download binary:** grab the latest release from [GitHub Releases](https://github.com/danielbrodie/metabase-cli/releases), extract, and put `metabase` on your `$PATH`.

**Build from source:**
```bash
go install github.com/danielbrodie/metabase-cli@latest
```

## Auth

```bash
metabase auth login --url http://localhost:3000 --email you@example.com
metabase auth status
metabase auth token
metabase auth logout
```

Credentials are stored in `~/.config/metabase/config.json`.

## Commands

### Search
```bash
metabase search "revenue"
metabase search "revenue" --models card
metabase search "revenue" --models dashboard
metabase search "revenue" --limit 10
```

### Databases
```bash
metabase databases list
metabase databases get 1
metabase databases get 1 --include-tables
metabase databases get 1 --include-fields
metabase databases metadata 1
metabase databases schemas 1
```

### Collections
```bash
metabase collections tree
metabase collections tree --search "analytics" -L 3
metabase collections get 42
metabase collections get root
metabase collections items 42
metabase collections items 42 --models card
metabase collections create --name "My Collection" --parent-id 42
metabase collections update 42 --name "New Name"
metabase collections archive 42
```

### Cards
```bash
metabase cards list --filter mine
metabase cards list --collection-id 42
metabase cards get 123
metabase cards run 123
metabase cards run 123 --limit 100 --json
metabase cards import --file card.json
metabase cards import --file card.json --id 123    # update existing
metabase cards archive 123
metabase cards delete 123 --force
```

### Dashboards
```bash
metabase dashboards list --collection-id 42
metabase dashboards get 456
metabase dashboards get 456 --include-cards
metabase dashboards export 456
metabase dashboards import --file dashboard.json
metabase dashboards import --file dashboard.json --id 456    # update existing
metabase dashboards revisions 456
metabase dashboards revert 456 789
metabase dashboards archive 456
metabase dashboards delete 456 --force
```

### Resolve
```bash
metabase resolve 'http://localhost:3000/dashboard/5'
metabase resolve 'http://localhost:3000/question/42'
```

## Flags

| Flag | Description |
|------|-------------|
| `-p, --profile` | Named profile (default: `"default"`) |
| `--json` | JSON output |
| `-v, --verbose` | Verbose output |

## Profiles

Multiple Metabase instances are supported via named profiles:

```bash
metabase auth login --url https://metabase.prod.example.com --email me@example.com --profile prod
metabase cards list --profile prod
```

## AI Agent Integration

A `SKILL.md` is included in `skill/` for use with AI coding agents that support skill-based tool discovery. Drop it into your agent's skills directory to enable natural language access to Metabase commands.

## License

MIT
