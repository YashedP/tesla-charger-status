# CONTINUITY

[PLANS]
- 2026-02-22T22:07Z [USER] Implement Go Tesla wrapper with OAuth and one `GET /v1/is-charging` endpoint returning `true`/`false`; remove `/healthz`; keep SQLite/key paths fixed and add Python key scripts.

[DECISIONS]
- 2026-02-22T22:07Z [USER] Router is `chi`.
- 2026-02-22T22:07Z [USER] Fixed DB path `./data/tesla.sqlite` and fixed key path `./secrets/token_enc_key.b64` (not env vars).
- 2026-02-22T22:07Z [USER] Failure semantics for live fetch errors map to `false`.
- 2026-02-22T22:07Z [CODE] Required env vars retained: `TESLA_CLIENT_ID`, `TESLA_CLIENT_SECRET`, `TESLA_REDIRECT_URI`, `TESLA_VIN`, `SHORTCUT_BEARER_TOKEN`, `TESLA_BASE_URL`.
- 2026-02-23T04:32Z [USER] Replaced `TESLA_REDIRECT_URI` env var with `APP_BASE_URL`, deriving callback URI as `<APP_BASE_URL>/oauth/callback` (supersedes older redirect-env decision).
- 2026-02-23T06:10Z [USER] Obsidian writing skill should be named `obsidian-page`, stored canonically in `~/.codex/skills`, with strict page-quality standards.
- 2026-02-23T07:18Z [USER] `ai-config` should sync Codex user config + non-system skills and Claude `CLAUDE.md`, `settings.json`, plugin manifests, and `plugins/marketplaces/personal`, with idempotent symlink install behavior.

[PROGRESS]
- 2026-02-22T22:07Z [CODE] Created service structure for config loading, OAuth flow, encrypted token store, Tesla Fleet API client, and chi router.
- 2026-02-23T03:10Z [CODE] Added project files for server, internal packages, tests, Docker workflow, `.env.example`, and Python key generation/validation scripts.
- 2026-02-23T03:10Z [CODE] Pinned Go module dependencies in `go.mod` using current versions from Go module proxy metadata.
- 2026-02-23T03:15Z [CODE] Rewrote `README.md` as a full personal-use onboarding guide covering prerequisites, env vars, OAuth bootstrap, iPhone Shortcuts wiring, Docker usage, and troubleshooting.
- 2026-02-23T03:16Z [CODE] Added `Makefile` targets for key script workflow (`key-generate`, `key-validate`, `key-scripts`) and backend startup (`run`).
- 2026-02-23T03:19Z [CODE] Updated `.gitignore` to ignore Python cache artifacts (`__pycache__/`, `*.py[cod]`) per user request.
- 2026-02-23T03:26Z [CODE] Refactored `internal/store/sqlite.go` to use `go:embed` with SQL files in `internal/store/sql/` (`migrate.sql`, `upsert_token.sql`, `select_token.sql`) instead of inline query strings.
- 2026-02-23T03:50Z [TOOL] Ran `go mod tidy` successfully; module dependency graph resolved and `go.sum` populated.
- 2026-02-23T04:02Z [CODE] Added automatic `.env` loading for native runs in `cmd/server/main.go` via `github.com/joho/godotenv` (non-overriding load).
- 2026-02-23T04:02Z [TOOL] Ran `gofmt`, `go mod tidy`, and `go test ./...` successfully after adding `godotenv`.
- 2026-02-23T04:10Z [USER] Replaced `0o700` magic number in `cmd/server/main.go` with named constant `privateDirPerm` and added explanatory comments in `ensureParentDirs`.
- 2026-02-23T04:17Z [USER] Hardened key loading by enforcing owner-only key-file permissions in `internal/crypto/aesgcm.go` before reading/decode.
- 2026-02-23T04:17Z [CODE] Added regression test `TestLoadKeyFromFileRejectsOpenPermissions` in `internal/crypto/aesgcm_test.go`.
- 2026-02-23T04:21Z [USER] Removed magic driver string in `internal/store/sqlite.go` by introducing `sqliteDriverName` constant and using it in `sql.Open`.
- 2026-02-23T04:32Z [CODE] Updated config loader, `.env.example`, and `README.md` to use `APP_BASE_URL`; redirect URI is now computed in code using `/oauth/callback`.
- 2026-02-23T04:52Z [USER] Replaced low-level Fleet API call boilerplate in `internal/tesla/client.go` with Resty-based request flow while preserving `GetChargingState` interface/behavior.
- 2026-02-23T04:52Z [CODE] Added Tesla client defaults (`Accept: application/json`, stable `User-Agent`) and retry policy (3 retries, 100-300ms backoff, retry on transport/5xx/429).
- 2026-02-23T04:52Z [CODE] Expanded `internal/tesla/client_test.go` to cover header defaults, transient retry success, HTTP error handling, and API envelope errors.
- 2026-02-23T05:33Z [USER] Removed Tesla error-body clipping in `internal/tesla/client.go` by deleting `maxErrorBodyBytes`/`clipBody` and returning full upstream body text in error messages.
- 2026-02-23T06:10Z [CODE] Created skill files at `~/.codex/skills/obsidian-page`: `SKILL.md`, `agents/openai.yaml`, and references for Obsidian page quality/templates.
- 2026-02-23T07:18Z [CODE] Created `~/ai-config` structure with Codex/Claude managed files, copied non-system Codex skills, copied Claude personal marketplace plugins, and added `install/install.sh` + `install/verify.sh`.
- 2026-02-23T07:18Z [CODE] Fixed symlink resolution bug in installer/verifier for absolute links and confirmed repeatable no-op behavior on subsequent installs.
- 2026-02-23T07:22Z [TOOL] Committed `~/ai-config` root commit (`baa154b`) and pushed `main` to `origin` (`https://github.com/YashedP/ai-config.git`).
- 2026-02-23T07:29Z [USER] Requested dedicated `install/sync.sh` to import new local Codex/Claude additions into `~/ai-config`.
- 2026-02-23T07:29Z [CODE] Added `~/ai-config/install/sync.sh` with `--dry-run`, safe canonical-path checks, Codex skill import (excluding `.system`), and Claude user/plugin import.
- 2026-02-23T07:29Z [CODE] Updated `~/ai-config/README.md` to document sync commands and recommended `sync -> install -> verify -> commit/push` workflow.
- 2026-02-23T07:30Z [USER] Requested a `Makefile` in `~/ai-config` for script shortcuts.
- 2026-02-23T07:30Z [CODE] Added `~/ai-config/Makefile` targets: `sync`, `sync-dry-run`, `install`, `install-dry-run`, `verify`, and `apply`.
- 2026-02-23T07:30Z [CODE] Extended `~/ai-config/README.md` with Make-based usage examples.
- 2026-02-23T07:31Z [CODE] Refined `~/ai-config/README.md` update workflow to prioritize `make apply` and added script-only equivalent flow.
- 2026-02-23T07:31Z [TOOL] Committed `~/ai-config` changes as `62534a0` and pushed `main` to `origin`.

