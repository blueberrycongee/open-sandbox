# Feature Specification: MCP Integration

**Feature Branch**: `002-mcp-integration`
**Created**: 2026-01-22
**Status**: Draft
**Input**: User request: "MCP interface for open-sandbox"

## Positioning

Expose sandbox capabilities through a minimal, extensible MCP server so LLM agents can invoke tools safely and incrementally.

## User Scenarios & Testing (mandatory)

### User Story 1 - MCP capability discovery (Priority: P1)

As a client, I can query MCP capabilities to discover available tools and permissions.

**Why this priority**: MCP clients need a stable entry point to understand tool availability.

**Independent Test**: `mcp.capabilities` returns a tool list with versions and permissions.

**Acceptance Scenarios**:

1. **Given** MCP server is running, **When** `mcp.capabilities` is requested, **Then** return tool list and metadata.
2. **Given** MCP server is running, **When** an unknown method is requested, **Then** return a unified error response.

---

### User Story 2 - Browser tools via MCP (Priority: P1)

As a client, I can navigate and capture screenshots through MCP tools.

**Why this priority**: Browser actions are core to sandbox value.

**Independent Test**: Navigate then capture a screenshot under `<SANDBOX_WORKSPACE>`.

**Acceptance Scenarios**:

1. **Given** browser is available, **When** `browser.navigate` is called with a URL, **Then** navigation succeeds.
2. **Given** navigation succeeded, **When** `browser.screenshot` is called with a path under `<SANDBOX_WORKSPACE>`, **Then** file exists.

---

### User Story 3 - File/Shell/Code tools via MCP (Priority: P1)

As a client, I can read/write files and execute shell/code through MCP tools.

**Why this priority**: Enables end-to-end workflows for LLM agents.

**Independent Test**: Write -> read -> process -> write -> read using MCP calls.

**Acceptance Scenarios**:

1. **Given** MCP is running, **When** `file.write` writes to `<SANDBOX_WORKSPACE>`, **Then** `file.read` returns content.
2. **Given** MCP is running, **When** `shell.exec` runs `echo test`, **Then** stdout includes `test`.
3. **Given** MCP is running, **When** `code.exec` processes a file, **Then** output file exists under `<SANDBOX_WORKSPACE>`.

---

## Requirements (mandatory)

### Functional Requirements

- **FR-001**: Provide MCP server using JSON-RPC 2.0 over stdio (primary transport).
- **FR-001a**: Provide an HTTP transport endpoint for MCP requests/responses.
- **FR-001b**: Provide an SSE (Server-Sent Events) transport for MCP requests/responses.
- **FR-001c**: Declare the supported MCP protocol version explicitly; reject incompatible versions with a traceable error.
- **FR-001d**: HTTP JSON-RPC endpoint path must be `POST /mcp`.
- **FR-001e**: SSE endpoint path must be `GET /mcp/sse`.
- **FR-002**: Implement `mcp.capabilities` with tool list, versions, and permissions.
- **FR-003**: Implement `browser.navigate` and `browser.screenshot` tools.
- **FR-004**: Implement `file.read/write/list/search/replace` tools.
- **FR-005**: Implement `shell.exec` tool.
- **FR-006**: Implement `code.exec` tool.
- **FR-007**: Enforce workspace boundary for all file and screenshot paths.
- **FR-008**: All MCP responses use unified error schema (code/message/trace_id).
- **FR-009**: Support `Authorization: Bearer <token>` on HTTP/SSE when auth is enabled; ignore when disabled.
- **FR-010**: Include tool-level permission metadata in `mcp.capabilities` using a minimal schema: `allow` (bool), `scope` (string: `workspace`|`network`|`exec`), optional `reason` (string).
- **FR-011**: JSON-RPC error `data` must include unified error details (e.g., trace_id, kind).
- **FR-012**: Auth toggles must include `MCP_AUTH_ENABLED` (default `false`) and one of `MCP_AUTH_JWT_SECRET` or `MCP_AUTH_JWT_PUBLIC_KEY`. Optional: `MCP_AUTH_AUDIENCE`, `MCP_AUTH_ISSUER`.

### Non-Functional Requirements

- **NFR-001**: Minimal dependencies (standard library preferred).
- **NFR-002**: No breaking changes to existing HTTP APIs.
- **NFR-003**: Demo-ready on a single machine.
- **NFR-004**: Strict TDD; tests before implementation.
- **NFR-005**: Runtime artifacts stay under `SANDBOX_ROOT`.

## Success Criteria (mandatory)

- **SC-001**: `mcp.capabilities` returns tool list and metadata.
- **SC-002**: `browser.navigate` + `browser.screenshot` produce a screenshot under `<SANDBOX_WORKSPACE>`.
- **SC-003**: `file.read/write` work under `<SANDBOX_WORKSPACE>`.
- **SC-004**: `shell.exec` returns stdout/stderr/exit code.
- **SC-005**: `code.exec` writes output under `<SANDBOX_WORKSPACE>`.

## Assumptions & Open Questions

- [NEEDS CLARIFICATION] Client compatibility details for HTTP/SSE transport.
- [NEEDS CLARIFICATION] Do we need explicit allow/deny policies per tool?
- [NEEDS CLARIFICATION] Default `SANDBOX_WORKSPACE` values by platform (Windows vs. Linux/WSL) and how overrides are configured.
