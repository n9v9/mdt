package main

import (
	"context"
	"io"
)

type ctxKeyType int

const ctxKeyInputFile ctxKeyType = iota

func withInputFile(ctx context.Context, r io.ReadCloser) context.Context {
	return context.WithValue(ctx, ctxKeyInputFile, r)
}

func inputFile(ctx context.Context) io.ReadCloser {
	return ctx.Value(ctxKeyInputFile).(io.ReadCloser)
}
