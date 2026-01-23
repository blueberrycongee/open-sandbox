# Data Model: Codex MCP Alignment

## Entities

### MCP Server

- **Description**: Entry point for Codex MCP clients.
- **Attributes**:
  - `protocol_version` (string)
  - `transports` (list: stdio, http, sse)
  - `auth_enabled` (boolean)
- **Validation**:
  - Reject incompatible protocol versions.
  - Enforce bearer token checks when auth is enabled.

### Tool

- **Description**: Callable capability exposed to MCP clients.
- **Attributes**:
  - `name` (string)
  - `version` (string)
  - `permissions` (allow, scope, reason)
  - `schema` (input/output JSON schema)

### Tool Schema

- **Description**: Declarative input/output schema for tool invocation.
- **Attributes**:
  - `input` (JSON Schema object)
  - `output` (JSON Schema object)

### Tool Invocation

- **Description**: A single call from a client to a tool.
- **Attributes**:
  - `method` (string)
  - `params` (object)
  - `result` (object)
  - `error` (code/message/trace_id)

### Transport Mode

- **Description**: Connection type used by a client.
- **Attributes**:
  - `type` (stdio, http, sse)
  - `endpoint` (string, optional for stdio)

### Auth Configuration

- **Description**: Auth settings for MCP HTTP/SSE.
- **Attributes**:
  - `enabled` (boolean)
  - `bearer_token` (string)
  - `audience` (string, optional)
  - `issuer` (string, optional)
