package gentlemanlogger

import (
	"io"

	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// New creates logger plugin instance
func New(out io.Writer) plugin.Plugin {
	p := plugin.New()
	p.SetHandler("request", func(c *context.Context, h context.Handler) {
		h.Next(c)
	})
	p.SetHandler("response", func(c *context.Context, h context.Handler) {
		h.Next(c)
	})
	return p
}
