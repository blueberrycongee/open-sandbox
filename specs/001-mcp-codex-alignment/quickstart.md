# Quickstart: Codex MCP Alignment

## Local Setup (HTTP/SSE)

1) Run the sandbox server:

```powershell
go run ./cmd/server
```

2) Configure Codex MCP in `~/.codex/config.toml`:

```toml
[mcp_servers.open_sandbox]
url = "http://localhost:8080/mcp"
```

3) Optional SSE transport:

```toml
[mcp_servers.open_sandbox_sse]
url = "http://localhost:8080/mcp/sse"
```

4) Optional auth (HTTP/SSE only):

```toml
[mcp_servers.open_sandbox]
bearer_token_env_var = "MCP_AUTH_TOKEN"
```

## Local Setup (STDIO)

1) Build and run the MCP stdio entrypoint:

```powershell
go run ./cmd/mcp
```

2) Configure Codex MCP in `~/.codex/config.toml`:

```toml
[mcp_servers.open_sandbox_stdio]
command = "go"
args = ["run", "./cmd/mcp"]
```

## Smoke Test (HTTP)

```powershell
# Initialize
Invoke-RestMethod -Uri http://localhost:8080/mcp -Method Post -ContentType application/json -Body '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocol_version":"1.0"}}'

# Tool discovery
Invoke-RestMethod -Uri http://localhost:8080/mcp -Method Post -ContentType application/json -Body '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

## Smoke Test (STDIO)

```powershell
# Pipe JSON-RPC requests into the stdio server
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocol_version":"1.0"}}' | go run ./cmd/mcp
'{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | go run ./cmd/mcp
```
