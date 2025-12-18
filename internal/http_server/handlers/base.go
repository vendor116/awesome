package handlers

import (
	"github.com/vendor116/awesome/pkg/openapi"
)

var _ openapi.StrictServerInterface = (*BaseHandler)(nil)

type BaseHandler struct {
	VersionHandler
}

func NewBaseHandler() BaseHandler {
	return BaseHandler{
		VersionHandler: NewVersionHandler(),
	}
}
