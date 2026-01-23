param(
  [string]$BaseUrl = "http://localhost:8080"
)

$ErrorActionPreference = "Stop"

function Assert-NoError {
  param(
    [Parameter(Mandatory = $true)]$Response,
    [Parameter(Mandatory = $true)][string]$Context
  )

  if ($null -ne $Response.error) {
    throw "MCP $Context failed: $($Response.error | ConvertTo-Json -Compress)"
  }
}

function Invoke-McpHttp {
  param(
    [Parameter(Mandatory = $true)][string]$Method,
    [Parameter(Mandatory = $true)]$Params,
    [Parameter(Mandatory = $true)][int]$Id
  )

  $payload = @{ jsonrpc = "2.0"; id = $Id; method = $Method; params = $Params } | ConvertTo-Json -Compress
  return Invoke-RestMethod -Uri "$BaseUrl/mcp" -Method Post -ContentType "application/json" -Body $payload
}

Write-Host "[HTTP] initialize"
$initResp = Invoke-McpHttp -Method "initialize" -Params @{ protocol_version = "1.0" } -Id 1
Assert-NoError -Response $initResp -Context "initialize"

Write-Host "[HTTP] tools/list"
$listResp = Invoke-McpHttp -Method "tools/list" -Params @{} -Id 2
Assert-NoError -Response $listResp -Context "tools/list"
if (-not $listResp.result.tools -or $listResp.result.tools.Count -eq 0) {
  throw "MCP tools/list returned no tools"
}

Write-Host "[HTTP] tools/call file.write"
$writeResp = Invoke-McpHttp -Method "tools/call" -Params @{ name = "file.write"; arguments = @{ path = "mcp-smoke-script.txt"; content = "smoke" } } -Id 3
Assert-NoError -Response $writeResp -Context "tools/call file.write"

Write-Host "[HTTP] tools/call file.read"
$readResp = Invoke-McpHttp -Method "tools/call" -Params @{ name = "file.read"; arguments = @{ path = "mcp-smoke-script.txt" } } -Id 4
Assert-NoError -Response $readResp -Context "tools/call file.read"
if ($readResp.result.structuredContent.content -ne "smoke") {
  throw "MCP file.read returned unexpected content: $($readResp.result.structuredContent.content)"
}

Write-Host "[STDIO] initialize/tools/list/tools/call"
$stdioRequests = @(
  (@{ jsonrpc = "2.0"; id = 1; method = "initialize"; params = @{ protocol_version = "1.0" } } | ConvertTo-Json -Compress),
  (@{ jsonrpc = "2.0"; id = 2; method = "tools/list"; params = @{} } | ConvertTo-Json -Compress),
  (@{ jsonrpc = "2.0"; id = 3; method = "tools/call"; params = @{ name = "file.write"; arguments = @{ path = "mcp-smoke-stdio.txt"; content = "smoke" } } } | ConvertTo-Json -Compress),
  (@{ jsonrpc = "2.0"; id = 4; method = "tools/call"; params = @{ name = "file.read"; arguments = @{ path = "mcp-smoke-stdio.txt" } } } | ConvertTo-Json -Compress)
) -join "`n"

$stdioOutput = $stdioRequests | go run ./cmd/mcp
$lines = $stdioOutput -split "`r?`n" | Where-Object { $_ -ne "" }
if ($lines.Count -lt 4) {
  throw "STDIO expected 4 responses, got $($lines.Count)"
}

$stdioInit = $lines[0] | ConvertFrom-Json
Assert-NoError -Response $stdioInit -Context "stdio initialize"
$stdioList = $lines[1] | ConvertFrom-Json
Assert-NoError -Response $stdioList -Context "stdio tools/list"
if (-not $stdioList.result.tools -or $stdioList.result.tools.Count -eq 0) {
  throw "STDIO tools/list returned no tools"
}
$stdioWrite = $lines[2] | ConvertFrom-Json
Assert-NoError -Response $stdioWrite -Context "stdio tools/call file.write"
$stdioRead = $lines[3] | ConvertFrom-Json
Assert-NoError -Response $stdioRead -Context "stdio tools/call file.read"
if ($stdioRead.result.structuredContent.content -ne "smoke") {
  throw "STDIO file.read returned unexpected content: $($stdioRead.result.structuredContent.content)"
}

Write-Host "MCP smoke test passed (HTTP + STDIO)"
