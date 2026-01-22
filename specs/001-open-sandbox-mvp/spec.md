# Feature Specification: open-sandbox MVP

**Feature Branch**: `mvp/spec`
**Created**: 2026-01-22
**Status**: Draft
**Input**: User description: "/speckit.specify open-sandbox MVP spec"

## Positioning

- Build an integrated sandbox for AI/human co-development with browser, terminal, file, IDE, Jupyter, and code execution.
- MVP must be demo-ready and usable; deployment should be simple (single machine/single container).

## User Scenarios & Testing (mandatory)

### User Story 1 - Unified entry API for sandbox capabilities (Priority: P1)

Users discover and call sandbox capabilities (Browser/VNC/IDE/Jupyter/Shell/File/Code) via a single HTTP entrypoint to complete an end-to-end task.

**Why this priority**: The unified entry is the starting point for all capabilities; without it, the MVP demo cannot run.

**Independent Test**: `GET /v1/sandbox` returns a capability list and health status.

**Acceptance Scenarios**:

1. **Given** service is running, **When** `GET /v1/sandbox`, **Then** return a unified response with status/data/error (or equivalent).
2. **Given** service is running, **When** request unknown resource, **Then** return a unified, traceable error response.

---

### User Story 2 - Headed browser + VNC visual control (Priority: P1)

Users drive a headed browser via API and can visually take over via VNC for navigation and screenshots.

**Why this priority**: Browser + VNC are the most visible MVP capabilities.

**Independent Test**: Navigate and take a screenshot; VNC page is accessible and interactive.

**Acceptance Scenarios**:

1. **Given** browser is available, **When** navigate to `https://example.com`, **Then** return success and VNC shows the page.
2. **Given** page is open, **When** capture screenshot to `<SANDBOX_WORKSPACE>/screenshots/example.png`, **Then** file exists and File API can read metadata. If running in container, the equivalent path is `/workspace/screenshots/example.png`.

---

### User Story 3 - Shell/File/Code closed loop (Priority: P1)

Users use Shell, File, and Code execution to read/write/process files and verify output.

**Why this priority**: A closed-loop flow is the minimal proof of value.

**Independent Test**: Run Shell -> File -> Code -> File and verify output.

**Acceptance Scenarios**:

1. **Given** Shell API is available, **When** run `echo test`, **Then** stdout includes `test`.
2. **Given** `<SANDBOX_WORKSPACE>/output.txt` exists, **When** read via File API, **Then** return verifiable content. If running in container, the equivalent path is `/workspace/output.txt`.

---

### User Story 4 - Access Jupyter Lab & Code Server (Priority: P2)

Users can reach Jupyter Lab and code-server via the unified entry.

**Why this priority**: Improves demo value but not required for minimal loop.

**Independent Test**: Access endpoints and receive valid pages.

**Acceptance Scenarios**:

1. **Given** service is running, **When** `http://localhost:8080/jupyter`, **Then** page loads.
2. **Given** service is running, **When** `http://localhost:8080/code-server/`, **Then** page loads.

---

### Edge Cases

- File API path not found, permission denied, binary reads: error format?
- Shell/Code timeout, command not found, non-zero exit: error format?
- Browser not ready or CDP connection failed: error format?
- VNC not ready or disconnected: error format?
- Large file read/write and screenshot size limits?

## Requirements (mandatory)

### Functional Requirements

- **FR-001**: Provide a unified entry HTTP API for capability discovery and health.
- **FR-002**: Provide a headed browser with CDP address discovery and external connectivity.
- **FR-003**: Provide browser action APIs (at least navigate and screenshot).
- **FR-004**: Provide VNC visual takeover entry (desktop visible and controllable).
- **FR-005**: Provide Shell API to execute non-interactive commands and return stdout/stderr/exit code.
- **FR-006**: Provide File API: read, write, list, search, replace.
- **FR-007**: Provide code execution (Python and Node minimal viability).
- **FR-008**: Provide access to Jupyter Lab and code-server.
- **FR-009**: All APIs return a consistent structure (status/data/error or equivalent).
- **FR-010**: Errors are traceable (error code/message/trace_id or equivalent).
- **FR-011**: Auth off by default, with JWT toggle placeholder.
- **FR-012**: Absolute paths only. Host workspace is configurable via `SANDBOX_WORKSPACE` (default `<repo_root>/workspace`, create if missing).
- **FR-013**: If a container is used, the host workspace must be mounted to `/workspace` and `/workspace` is used inside the container.

### Non-Functional Requirements

- **NFR-001**: Runs on Windows host and local Docker.
- **NFR-002**: Docs must list ports, environment variables, startup instructions, and limitations/TODO.
- **NFR-003**: No hard performance targets, but avoid obvious blocking (timeouts/async strategy documented).
- **NFR-004**: Atomic development and commits.
- **NFR-005**: Strict TDD: tests first, every implementation has tests.
- **NFR-006**: Deployment is single machine/single container and demo-ready for local use.

## Success Criteria (mandatory)

### Measurable Outcomes

- **SC-001**: After startup, these endpoints are accessible on Windows localhost:
  - `http://localhost:8080/v1/sandbox`
  - `http://localhost:8080/vnc/index.html`
  - `http://localhost:8080/jupyter`
  - `http://localhost:8080/code-server/`
- **SC-002**: End-to-end API flow with verifiable output:
  1) Open `https://example.com`
  2) Save screenshot to `<SANDBOX_WORKSPACE>/screenshots/example.png`
  3) File API reads screenshot metadata or saves page text to file
  4) Code exec processes file and writes `<SANDBOX_WORKSPACE>/output.txt`
  5) File API reads `<SANDBOX_WORKSPACE>/output.txt` and verifies content
  6) If running in container, the equivalent path is `/workspace/...`
- **SC-003**: Shell API runs `echo test` and returns output; `ls`/`dir` lists the workspace.
- **SC-004**: File API supports read/write/list/search/replace with consistent errors.
- **SC-005**: Browser info returns CDP address and external tool can connect and perform one `newPage` navigation.
- **SC-006**: Docs list ports, env vars, startup, limitations/known issues.
- **SC-007**: All runtime artifacts stay under `SANDBOX_ROOT` (default `<repo_root>`, e.g., `<SANDBOX_ROOT>/.cache`).

## Assumptions & Open Questions

- [NEEDS CLARIFICATION] Code execution isolation strategy (container/OS-level/none)?
- [NEEDS CLARIFICATION] Is proxy/port forwarding in scope for MVP?
- [NEEDS CLARIFICATION] Is MCP Hub required, and if yes, what interface?
