# Implementation Plan: Codex MCP Alignment

**Branch**: `001-mcp-codex-alignment` | **Date**: 2026-01-23 | **Spec**: `specs/001-mcp-codex-alignment/spec.md`
**Input**: Feature specification from `/specs/001-mcp-codex-alignment/spec.md`

## Summary

Align the sandbox MCP server with Codex expectations by adding standard MCP discovery/invocation semantics, exposing a stdio entrypoint alongside HTTP/SSE, documenting Codex configuration, and providing a repeatable smoke test.

## Technical Context

**Language/Version**: Go 1.24  
**Primary Dependencies**: Standard library; `github.com/golang-jwt/jwt/v5` (MCP auth)  
**Storage**: Filesystem only (workspace under `SANDBOX_ROOT`)  
**Testing**: `go test ./...` (unit + integration)  
**Target Platform**: Windows + WSL + Docker (single-node)  
**Project Type**: Single service (monorepo)  
**Performance Goals**: Interactive local usage; simple MCP calls should respond in under ~2 seconds on a single machine  
**Constraints**: Strict TDD, unified error model, minimal dependencies, no breaking changes to existing HTTP APIs, runtime artifacts under `SANDBOX_ROOT`  
**Scale/Scope**: Single-machine demo and smoke-test usage (single-user)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- MVP First, Demo-Ready: PASS
- Single-Node, Single-Container: PASS
- Simplicity Over Cleverness: PASS
- Test-First (Non-Negotiable): PASS
- Safe-by-Default for MVP: PASS
- Commenting Standard (English only): PASS
- Unified Error Model: PASS
- Minimal Dependencies: PASS

## Project Structure

### Documentation (this feature)

```text
specs/001-mcp-codex-alignment/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
├── server/                 # Existing HTTP server
└── mcp/                    # NEW: stdio MCP entrypoint

internal/
├── api/handlers/mcp.go     # HTTP/SSE MCP routes
├── mcp/                    # MCP core (types, server, registry, auth)
└── mcp/tools/              # MCP tool bindings

tests/
├── unit/                   # MCP unit tests
└── integration/            # MCP HTTP/SSE + tool smoke tests
```

**Structure Decision**: Single service with a dedicated stdio CLI entrypoint (`cmd/mcp`) and existing HTTP server routes under `internal/api/handlers`.

## Phase 0: Outline & Research

### Research Questions

- Codex MCP client expectations for discovery/invocation methods and schema metadata.
- Best-practice transport posture for stdio vs HTTP/SSE in local tooling.
- Minimum schema representation sufficient for Codex tool invocation.
- Minimal smoke test that proves Codex-aligned MCP behavior without external tooling.

### Output

- `specs/001-mcp-codex-alignment/research.md` with resolved decisions.

## Phase 1: Design & Contracts

### Data Model

- Extract MCP entities (server, tool, schema, transport, auth) into `data-model.md`.

### Contracts

- Define JSON-RPC request/response shapes and HTTP/SSE endpoints in `/contracts/`.

### Quickstart

- Provide Codex configuration steps for stdio and HTTP/SSE connections.

### Agent Context Update

- Run `.specify/scripts/powershell/update-agent-context.ps1 -AgentType codex`.

## Post-Design Constitution Check

- MVP First, Demo-Ready: PASS
- Single-Node, Single-Container: PASS
- Simplicity Over Cleverness: PASS
- Test-First (Non-Negotiable): PASS
- Safe-by-Default for MVP: PASS
- Commenting Standard (English only): PASS
- Unified Error Model: PASS
- Minimal Dependencies: PASS

## Complexity Tracking

> No constitution violations requiring justification.