[DISCOVERIES]
- 2026-02-22T22:07Z [TOOL] Host environment lacks `go` binary and Docker runtime (`docker` command unavailable in WSL distro). Validation must be run once toolchain is available.
- 2026-02-23T03:10Z [TOOL] Python runtime is available; key scripts compile and pass functional test using `/tmp/token_enc_key_test.b64`.
- 2026-02-23T03:16Z [TOOL] Make dry-run confirms target command wiring: `make -n key-scripts` and `make -n run`.
- 2026-02-23T03:50Z [TOOL] Go toolchain is now available in shell context (confirmed by successful `go mod tidy` execution).
- 2026-02-23T04:02Z [CODE] `godotenv.Load()` preserves existing process environment values while filling unset values from `.env`, which keeps Docker/env-injected settings authoritative.
- 2026-02-23T04:17Z [CODE] Permission checks skip Windows due to unreliable POSIX mode bits; Unix-like systems reject key files with any group/other permission bits set (`mode & 0o077 != 0`).
- 2026-02-23T04:52Z [TOOL] `go get github.com/go-resty/resty/v2@v2.17.2` succeeded and introduced `golang.org/x/net` transitively.
- 2026-02-23T07:18Z [TOOL] Initial verify failures revealed absolute-symlink path resolution issue (`dirname + absolute readlink`), causing false drift; resolved by explicit absolute/relative link-target handling.
- 2026-02-23T07:29Z [TOOL] `install/sync.sh --dry-run` reports `OK` for already-linked managed paths and skips `.system`, confirming idempotent no-op behavior on current machine.
- 2026-02-23T07:30Z [TOOL] `make sync-dry-run` and `make install-dry-run` execute successfully and route to the expected shell scripts.
- 2026-02-23T07:31Z [TOOL] `make help` in `~/ai-config` confirms target wiring after README and Makefile updates.

[OUTCOMES]
- 2026-02-22T22:07Z [ASSUMPTION] UNCONFIRMED until tooling is available: compile/test status.
- 2026-02-23T03:10Z [CODE] Implementation completed per agreed plan with fixed SQLite/key paths and no `/healthz`; Go build/test remains UNCONFIRMED in this environment due missing Go/Docker toolchain.
- 2026-02-23T04:17Z [TOOL] `go test ./internal/crypto -v` and `go test ./...` pass after permission-hardening changes.
- 2026-02-23T04:21Z [TOOL] `go test ./...` passes after sqlite driver string constant refactor.
- 2026-02-23T04:32Z [TOOL] `go test ./...` passes after env-var rename and redirect URI derivation changes.
- 2026-02-23T04:52Z [TOOL] `go test ./...` passes after Resty migration for Tesla API calls.
- 2026-02-23T05:33Z [TOOL] `go test ./...` passes after removing Tesla error-body clipping logic.
- 2026-02-23T06:10Z [TOOL] Skill directory and metadata validated via file listing and content checks.
- 2026-02-23T07:18Z [TOOL] `~/ai-config/install/install.sh` now idempotent (`OK` on rerun) and `~/ai-config/install/verify.sh` passes for all managed Codex and Claude links.
- 2026-02-23T07:22Z [TOOL] `~/ai-config` is now committed and remote-tracked; branch `main` tracks `origin/main`.
- 2026-02-23T07:29Z [TOOL] `bash -n` passes for `install/install.sh`, `install/verify.sh`, and new `install/sync.sh`; pending local git changes in `~/ai-config`: `README.md` modified, `install/sync.sh` added.
- 2026-02-23T07:30Z [TOOL] `make help` works and documents all new targets; pending local git changes in `~/ai-config`: `README.md` modified, `Makefile` added, `install/sync.sh` added.
- 2026-02-23T07:31Z [TOOL] `~/ai-config` is clean after commit `62534a0`, and remote `main` is up to date.
