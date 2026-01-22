open-sandbox
=============

One-node, one-container sandbox for AI/human co-development: browser + VNC + IDE + Jupyter + shell + file + code execution.

Project Goals
-------------
- Deliver a demo-ready MVP that is actually usable on a single machine.
- Provide a unified HTTP API for browser, shell, file, and code execution workflows.
- Ensure all runtime artifacts (cache/logs/build outputs) stay on D:\.
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
- Auth can be off by default, but has a JWT toggle placeholder.
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
docker run --rm -p 8080:8080 -v D:\Desktop\sandbox\open-sandbox\workspace:/workspace open-sandbox
```

Ports
-----
- API + static pages: 8080 (default)

Environment Variables
---------------------
- `SANDBOX_ADDR` (default `:8080`)
- `SANDBOX_BROWSER_CDP` (CDP websocket address for a headed browser)

Runtime Artifacts
-----------------
All cache, logs, and build outputs must live under `D:\Desktop\sandbox\open-sandbox`:
- Cache: `D:\Desktop\sandbox\open-sandbox\.cache`
- Logs: `D:\Desktop\sandbox\open-sandbox\logs`
- Build outputs: `D:\Desktop\sandbox\open-sandbox\build`

Limitations / TODO
------------------
- Browser actions are stubbed; CDP integration and real navigation/screenshot capture need implementation.
- VNC, Jupyter Lab, and code-server endpoints are placeholders pending external service wiring.
- JWT auth toggle is a placeholder and not enforced yet.
- Acceptance criteria are not met for: CDP connectable browser control, VNC takeover, and real Jupyter/code-server access.

Docs & Specs
------------
- Constitution: `.specify/memory/constitution.md`
- Spec, plan, tasks: generated via `/speckit.*`

License
-------
Apache-2.0
