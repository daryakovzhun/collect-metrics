package handler

import (
	"embed"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
	"html/template"
	"io/fs"
)

//go:embed templates/*
var templatesFS embed.FS

type Config struct {
	Key string
}

type Handler struct {
	cfg         *Config
	domain      controller.IServerController
	metricsTmpl *template.Template
}

func New(cfg *Config, domain controller.IServerController) *Handler {
	tmplFS, _ := fs.Sub(templatesFS, "templates")

	return &Handler{
		cfg:    cfg,
		domain: domain,
		metricsTmpl: template.Must(
			template.
				New("metrics.html").
				ParseFS(tmplFS, "metrics.html"),
		),
	}
}
