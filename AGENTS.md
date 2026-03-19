# Repo AGENTS Notes

## Container-first workflow

Use containerized tooling by default:

- Build: `docker build -t tesla-charger-status .`
- Run service: `docker compose up --build`
- Generate key: `docker compose run --rm tesla-charger-status python3 /app/scripts/gen_token_key.py`

If container tooling is unavailable on host, install/configure Docker Desktop WSL integration before running project checks.
