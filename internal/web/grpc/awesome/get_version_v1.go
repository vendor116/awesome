package awesome

import (
	"context"

	awesomepb "github.com/vendor116/awesome/pkg/protobuf/awesome"
)

func (s *Server) GetVersionV1(
	_ context.Context,
	_ *awesomepb.GetVersionV1Request,
) (*awesomepb.GetVersionV1Response, error) {
	return &awesomepb.GetVersionV1Response{
		Version: "dev",
	}, nil
}
