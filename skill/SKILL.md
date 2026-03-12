---
name: metabasecli
disable-model-invocation: true
description: CLI for querying Metabase (cards/questions, dashboards, collections, databases, search). Use when working with Metabase to search entities, browse collections, run saved questions, export/import dashboards, explore database metadata, or resolve Metabase URLs. Triggered by requests involving Metabase data, dashboard queries, saved questions, collection browsing, or Metabase automation.
trigger-keywords: metabase, metabase card, metabase question, metabase query, metabase dashboard, metabase collection, metabase database
allowed-tools: Bash(metabase --help), Bash(metabase auth status:*), Bash(metabase auth login:*), Bash(metabase auth token:*), Bash(metabase auth logout:*), Bash(metabase search:*), Bash(metabase resolve:*), Bash(metabase databases list:*), Bash(metabase databases get:*), Bash(metabase databases metadata:*), Bash(metabase databases schemas:*), Bash(metabase collections tree:*), Bash(metabase collections get:*), Bash(metabase collections items:*), Bash(metabase collections create:*), Bash(metabase collections update:*), Bash(metabase collections archive:*), Bash(metabase cards list:*), Bash(metabase cards get:*), Bash(metabase cards run:*), Bash(metabase cards import:*), Bash(metabase cards archive:*), Bash(metabase cards delete:*), Bash(metabase dashboards list:*), Bash(metabase dashboards get:*), Bash(metabase dashboards export:*), Bash(metabase dashboards import:*), Bash(metabase dashboards revisions:*), Bash(metabase dashboards revert:*), Bash(metabase dashboards archive:*), Bash(metabase dashboards delete:*)
---

# metabasecli

Command-line interface for Metabase API operations.

## Terminology

- Cards = Questions/Queries in the Metabase UI
- Dashboards = Dashboards
- Collections = Collections (folders for organizing cards and dashboards)
- Databases = Connected data sources

## Flag Placement

Always place flags after the full command path, not between `metabase` and the command group.

Correct:
- `metabase databases list --profile prod --json`
- `metabase cards run 123 --limit 100`

Wrong:
- `metabase -p prod databases list`

## Global Flags

- `-p`, `--profile` — Named profile to use (default: "default")
- `-v`, `--verbose` — Verbose output
- `--json` — JSON output

## JSON Output Format

- Success: `{"success": true, "data": ...}`
- Error: `{"success": false, "error": {"code": "API_ERROR", "message": "..."}}`

## Authentication

- `metabase auth status`
- `metabase auth login --url http://localhost:3000 --email you@example.com`
- `metabase auth token`
- `metabase auth logout`

## Search

- `metabase search "revenue"`
- `metabase search "revenue" --models card`
- `metabase search "revenue" --models dashboard`
- `metabase search "revenue" --collection-id 42`
- `metabase search "revenue" --limit 10`

## Resolve URLs

- `metabase resolve 'https://metabase.example.com/question/123'`
- `metabase resolve 'https://metabase.example.com/dashboard/456'`
- `metabase resolve '/collection/789'`
- `metabase resolve '/browse/databases/1'`

## Databases

- `metabase databases list`
- `metabase databases get 1`
- `metabase databases get 1 --include-tables`
- `metabase databases get 1 --include-fields`
- `metabase databases metadata 1`
- `metabase databases schemas 1`

## Collections

- `metabase collections tree`
- `metabase collections tree --search "analytics" -L 3`
- `metabase collections get 42`
- `metabase collections get root`
- `metabase collections items 42`
- `metabase collections items 42 --models card`
- `metabase collections create --name "My Collection" --parent-id 42`
- `metabase collections update 42 --name "New Name"`
- `metabase collections archive 42`

## Cards

- `metabase cards list --filter mine`
- `metabase cards list --collection-id 42`
- `metabase cards get 123`
- `metabase cards run 123`
- `metabase cards run 123 --limit 100 --json`
- `metabase cards import --file card.json`
- `metabase cards import --file card.json --id 123`
- `metabase cards archive 123`
- `metabase cards delete 123 --force`

## Dashboards

- `metabase dashboards list --collection-id 42`
- `metabase dashboards get 456`
- `metabase dashboards get 456 --include-cards`
- `metabase dashboards export 456`
- `metabase dashboards import --file dashboard.json`
- `metabase dashboards import --file dashboard.json --id 456`
- `metabase dashboards revisions 456`
- `metabase dashboards revert 456 789`
- `metabase dashboards archive 456`
- `metabase dashboards delete 456 --force`
