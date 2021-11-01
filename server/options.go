package server

import (
	"samwang0723/jarvis/handlers"
	structuredlog "samwang0723/jarvis/logger/structured"
)

type Options struct {
	Name             string
	Logger           structuredlog.ILogger
	Handler          handlers.IHandler
	ProfilingEnabled bool
	Config           interface{}

	// Before funcs
	BeforeStart []func() error
	BeforeStop  []func() error
}

type Option func(o *Options)

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func Handler(handler handlers.IHandler) Option {
	return func(o *Options) {
		o.Handler = handler
	}
}

func Logger(logger structuredlog.ILogger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}
