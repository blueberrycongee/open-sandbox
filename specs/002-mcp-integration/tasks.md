# Tasks: MCP Integration

**Input**: `specs/002-mcp-integration/spec.md`, `specs/002-mcp-integration/plan.md`
**Tests**: Required (strict TDD)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel
- **[Story]**: User story ID

---

## Phase 1: MCP Core (Blocking)

- [ ] T001 [P] Unit tests for JSON-RPC parsing + protocol version rejection + unified error schema in `tests/unit/mcp_types_test.go`
- [ ] T002 Create MCP JSON-RPC types in `internal/mcp/types.go`
- [ ] T003 Implement stdio server loop in `internal/mcp/server.go`
- [ ] T004 Implement protocol version validation with traceable errors in `internal/mcp/server.go`
- [ ] T005 Implement unified error schema and JSON-RPC error `data` details in `internal/mcp/types.go`
- [ ] T006 Implement `mcp.capabilities` in `internal/mcp/handlers.go`
- [ ] T007 [P] Unit tests for `mcp.capabilities` payload (versions + permissions metadata) in `tests/unit/mcp_capabilities_test.go`

---

## Phase 2: Tool Registry

- [ ] T008 [P] Unit tests for registry + permission metadata schema in `tests/unit/mcp_registry_test.go`
- [ ] T009 Implement tool registry + permission metadata in `internal/mcp/registry.go`

---

## Phase 3: Tool Mappings

- [ ] T010 [P] Integration test for MCP tool calls + workspace boundary in `tests/integration/mcp_tools_test.go`
- [ ] T011 Implement `browser.navigate` and `browser.screenshot` in `internal/mcp/tools/browser.go`
- [ ] T012 Implement `file.*` tools in `internal/mcp/tools/file.go`
- [ ] T013 Implement `shell.exec` tool in `internal/mcp/tools/shell.go`
- [ ] T014 Implement `code.exec` tool in `internal/mcp/tools/code.go`

---

## Phase 4: HTTP Transport

- [ ] T015 [P] Integration test for HTTP round-trip at `POST /mcp` + auth enabled/disabled in `tests/integration/mcp_http_test.go`
- [ ] T016 Implement HTTP endpoint for MCP requests/responses in `internal/mcp/server.go`

---

## Phase 5: SSE Transport

- [ ] T017 [P] Integration test for SSE round-trip at `GET /mcp/sse` + auth enabled/disabled in `tests/integration/mcp_sse_test.go`
- [ ] T018 Implement SSE endpoint for MCP requests/responses in `internal/mcp/server.go`

---

## Phase 6: Auth Controls

- [ ] T019 [P] Unit tests for auth toggles + JWT config validation in `tests/unit/mcp_auth_test.go`
- [ ] T020 Implement `MCP_AUTH_ENABLED` + JWT validation for HTTP/SSE in `internal/mcp/auth.go`

---

## Phase 7: Docs

- [ ] T021 Update `README.md` with MCP usage and examples
- [ ] T022 Add MCP usage to `specs/open-sandbox-mvp/quickstart.md`
- [ ] T023 Add README NFR verification checklist (minimal deps, no breaking HTTP APIs, demo-ready single machine, strict TDD, SANDBOX_ROOT usage)
