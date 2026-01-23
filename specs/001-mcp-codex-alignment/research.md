# Research Log: Codex MCP Alignment

**Date**: 2026-01-23

## Decision 1: MCP method alignment for Codex

- **Decision**: Implement standard MCP discovery and invocation methods (`initialize`, `tools/list`, `tools/call`) and keep `mcp.capabilities` for backward compatibility.
- **Rationale**: Codex expects MCP-standard discovery and tool invocation semantics; adding canonical methods maximizes compatibility without removing existing behavior.
- **Alternatives considered**: Continue with only `mcp.capabilities` and direct tool methods. Rejected because Codex clients may not recognize non-standard method names.

## Decision 2: Stdio entrypoint

- **Decision**: Add a dedicated stdio CLI entrypoint (`cmd/mcp`) that runs the MCP server over stdin/stdout.
- **Rationale**: Codex supports stdio MCP servers and requires a local command to launch them.
- **Alternatives considered**: HTTP-only configuration. Rejected because it excludes Codex stdio flow.

## Decision 3: HTTP/SSE request format

- **Decision**: Keep JSON-RPC 2.0 request bodies on `POST /mcp` and single-event SSE responses on `GET /mcp/sse` with a URL-encoded request payload; suppress responses for notifications (`id: null`).
- **Rationale**: Matches existing implementation, minimizes breaking change risk, and remains compatible with JSON-RPC semantics.
- **Alternatives considered**: Long-lived SSE streams and bidirectional streaming. Rejected for MVP scope.

## Decision 4: Tool schema representation

- **Decision**: Provide JSON Schema (Draft 2020-12 subset) for tool inputs/outputs using plain Go structs (no new dependency).
- **Rationale**: Keeps dependencies minimal while still enabling Codex to validate inputs and generate calls.
- **Alternatives considered**: Full schema libraries. Rejected due to added dependency weight.

## Decision 5: Smoke test scope

- **Decision**: Provide a repeatable smoke test that verifies discovery and a minimal tool call across stdio and HTTP transports (no external Codex tooling required).
- **Rationale**: Ensures MVP alignment without requiring Codex-specific automation.
- **Alternatives considered**: Rely only on manual Codex runs. Rejected because it is not repeatable or automatable.

## Implementation Notes (2026-01-23)

- Standard MCP methods, tool schemas, and stdio entrypoint were implemented as planned with no deviations.
