# open-sandbox MVP Quickstart

## Local (Windows)

1) Ensure Go 1.24+ is installed.
2) From the repo root:

```powershell
go test ./...
go run ./cmd/server
```

3) Open endpoints:
- API: http://localhost:8080/v1/sandbox
- VNC placeholder: http://localhost:8080/vnc/index.html
- Jupyter placeholder: http://localhost:8080/jupyter
- Code-server placeholder: http://localhost:8080/code-server/
- MCP HTTP: http://localhost:8080/mcp
- MCP SSE: http://localhost:8080/mcp/sse

MCP example request (HTTP):
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "mcp.capabilities",
  "params": {
    "protocol_version": "1.0"
  }
}
```

## Docker (Placeholder)

```powershell
docker build -t open-sandbox .
docker run --rm -p 8080:8080 -v <SANDBOX_WORKSPACE>:/workspace open-sandbox
```

## Environment

- `SANDBOX_ADDR` sets the HTTP listen address (default `:8080`)
- `SANDBOX_BROWSER_CDP` sets the CDP websocket address returned by the browser info API
- `SANDBOX_ROOT` sets the base path for runtime artifacts
- `SANDBOX_WORKSPACE` overrides the workspace path
