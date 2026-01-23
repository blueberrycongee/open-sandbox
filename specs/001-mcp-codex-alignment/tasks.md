# Tasks: Codex MCP Alignment

**Input**: Design documents from `/specs/001-mcp-codex-alignment/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required (strict TDD from constitution)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create MCP stdio entrypoint directory in `cmd/mcp/`
- [x] T002 [P] Create MCP stdio command scaffold in `cmd/mcp/main.go`
- [x] T003 [P] Add Codex MCP configuration stub in `specs/001-mcp-codex-alignment/quickstart.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Define MCP standard method constants and request/response structs in `internal/mcp/types.go`
- [x] T005 [P] Extend MCP error helpers for standard methods in `internal/mcp/types.go`
- [x] T006 Add schema container types for tools in `internal/mcp/types.go`
- [x] T007 Update MCP registry to store tool schemas in `internal/mcp/registry.go`
- [x] T008 Add smoke test helpers for MCP in `tests/integration/mcp_smoke_helpers_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Codex connects and discovers tools (Priority: P1) 🎯 MVP

**Goal**: Implement MCP discovery with schemas and versions for Codex tool listing

**Independent Test**: Codex-compatible tool discovery returns tool list with schemas and versions

### Tests for User Story 1 (REQUIRED) ⚠️

- [x] T009 [P] [US1] Unit test for standard discovery methods in `tests/unit/mcp_discovery_test.go`
- [x] T010 [P] [US1] Integration test for tool discovery over HTTP in `tests/integration/mcp_discovery_http_test.go`
- [x] T011 [P] [US1] Integration test for tool discovery over stdio in `tests/integration/mcp_discovery_stdio_test.go`

### Implementation for User Story 1

- [x] T012 [US1] Implement `initialize` handler in `internal/mcp/server.go`
- [x] T013 [US1] Implement `tools/list` handler in `internal/mcp/server.go`
- [x] T014 [US1] Add tool schema population for existing tools in `internal/api/handlers/mcp.go`
- [x] T015 [US1] Add step-by-step Codex configuration guidance in `README.md`

**Checkpoint**: User Story 1 is independently functional and testable

---

## Phase 4: User Story 2 - Codex executes sandbox tools (Priority: P1)

**Goal**: Implement standard tool invocation pathway and consistent results

**Independent Test**: Standard `tools/call` executes file and shell tools successfully

### Tests for User Story 2 (REQUIRED) ⚠️

- [x] T016 [P] [US2] Unit test for `tools/call` routing in `tests/unit/mcp_tools_call_test.go`
- [x] T017 [P] [US2] Integration test for tool execution over HTTP in `tests/integration/mcp_tools_call_http_test.go`
- [x] T018 [P] [US2] Integration test for tool execution over stdio in `tests/integration/mcp_tools_call_stdio_test.go`

### Implementation for User Story 2

- [x] T019 [US2] Implement `tools/call` handler in `internal/mcp/server.go`
- [x] T020 [US2] Map standard tool invocation to existing handlers in `internal/mcp/tools/` (update relevant files)
- [x] T021 [US2] Ensure unified errors for tool invocation in `internal/mcp/server.go`

**Checkpoint**: User Story 2 is independently functional and testable

---

## Phase 5: User Story 3 - Codex connects via local process or network transport (Priority: P2)

**Goal**: Provide stdio entrypoint and repeatable smoke test across transports

**Independent Test**: Run tool discovery and execution via stdio and HTTP with identical results

### Tests for User Story 3 (REQUIRED) ⚠️

- [x] T022 [P] [US3] Integration smoke test for HTTP transport in `tests/integration/mcp_smoke_http_test.go`
- [x] T023 [P] [US3] Integration smoke test for stdio transport in `tests/integration/mcp_smoke_stdio_test.go`

### Implementation for User Story 3

- [x] T024 [US3] Implement stdio server command in `cmd/mcp/main.go`
- [x] T025 [US3] Add stdio wiring to MCP server in `internal/mcp/server.go`
- [x] T026 [US3] Document stdio setup in `specs/001-mcp-codex-alignment/quickstart.md`

**Checkpoint**: User Story 3 is independently functional and testable

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T027 [P] Add repeatable smoke test script in `scripts/mcp_smoke.ps1`
- [x] T028 [P] Update MCP contract schema to include standard discovery methods in `specs/001-mcp-codex-alignment/contracts/mcp-http.yaml`
- [x] T029 [P] Update feature research notes if implementation deviates in `specs/001-mcp-codex-alignment/research.md`
- [x] T030 Run quickstart smoke steps in `specs/001-mcp-codex-alignment/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Schemas before discovery handlers
- Discovery before tool invocation
- Core implementation before integration

### Parallel Opportunities

- Setup tasks T002 and T003 can run in parallel
- Foundational tasks T005, T006, T007, T008 can run in parallel after T004
- Story tests within each story can run in parallel
- Story phases can proceed in parallel once foundations are complete

---

## Parallel Example: User Story 1

```bash
Task: "Unit test for standard discovery methods in tests/unit/mcp_discovery_test.go"
Task: "Integration test for tool discovery over HTTP in tests/integration/mcp_discovery_http_test.go"
Task: "Integration test for tool discovery over stdio in tests/integration/mcp_discovery_stdio_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
