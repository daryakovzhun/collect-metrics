package handler

import (
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
)

type Handler struct {
	domain controller.IController
}

func New(domain controller.IController) *Handler {
	return &Handler{
		domain: domain,
	}
}
