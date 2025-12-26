package awesome

import (
	awesomepb "github.com/vendor116/awesome/pkg/protobuf/awesome"
)

var _ awesomepb.AwesomeServer = (*Server)(nil)

type Server struct {
	awesomepb.UnimplementedAwesomeServer
}

func NewServer() *Server {
	return &Server{}
}
