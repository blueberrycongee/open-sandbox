package remote

type ServerConfig struct {
	Name               string            `json:"name"`
	URL                string            `json:"url"`
	Transport          string            `json:"transport,omitempty"`
	AuthorizationToken string            `json:"authorization_token,omitempty"`
	Headers            map[string]string `json:"headers,omitempty"`
	ToolAllow          []string          `json:"tool_allow,omitempty"`
	ToolDeny           []string          `json:"tool_deny,omitempty"`
}

type ConfigFile struct {
	Servers []ServerConfig `json:"servers"`
}
