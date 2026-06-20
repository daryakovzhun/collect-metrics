package handler

import (
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
)

type Handler struct {
	domain controller.IServerController
}

func New(domain controller.IServerController) *Handler {
	return &Handler{
		domain: domain,
	}
}
