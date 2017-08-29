package logger

import (
	"io"

	"github.com/izumin5210/httplogger"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// New creates logger plugin instance
func New(out io.Writer) plugin.Plugin {
	return plugin.NewRequestPlugin(func(c *context.Context, h context.Handler) {
		c.Client.Transport = httplogger.NewRoundTripper(out, c.Client.Transport)
		h.Next(c)
	})
}
