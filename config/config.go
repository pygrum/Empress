package options

import _ "embed"

//go:embed options.json
var options []byte

var (
	AgentID string
	C2Host  string
	C2Port  string
)
