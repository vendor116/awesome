package v1

import (
	v1 "github.com/vendor116/awesome/pkg/rest/v1"
)

var _ v1.StrictServerInterface = (*Server)(nil)

type Server struct {
	VersionHandlers
}

func NewServer() *Server {
	return &Server{
		VersionHandlers: NewVersionHandler(),
	}
}
