# Research Log: open-sandbox MVP

## 2026-01-22

### Integration Test Run

Command:
```
go test ./...
```

Result: PASS (unit + integration)

Output:
```
?   	open-sandbox/cmd/server	[no test files]
?   	open-sandbox/internal/api	[no test files]
?   	open-sandbox/internal/api/handlers	[no test files]
?   	open-sandbox/internal/browser	[no test files]
?   	open-sandbox/internal/codeexec	[no test files]
?   	open-sandbox/internal/config	[no test files]
?   	open-sandbox/internal/file	[no test files]
?   	open-sandbox/internal/shell	[no test files]
?   	open-sandbox/pkg/types	[no test files]
ok  	open-sandbox/tests/integration	1.954s
ok  	open-sandbox/tests/unit	0.017s
```

Notes:
- Browser, VNC, Jupyter, and code-server are exercised via real handlers (no placeholder HTML).
- Browser tests run in headless mode for stability during automated runs.
- Code execution tests skip if Python/Node is not installed.
