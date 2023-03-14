package grpc

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func NewStatus(grpcStatus *status.Status) *Status {
	return &Status{Sta: grpcStatus}
}

type Status struct {
	Sta *status.Status
}

func (sta *Status) Error() string {
	return sta.Sta.Err().Error()
}

func ErrInvalidRequestField(field string, err error, moreDetails ...string) *Status {

	details := []proto.Message{
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{Field: field, Description: err.Error()},
			},
		},
		&errdetails.DebugInfo{
			Detail:       err.Error(),
			StackEntries: util.ErrTrace(err, 5),
		},
	}
	for _, v := range moreDetails {
		details = append(
			details,
			&errdetails.DebugInfo{
				Detail: v,
			},
		)
	}

	s, err := status.New(codes.InvalidArgument, http.StatusText(http.StatusBadRequest)).WithDetails(details...)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrFromStore(err error, id string) *Status {
	if staErr, ok := status.FromError(err); ok {
		return NewStatus(staErr)
	}
	if err == sql.ErrNoRows {
		return ErrNotFound(id)
	}
	if err == rw.ErrNotPermitted {
		return ErrPermissionDenied(err.Error())
	}
	return ErrInternal(err)
}

func ErrInternal(err error) *Status {
	if err == nil {
		return NewStatus(status.New(codes.Internal, http.StatusText(http.StatusInternalServerError)))
	}
	// do not wrap existing grpc errors
	if sta, ok := status.FromError(err); ok {
		return NewStatus(sta)
	}
	s, errErr := status.New(codes.Internal, http.StatusText(http.StatusInternalServerError)).WithDetails(
		&errdetails.DebugInfo{
			Detail:       err.Error(),
			StackEntries: util.ErrTrace(err, 5),
		},
		&errdetails.ErrorInfo{Reason: err.Error()},
	)
	if errErr != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrNotFound(id string) *Status {
	s, err := status.New(codes.NotFound, http.StatusText(http.StatusNotFound)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("no record found with ID: %s", id),
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrServerConfused() *Status {
	s, err := status.New(codes.InvalidArgument, http.StatusText(http.StatusBadRequest)).WithDetails(
		&errdetails.DebugInfo{
			Detail: "Server was confused by request",
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrAuthFailed() *Status {
	s, err := status.New(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized)).WithDetails(
		&errdetails.DebugInfo{
			Detail: "Authorization request failed as the verification code did not match. It may have already expired.",
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrUnauthorized(reason string) *Status {
	s, err := status.New(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("Authorization failed with reason: %s", reason),
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrPermissionDenied(reason string) *Status {
	s, err := status.New(codes.PermissionDenied, http.StatusText(http.StatusForbidden)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("Permission was denied for action: %s", reason),
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrFailedPrecondition(reason string) *Status {
	s, err := status.New(codes.FailedPrecondition, http.StatusText(http.StatusBadRequest)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("Precondition failed: %s", reason),
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrNotImplemented() *Status {
	return NewStatus(status.New(codes.Unimplemented, http.StatusText(http.StatusNotImplemented)))
}

func ErrRateLimited() *Status {
	s, err := status.New(codes.Unavailable, http.StatusText(http.StatusServiceUnavailable)).WithDetails(
		&errdetails.DebugInfo{
			Detail: "Too many requests",
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}

func ErrThirdParty(reason string) *Status {
	s, err := status.New(codes.Unavailable, http.StatusText(http.StatusServiceUnavailable)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("External service was unable to process request: %s", reason),
		},
	)
	if err != nil {
		return NewStatus(status.New(codes.Internal, "failed to create error"))
	}
	return NewStatus(s)
}
