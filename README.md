open-sandbox
=============

One-node, one-container sandbox for AI/human co-development: browser + VNC + IDE + Jupyter + shell + file + code execution.

Project Goals
-------------
- Deliver a demo-ready MVP that is actually usable on a single machine.
- Provide a unified HTTP API for browser, shell, file, and code execution workflows.
- Ensure all runtime artifacts (cache/logs/build outputs) stay under `SANDBOX_ROOT`.
- All code comments must be English-only and follow best practices (intent/why, concise, no obvious restatements).

MVP Scope (Must-Have)
---------------------
1) Unified HTTP entry API
2) Headed browser with CDP (address, screenshot, actions)
3) VNC takeover for visual control
4) Shell API (non-interactive at minimum)
5) File API (read/write/list/search/replace)
6) Code execution (Python/Node minimal viable)
7) Jupyter Lab & code-server accessible

Non-Functional Requirements
---------------------------
- Runs on Windows and local Docker.
- Docs include ports, env vars, and startup instructions.
- MCP auth can be off by default, but supports JWT verification for HTTP/SSE transport.
- No strict perf targets, but avoid obvious blocking.
- Atomic development & commits.
- TDD required: tests first, then implementation.

Quick Start
-----------
Local (Windows)
```
go test ./...
go run ./cmd/server
```

Docker (placeholder)
```
docker build -t open-sandbox .
docker run --rm -p 8080:8080 -v <SANDBOX_WORKSPACE>:/workspace open-sandbox
```

Ports
-----
- API + static pages: 8080 (default)

MCP Integration
---------------
- HTTP JSON-RPC: `POST /mcp`
- SSE JSON-RPC: `GET /mcp/sse?request=<urlencoded JSON>`

Example JSON-RPC payload:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "mcp.capabilities",
  "params": {
    "protocol_version": "1.0"
  }
}
```

Codex MCP Setup (HTTP)
----------------------
1) Start the server: `go run ./cmd/server`
2) Open `~/.codex/config.toml` and add:
```
[mcp_servers.open_sandbox]
url = "http://localhost:8080/mcp"
tool_timeout_sec = 60
```
3) If MCP auth is enabled, export a bearer token and add:
```
[mcp_servers.open_sandbox]
bearer_token_env_var = "MCP_AUTH_TOKEN"
```
4) In Codex, open the MCP panel (`/mcp`) and verify `open_sandbox` is connected.


Environment Variables
---------------------
- `SANDBOX_ADDR` (default `:8080`)
- `SANDBOX_ROOT` (base directory for runtime artifacts; defaults to repo root when available)
- `SANDBOX_WORKSPACE` (absolute workspace path; defaults to `<SANDBOX_ROOT>/workspace`)
- `SANDBOX_CACHE_ROOT` (defaults to `<SANDBOX_ROOT>/.cache`)
- `SANDBOX_LOGS_ROOT` (defaults to `<SANDBOX_ROOT>/logs`)
- `SANDBOX_BUILD_ROOT` (defaults to `<SANDBOX_ROOT>/build`)
- `SANDBOX_BROWSER_BIN` (path to Chrome/Chromium binary)
- `SANDBOX_BROWSER_CDP` (existing CDP websocket address; skips launching a new browser)
- `SANDBOX_CDP_HOST` (default `127.0.0.1`)
- `SANDBOX_CDP_PORT` (default `9222`)
- `SANDBOX_BROWSER_HEADLESS` (default `false`)
- `SANDBOX_BROWSER_NAV_TIMEOUT_SEC` (default `15`, navigation timeout)
- `SANDBOX_BROWSER_SCREENSHOT_TIMEOUT_SEC` (default `15`, screenshot timeout)
- `SANDBOX_JUPYTER_URL` (reverse proxy target, e.g. `http://localhost:8888`)
- `SANDBOX_CODESERVER_URL` (reverse proxy target, e.g. `http://localhost:8081`)
- `MCP_AUTH_ENABLED` (default `false`)
- `MCP_AUTH_JWT_SECRET` (HMAC secret for HS256/384/512)
- `MCP_AUTH_JWT_PUBLIC_KEY` (PEM public key for RS/ES/EdDSA)
- `MCP_AUTH_AUDIENCE` (optional audience validation)
- `MCP_AUTH_ISSUER` (optional issuer validation)

Runtime Artifacts
-----------------
All cache, logs, and build outputs must live under `SANDBOX_ROOT` (default `<repo_root>`). Example Windows defaults:
- Cache: `D:\Desktop\sandbox\open-sandbox\.cache`
- Logs: `D:\Desktop\sandbox\open-sandbox\logs`
- Build outputs: `D:\Desktop\sandbox\open-sandbox\build`

Limitations / TODO
------------------
- Browser requires a locally installed Chrome/Chromium or an existing CDP endpoint.
- VNC view is a live browser screen with click support, not a full desktop capture.
- Jupyter Lab and code-server are proxied endpoints; the upstream services must be running.
- MCP HTTP/SSE endpoints support JWT auth when enabled; other MVP endpoints are unauthenticated.

NFR Verification Checklist
--------------------------
- [ ] Minimal dependencies (stdlib preferred; JWT library allowed for MCP auth)
- [ ] No breaking changes to existing HTTP APIs
- [ ] Demo-ready on a single machine
- [ ] Strict TDD followed for this feature
- [ ] Runtime artifacts stay under `SANDBOX_ROOT`

Docs & Specs
------------
- Constitution: `.specify/memory/constitution.md`
- Spec, plan, tasks: generated via `/speckit.*`

License
-------
Apache-2.0
