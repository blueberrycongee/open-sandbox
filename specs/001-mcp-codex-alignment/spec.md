# Feature Specification: Codex MCP Alignment

**Feature Branch**: `001-mcp-codex-alignment`  
**Created**: 2026-01-23  
**Status**: Draft  
**Input**: User description: "针对这个对齐 我们新开一个分支/pr"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Codex connects and discovers tools (Priority: P1)

As a Codex user, I can connect Codex to the local sandbox MCP server and discover available tools with their capabilities.

**Why this priority**: Tool discovery is the minimum viable step for any MCP client to operate.

**Independent Test**: Configure Codex once, connect successfully, and receive a tool list with metadata.

**Acceptance Scenarios**:

1. **Given** the MCP server is running, **When** Codex connects, **Then** it receives a capabilities response listing tools and permissions.
2. **Given** the MCP server is running, **When** Codex requests tool discovery, **Then** the response includes tool schemas and versions.

---

### User Story 2 - Codex executes sandbox tools (Priority: P1)

As a Codex user, I can call sandbox tools through MCP and receive correct results.

**Why this priority**: End-to-end tool execution is the core value of MCP integration.

**Independent Test**: Perform a file read/write and a shell execution via MCP and confirm outputs.

**Acceptance Scenarios**:

1. **Given** a valid MCP tool request, **When** Codex invokes a tool, **Then** the result matches the expected output.
2. **Given** an invalid tool request, **When** Codex invokes a tool, **Then** it receives a structured error response.

---

### User Story 3 - Codex connects via local process or network transport (Priority: P2)

As an operator, I can expose MCP via both local-process and network transports so Codex can connect in different environments.

**Why this priority**: Some Codex setups use local process transport, while others use network access.

**Independent Test**: Connect once via local-process transport and once via network transport without changing tool behavior.

**Acceptance Scenarios**:

1. **Given** local-process transport is configured, **When** Codex connects, **Then** tool discovery and execution work as expected.
2. **Given** network transport is configured, **When** Codex connects, **Then** tool discovery and execution work as expected.

---

### Edge Cases

- What happens when the client sends a notification request with a null id?
- How does the system handle unsupported protocol versions?
- How does the system respond to malformed requests or missing required fields?
- What happens when authentication is enabled but the request lacks a valid token?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The MCP server MUST support standard MCP discovery and tool invocation semantics required by Codex clients.
- **FR-002**: The MCP server MUST expose both local-process transport and network transport modes as defined by the MCP specification.
- **FR-003**: The MCP server MUST provide tool metadata including name, version, permissions, and input/output schemas.
- **FR-004**: The MCP server MUST allow tool execution with consistent results across transports.
- **FR-005**: The MCP server MUST return structured, traceable errors for invalid requests and tool failures.
- **FR-006**: The MCP server MUST reject incompatible protocol versions with a clear error response.
- **FR-007**: The MCP server MUST suppress responses for notification requests (null id).
- **FR-008**: When authentication is enabled, the MCP server MUST enforce bearer token checks for network transport; when disabled, it MUST ignore missing tokens.
- **FR-009**: The documentation MUST include step-by-step Codex configuration guidance for connecting to the MCP server.
- **FR-010**: The system MUST provide a repeatable smoke test that verifies Codex connectivity and tool execution across transports.
- **FR-011**: The MCP server MUST implement standard methods: `initialize`, `tools/list`, `tools/call`.
- **FR-012**: Tool schemas MUST be provided as JSON Schema (Draft 2020-12 subset) for both inputs and outputs.

### Non-Functional Requirements

- **NFR-001**: Strict TDD is required; tests MUST be written before implementation and must fail before code is added.

### Key Entities *(include if feature involves data)*

- **MCP Server**: The service that exposes tools and protocol behavior to MCP clients.
- **Tool**: A callable capability (browser/file/shell/code) with metadata and schemas.
- **Tool Schema**: The declared inputs and outputs required for tool invocation.
- **Transport Mode**: The connection method used by clients (local-process or network).
- **Auth Configuration**: The settings that enable or disable bearer token enforcement.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A first-time Codex user can connect and list tools within 5 minutes using the documented setup steps.
- **SC-002**: In a scripted smoke test, 3 consecutive runs of tool discovery and execution succeed without manual intervention.
- **SC-003**: At least 95% of tool invocations in the smoke test return successful results on first attempt.
- **SC-004**: Protocol errors (invalid request, unsupported version) consistently return structured error responses in 100% of test cases.

## Assumptions

- Existing sandbox tools (browser/file/shell/code) are suitable for MCP exposure without changing their core behavior.
- Codex MCP clients follow the published MCP specification for discovery and tool invocation.
