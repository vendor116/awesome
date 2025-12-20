package v1

import (
	"context"

	v1 "github.com/vendor116/awesome/pkg/rest/v1"
)

type VersionHandlers struct{}

func NewVersionHandler() VersionHandlers {
	return VersionHandlers{}
}

func (vh VersionHandlers) GetVersion(
	_ context.Context,
	_ v1.GetVersionRequestObject,
) (v1.GetVersionResponseObject, error) {
	return v1.GetVersion200JSONResponse{
		Version: "dev",
	}, nil
}
