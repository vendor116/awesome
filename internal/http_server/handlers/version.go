package handlers

import (
	"context"

	"github.com/vendor116/awesome/pkg/openapi"
	"github.com/vendor116/awesome/pkg/version"
)

type VersionHandler struct{}

func NewVersionHandler() VersionHandler {
	return VersionHandler{}
}

func (vh VersionHandler) GetVersion(
	_ context.Context,
	_ openapi.GetVersionRequestObject,
) (openapi.GetVersionResponseObject, error) {
	return openapi.GetVersion200JSONResponse{
		Version: version.GetVersion(),
	}, nil
}
