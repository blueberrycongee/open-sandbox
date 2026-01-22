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

### WSL Environment Check

Command:
```
uname -a
```

Output:
```
Linux DESKTOP-MV4TM6D 6.6.87.2-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Thu Jun  5 18:30:46 UTC 2025 x86_64 x86_64 x86_64 GNU/Linux
```

### WSL Test Attempt

Command:
```
go test ./...
```

Output:
```
bash: line 1: go: command not found
```

Notes:
- Go is not installed in the WSL environment, so tests could not be executed there.

### WSL Test Run (After Go 1.24 Install)

Command:
```
export PATH=/usr/lib/go-1.24/bin:$PATH
export GOTOOLCHAIN=local
export GOMODCACHE=/mnt/c/Users/10758/go/pkg/mod
export GOPROXY=off
export GOSUMDB=off
go test ./...
```

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
ok  	open-sandbox/tests/integration	0.031s
ok  	open-sandbox/tests/unit	(cached)
```
