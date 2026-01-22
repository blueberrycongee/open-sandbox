# Tasks: open-sandbox MVP

**Input**: Design documents from `.specify/memory/` (`spec.md`, `plan.md`, `constitution.md`)
**Tests**: Required (strict TDD). All tests must be written first and fail before implementation.
**Organization**: Tasks grouped by milestone and user story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story ID (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Foundation (Blocking)

**Purpose**: Core infrastructure required before any user story work

- [X] T001 Create Go module and folder layout per plan in `cmd/`, `internal/`, `pkg/`, `tests/`
- [X] T002 [P] Add `.editorconfig` and `.gitattributes` in repo root
- [X] T003 [P] Add Makefile with `build` and `test` targets in `Makefile`
- [X] T004 Implement server entrypoint and HTTP wiring in `cmd/server/main.go`
- [X] T005 Define unified response and error types in `pkg/types/response.go`
- [X] T006 Implement error mapping + trace IDs in `internal/api/errors.go`
- [X] T007 [P] Add minimal logging wrapper (std `log`) in `internal/api/logging.go`
- [X] T008 Implement router/mux with net/http only in `internal/api/router.go`
- [X] T009 Implement error middleware for uniform JSON errors in `internal/api/middleware.go`
- [X] T010 Add absolute path config constants and workspace creation in `internal/config/paths.go`
- [X] T011 [P] Unit test response schema serialization in `tests/unit/response_test.go`
- [X] T012 [P] Unit test error mapping and trace ID behavior in `tests/unit/errors_test.go`
- [X] T013 [P] Unit test workspace path creation in `tests/unit/workspace_test.go`

**Checkpoint**: Foundation complete; user story work can begin

---

## Phase 2: User Story 1 - Unified entry API (P1)

**Goal**: `GET /v1/sandbox` returns capability list and health status

### Tests (write first)

- [X] T014 [P] [US1] Integration test for `GET /v1/sandbox` in `tests/integration/sandbox_index_test.go`
- [X] T015 [P] [US1] Integration test for unknown route error format in `tests/integration/not_found_test.go`

### Implementation

- [X] T016 [US1] Implement sandbox index handler in `internal/api/handlers/sandbox.go`
- [X] T017 [US1] Wire route `GET /v1/sandbox` in `internal/api/router.go`

**Checkpoint**: US1 independently testable

---

## Phase 3: User Story 2 - Browser + VNC (P1)

**Goal**: Headed browser control with VNC visual takeover

### Tests (write first)

- [X] T018 [P] [US2] Integration test for browser info (CDP address) in `tests/integration/browser_info_test.go`
- [X] T019 [P] [US2] Integration test for screenshot flow in `tests/integration/browser_screenshot_test.go`
- [X] T020 [P] [US2] Integration test for VNC page accessible in `tests/integration/vnc_page_test.go`

### Implementation

- [X] T021 [US2] Implement browser service scaffold in `internal/browser/browser.go`
- [X] T022 [US2] Implement browser handlers (info, navigate, screenshot) in `internal/api/handlers/browser.go`
- [X] T023 [US2] Implement VNC handler and static assets in `internal/api/handlers/vnc.go`
- [X] T024 [US2] Wire browser/VNC routes in `internal/api/router.go`

**Checkpoint**: Browser + VNC independently testable

---

## Phase 4: User Story 3 - Shell/File/Code closed loop (P1)

**Goal**: Shell, File, and Code execution form a verifiable loop

### Tests (write first)

- [X] T025 [P] [US3] Integration test for Shell `echo test` in `tests/integration/shell_exec_test.go`
- [X] T026 [P] [US3] Integration test for File CRUD/search/replace in `tests/integration/file_api_test.go`
- [X] T027 [P] [US3] Integration test for Code exec (python/node) in `tests/integration/code_exec_test.go`
- [X] T028 [P] [US3] End-to-end flow test in `tests/integration/flow_e2e_test.go`

### Implementation

- [X] T029 [US3] Implement shell service + handler in `internal/shell/shell.go` and `internal/api/handlers/shell.go`
- [X] T030 [US3] Implement file service + handler in `internal/file/file.go` and `internal/api/handlers/file.go`
- [X] T031 [US3] Implement code exec service + handler in `internal/codeexec/codeexec.go` and `internal/api/handlers/codeexec.go`
- [X] T032 [US3] Wire shell/file/code routes in `internal/api/router.go`

**Checkpoint**: Shell/File/Code independently testable; E2E flow passes

---

## Phase 5: User Story 4 - Jupyter & Code Server entry (P2)

**Goal**: Jupyter Lab and code-server are reachable

### Tests (write first)

- [X] T033 [P] [US4] Integration test for Jupyter entry in `tests/integration/jupyter_entry_test.go`
- [X] T034 [P] [US4] Integration test for Code Server entry in `tests/integration/codeserver_entry_test.go`

### Implementation

- [X] T035 [US4] Implement Jupyter proxy/redirect handler in `internal/api/handlers/jupyter.go`
- [X] T036 [US4] Implement Code Server proxy/redirect handler in `internal/api/handlers/codeserver.go`
- [X] T037 [US4] Wire Jupyter/Code Server routes in `internal/api/router.go`

**Checkpoint**: Jupyter & Code Server independently testable

---

## Phase 6: Docs, Ops, and Acceptance

**Purpose**: Documentation and acceptance alignment

- [X] T038 [P] Add README with ports/env/startup in `README.md`
- [X] T039 [P] Add limitations/TODO with assumptions in `README.md`
- [X] T040 [P] Add `.env.example` documenting JWT toggle in `.env.example`
- [X] T041 [P] Document runtime artifact locations under `D:\\Desktop\\sandbox\\open-sandbox` in `README.md`
- [X] T042 [P] Add quickstart doc in `specs/open-sandbox-mvp/quickstart.md`
- [X] T043 [P] Run integration tests and document results in `specs/open-sandbox-mvp/research.md`

---

## Dependencies & Execution Order

- Phase 1 blocks all other phases
- US1, US2, US3, US4 require Phase 1 complete
- Tests must be written and fail before implementation in each user story

## Notes

- Absolute paths only; host workspace is `D:\\Desktop\\sandbox\\open-sandbox\\workspace` (create if missing)
- If containerized, mount host workspace to `/workspace` and use `/workspace` inside the container
- All runtime artifacts (cache/logs/build) must stay on `D:\\Desktop\\sandbox\\open-sandbox`
