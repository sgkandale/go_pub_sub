package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	ErrInvalidRequest = status.Errorf(codes.FailedPrecondition, "invalid request")
	ErrNoMsgBody      = status.Errorf(codes.InvalidArgument, "no message body")
	ErrNoTopic        = status.Errorf(codes.InvalidArgument, "no topic specified")
)
