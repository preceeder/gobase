/*
File Name:  jsonHandler.go
Description:
Author:      Chenghu
Date:       2023/9/1 18:05
Change Activity:
*/
package logs

import (
	"golang.org/x/net/context"
	"io"
	"log/slog"
)

// JSONHandler is a Handler that writes Records to an io.Writer as
// line-delimited JSON objects.
type JsonHandler struct {
	slog.Handler
}

// NewJSONHandler creates a JSONHandler that writes to w,
// using the given options.
// If opts is nil, the default options are used.
func NewJsonHandler(w io.Writer, opts *slog.HandlerOptions) *JsonHandler {
	v := &JsonHandler{}
	v.Handler = slog.NewJSONHandler(w, opts)
	return v
}

func (h *JsonHandler) Enabled(c context.Context, level slog.Level) bool {
	return h.Handler.Enabled(c, level)
}

func (h *JsonHandler) Handle(c context.Context, r slog.Record) error {
	err := h.Handler.Handle(c, r)
	if err != nil {
		return err
	}
	return nil
}

func (h *JsonHandler) WithAttrs(as []slog.Attr) slog.Handler {
	return &JsonHandler{
		Handler: h.Handler.WithAttrs(as),
	}
}

func (h *JsonHandler) WithGroup(name string) slog.Handler {
	return &JsonHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
