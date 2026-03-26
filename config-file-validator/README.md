# Config File Validator

A CLI tool that reads and validates YAML/JSON configuration files before deployment — catching syntax errors and invalid values early.

## Features

- Parses and validates both **JSON** and **YAML** (`.yml`) files
- Checks for valid syntax (malformed JSON/YAML)
- Validates required fields are present and not empty
- Validates value ranges (e.g. port numbers between 1–65535)
- Reports **all** errors at once, not just the first one
- Clean output with ✓/✗ indicators and proper exit codes

## Usage

```bash
go run config-file-validator.go <filepath>
```

### Examples

```bash
# Validate a JSON config
go run config-file-validator.go sample.json

# Validate a YAML config
go run config-file-validator.go sample.yaml
```

### Sample Output

**Valid config:**
```
Validating sample.json...

  ✓ Valid JSON syntax
  ✓ All required fields present
  ✓ All values within valid ranges

✓ Configuration is valid!
```

**Invalid values:**
```
Validating invalid-values.yaml...

  ✓ Valid YAML syntax

✗ validation failed:
  ✗ 'name' is required but missing or empty
  ✗ 'port' must be between 1 and 65535, got 99999
  ✗ 'database.host' is required but missing or empty
```

## Expected Config Structure

```json
{
    "name": "my-service",
    "image": "example/my-service:v1.0.0",
    "port": 8080,
    "env": ["APP_ENV=dev"],
    "database": {
        "host": "db.internal",
        "port": 5432,
        "name": "my_db"
    }
}
```

## Validation Rules

| Field | Rule |
|---|---|
| `name` | Required, non-empty |
| `image` | Required, non-empty |
| `port` | Required, between 1 and 65535 |
| `database.host` | Required, non-empty |
| `database.port` | Required, between 1 and 65535 |
| `database.name` | Required, non-empty |

## Test Files

| File | Purpose |
|---|---|
| `sample.json` | Valid JSON config |
| `sample.yaml` | Valid YAML config |
| `broken.json` | Malformed JSON (missing comma) |
| `invalid-values.yaml` | Valid syntax, invalid values |

## Built With

- [Go](https://go.dev/)
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)
