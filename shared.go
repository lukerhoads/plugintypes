package plugintypes

import (
	"context"

	"github.com/hashicorp/go-plugin"
)

type InitArgs struct {
	ctx context.Context
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
