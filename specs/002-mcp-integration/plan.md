# Implementation Plan: MCP Integration

**Branch**: `002-mcp-integration` | **Date**: 2026-01-22 | **Spec**: `.specify/memory/spec.md`
**Input**: MCP integration feature specification

## Summary

Build a minimal MCP server (JSON-RPC 2.0 over stdio) that maps MCP tools to existing sandbox services with unified errors and workspace safety.

## Technical Context

**Language/Version**: Go 1.24+
**Primary Dependencies**: Standard library
**JWT Verification**: Allow `github.com/golang-jwt/jwt/v5` for token verification only; keep dependency scope minimal and document usage.
**Transport**: stdio JSON-RPC (primary) + HTTP + SSE
**Storage**: `<SANDBOX_WORKSPACE>` and `SANDBOX_ROOT`
**Testing**: Go testing (standard library)
**Target Platform**: Windows + WSL

## Constitution Check

- MVP First, Demo-Ready: PASS
- Single-Node, Single-Container: PASS
- Simplicity Over Cleverness: PASS
- Test-First (Non-Negotiable): PASS
- Safe-by-Default for MVP: PASS
- Commenting Standard (English only): PASS

## Project Structure

```text
internal/
  mcp/
    types.go
    server.go
    handlers.go
    registry.go
    auth.go
    tools/
      browser.go
      file.go
      shell.go
      code.go

tests/
  unit/
    mcp_types_test.go
    mcp_registry_test.go
    mcp_capabilities_test.go
    mcp_auth_test.go
  integration/
    mcp_tools_test.go
    mcp_http_test.go
    mcp_sse_test.go
```

## Milestones & Steps (TDD-first)

### 1) MCP Core
- Tests: JSON-RPC parsing, protocol version rejection, unified error schema.
- Define JSON-RPC request/response types.
- Implement stdio server loop.
- Implement `mcp.capabilities`.

### 2) Tool Registry
- Tests: registry lookup + permission metadata schema.
- Add registry with permissions/metadata.

### 3) Tool Mappings
- Tests: tool calls map to service calls, workspace boundary.
- Implement browser tools.
- Implement file tools.
- Implement shell and code tools.

### 4) HTTP Transport
- Tests: HTTP round-trip, `/mcp` path, auth enabled/disabled.
- Implement HTTP endpoint for MCP requests/responses.

### 5) SSE Transport
- Tests: SSE round-trip, `/mcp/sse` path, auth enabled/disabled.
- Implement SSE endpoint for MCP requests/responses.

### 6) Auth Controls
- Tests: `MCP_AUTH_ENABLED` off ignores Authorization, on enforces JWT requirements.
- Implement auth toggles and token validation (HTTP/SSE only).

### 7) Docs
- Document MCP usage and examples.

### 8) NFR Verification (Lightweight)
- Add README checklist covering minimal deps, no breaking HTTP APIs, demo-ready single machine, TDD, and SANDBOX_ROOT.

## Risks & Open Questions

- Transport selection if non-stdio clients are required.
- Security model for tool access.

## Change Management

- Atomic commits.
- Tests written before implementation.

## Complexity Tracking

N/A
