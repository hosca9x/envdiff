# envdiff

> Diff and reconcile `.env` files across environments with secret masking.

---

## Installation

```bash
go install github.com/yourname/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envdiff.git && cd envdiff && go build -o envdiff .
```

---

## Usage

Compare two `.env` files and highlight differences, with sensitive values automatically masked:

```bash
envdiff .env.development .env.production
```

**Example output:**

```
~ DB_HOST        dev-db.local        → prod-db.internal
+ REDIS_URL      [masked]
- DEBUG          true
= APP_NAME       myapp
```

### Flags

| Flag | Description |
|------|-------------|
| `--no-mask` | Disable secret masking |
| `--only-missing` | Show only keys missing from the target file |
| `--export` | Output a reconciled `.env` file to stdout |

### Reconcile environments

```bash
envdiff --export .env.staging .env.production > .env.reconciled
```

---

## How It Works

`envdiff` parses both files, detects added, removed, and changed keys, and automatically masks values for keys matching common secret patterns (e.g. `*_KEY`, `*_SECRET`, `*_TOKEN`, `*_PASSWORD`).

---

## License

MIT © [yourname](https://github.com/yourname)