package remote

import "path"

type ToolFilter struct {
	Allow     []string
	Deny      []string
	AllowGlob []string
	DenyGlob  []string
}

func (filter ToolFilter) Allowed(name string) bool {
	hasAllow := len(filter.Allow) > 0 || len(filter.AllowGlob) > 0
	if hasAllow {
		if matchesAny(name, filter.Allow, filter.AllowGlob) {
			return true
		}
		return false
	}
	if matchesAny(name, filter.Deny, filter.DenyGlob) {
		return false
	}
	return true
}

func matchesAny(name string, exact []string, globs []string) bool {
	for _, candidate := range exact {
		if candidate == name {
			return true
		}
	}
	for _, pattern := range globs {
		if pattern == "" {
			continue
		}
		matched, err := path.Match(pattern, name)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
