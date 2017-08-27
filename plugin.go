package gentlemanlogger

import (
	"io"

	"github.com/izumin5210/httplogger"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// New creates logger plugin instance
func New(out io.Writer) plugin.Plugin {
	p := plugin.New()
	logger := httplogger.NewLogger(out)
	p.SetHandler("request", func(c *context.Context, h context.Handler) {
		c.Request = logger.LogRequest(c.Request)
		h.Next(c)
	})
	p.SetHandler("response", func(c *context.Context, h context.Handler) {
		logger.LogResponse(c.Response)
		h.Next(c)
	})
	return p
}
