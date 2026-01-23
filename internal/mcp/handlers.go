package mcp

type CapabilitiesResponse struct {
	ProtocolVersion string     `json:"protocol_version"`
	Tools           []ToolInfo `json:"tools"`
}

func BuildCapabilities(registry *Registry) CapabilitiesResponse {
	return CapabilitiesResponse{
		ProtocolVersion: SupportedProtocolVersion,
		Tools:           registry.List(),
	}
}
