# Implementation Plan: open-sandbox MVP

**Branch**: `mvp/plan` | **Date**: 2026-01-22 | **Spec**: `.specify/memory/spec.md`
**Input**: Feature specification from `.specify/memory/spec.md`

## Summary

Deliver a single-node/single-container sandbox MVP with a unified HTTP entry and Browser/VNC/IDE/Jupyter/Shell/File/Code capabilities. Focus on demo readiness and a verifiable end-to-end flow, using Go + net/http with minimal dependencies and a unified error/response model.

## Technical Context

**Language/Version**: Go 1.24+ (tested: 1.24.11 on Windows)
**Primary Dependencies**: Standard library net/http only; no routing libraries
**Storage**: Absolute host workspace at `SANDBOX_WORKSPACE` (default `<repo_root>/workspace`)
**Testing**: Go testing (standard library)
**Target Platform**: Windows host + local Docker
**Project Type**: Single service
**Performance Goals**: No hard targets (avoid obvious blocking)
**Constraints**: All cache/log/build artifacts must stay under `SANDBOX_ROOT` (default `<repo_root>`)
**Container Mapping**: If containerized, mount `SANDBOX_WORKSPACE` to `/workspace` and use `/workspace` paths inside the container
**Scale/Scope**: Single-node/single-container MVP

## Constitution Check

- MVP First, Demo-Ready: PASS
- Single-Node, Single-Container: PASS
- Simplicity Over Cleverness: PASS
- Test-First (Non-Negotiable): PASS
- Safe-by-Default for MVP: PASS
- Commenting Standard (English only): PASS
- Default Branch is `main`: PASS

## Project Structure

### Documentation

```text
specs/open-sandbox-mvp/
  plan.md
  research.md
  data-model.md
  quickstart.md
  contracts/
  tasks.md
```

### Source Code (repository root)

```text
cmd/
  server/
    main.go

internal/
  api/
    router.go
    middleware.go
    handlers/
  browser/
  vnc/
  shell/
  file/
  codeexec/
  jupyter/
  codeserver/
  config/
  platform/

pkg/
  types/

tests/
  integration/
  unit/
```

**Structure Decision**: Single Go service. `cmd/server` contains bootstrap only; features live under `internal/*`. Unified error model and response schema in `internal/api` and `pkg/types`.

## Milestones & Steps (TDD-first)

### 1) Foundation: types, errors, router, workspace

- Define unified response and error types in `pkg/types`.
- Implement error mapping and trace IDs in `internal/api`.
- Build minimal router with net/http and handler registration.
- Enforce absolute path usage and create workspace if missing.
- Tests: response schema serialization, error mapping, workspace creation.

### 2) Unified entry API

- Implement `GET /v1/sandbox` capability discovery response.
- Include health status and service URLs.
- Tests: contract for response shape and errors.

### 3) Shell API (non-interactive)

- Implement `/v1/shell/exec` with stdout/stderr/exit code.
- Ensure working directory is the workspace and paths are absolute.
- Tests: `echo test`, `dir`/`ls`, non-zero exit handling.

### 4) File API

- Implement read/write/list/search/replace endpoints.
- Add size limits and error handling for missing/permission issues.
- Tests: CRUD flow in workspace, search/replace correctness.

### 5) Code execution API (Python/Node)

- Provide minimal execution wrapper for Python and Node commands.
- Restrict execution to workspace; capture stdout/stderr/exit code.
- Tests: process file size and write output to `output.txt`.

### 6) Browser API + VNC

- Define interfaces for headed browser control and VNC access.
- Implement endpoints for CDP address discovery, navigate, screenshot.
- VNC endpoint serves static UI and connects to the desktop session.
- Tests: CDP address presence, screenshot written to workspace.

### 7) Jupyter + Code Server access

- Provide reverse proxy or path mapping for `/jupyter` and `/code-server/`.
- Ensure links appear in unified entry response.
- Tests: endpoint returns HTTP 200 with HTML content.

### 8) Docs & demo validation

- Document ports, env vars, startup, limitations/known issues.
- Provide quickstart for Windows + Docker.
- Execute end-to-end flow and capture outputs under workspace.
- Tests: integration test for the full scenario if feasible.

## End-to-End Validation Flow

1) Start service locally on Windows.
2) `GET /v1/sandbox` returns capabilities.
3) Browser API opens `https://example.com` and writes screenshot to `<SANDBOX_WORKSPACE>/screenshots/example.png`.
4) File API reads screenshot metadata or saves page text.
5) Code exec processes file and writes `<SANDBOX_WORKSPACE>/output.txt`.
6) File API reads `output.txt` and returns contents.
7) VNC page loads at `http://localhost:8080/vnc/index.html`.
8) Jupyter and code-server endpoints return valid pages.

## Risks & Open Questions

- [NEEDS CLARIFICATION] Code execution isolation strategy (container/OS-level/none)?
- [NEEDS CLARIFICATION] Is proxy/port forwarding in scope for MVP?
- [NEEDS CLARIFICATION] Is MCP Hub required, and if yes, what interface?

## Change Management

- Atomic development and commits.
- Every change includes tests written first.
- All runtime artifacts stay under `SANDBOX_ROOT`.

## Complexity Tracking

N/A
