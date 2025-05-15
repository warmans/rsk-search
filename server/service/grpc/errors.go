package grpc

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

func ErrInvalidRequestField(field string, originalErr error, moreDetails ...string) error {

	// todo: moreDetails need to be appended dynamically but it's blocked by:
	// https://github.com/grpc/grpc-go/issues/6133

	if originalErr == nil {
		s, err := status.New(codes.InvalidArgument, http.StatusText(http.StatusBadRequest)).WithDetails(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{Field: field, Description: fmt.Sprintf("Invalid field: %s", strings.Join(moreDetails, ", "))},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create error")
		}
		return s.Err()
	}

	s, err := status.New(codes.InvalidArgument, http.StatusText(http.StatusBadRequest)).WithDetails(&errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{Field: field, Description: originalErr.Error()},
		},
	},
		&errdetails.DebugInfo{
			Detail:       originalErr.Error(),
			StackEntries: util.ErrTrace(originalErr, 7),
		})
	if err != nil {
		return fmt.Errorf("failed to create error")
	}
	return s.Err()
}

func ErrFromStore(err error, id string) error {
	if staErr, ok := status.FromError(err); ok {
		return staErr.Err()
	}
	if err == sql.ErrNoRows {
		return ErrNotFound(id)
	}
	if err == rw.ErrNotPermitted {
		return ErrPermissionDenied(err.Error())
	}
	return ErrInternal(err)
}

func ErrInternal(err error) error {
	if err == nil {
		return status.New(codes.Internal, http.StatusText(http.StatusInternalServerError)).Err()
	}
	err = errors.WithStack(err)

	// do not wrap existing grpc errors
	if sta, ok := status.FromError(err); ok {
		return sta.Err()
	}
	s, errErr := status.New(codes.Internal, http.StatusText(http.StatusInternalServerError)).WithDetails(
		&errdetails.DebugInfo{
			Detail:       err.Error(),
			StackEntries: util.ErrTrace(err, 5),
		},
		&errdetails.ErrorInfo{Reason: err.Error()},
	)
	if errErr != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrNotFound(id string) error {
	s, err := status.New(codes.NotFound, http.StatusText(http.StatusNotFound)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("no record found with ID: %s", id),
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrServerConfused() error {
	s, err := status.New(codes.InvalidArgument, http.StatusText(http.StatusBadRequest)).WithDetails(
		&errdetails.DebugInfo{
			Detail: "Server was confused by request",
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrAuthFailed() error {
	s, err := status.New(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized)).WithDetails(
		&errdetails.DebugInfo{
			Detail: "Authorization request failed as the verification code did not match. It may have already expired.",
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrUnauthorized(reason string) error {
	s, err := status.New(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("Authorization failed with reason: %s", reason),
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrPermissionDenied(reason string) error {
	s, err := status.New(codes.PermissionDenied, http.StatusText(http.StatusForbidden)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("Permission was denied for action: %s", reason),
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrFailedPrecondition(reason string) error {
	s, err := status.New(codes.FailedPrecondition, http.StatusText(http.StatusBadRequest)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("Precondition failed: %s", reason),
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrNotImplemented() error {
	return status.New(codes.Unimplemented, http.StatusText(http.StatusNotImplemented)).Err()
}

func ErrRateLimited() error {
	s, err := status.New(codes.Unavailable, http.StatusText(http.StatusServiceUnavailable)).WithDetails(
		&errdetails.DebugInfo{
			Detail: "Too many requests",
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}

func ErrThirdParty(reason string) error {
	s, err := status.New(codes.Unavailable, http.StatusText(http.StatusServiceUnavailable)).WithDetails(
		&errdetails.DebugInfo{
			Detail: fmt.Sprintf("External service was unable to process request: %s", reason),
		},
	)
	if err != nil {
		return status.New(codes.Internal, "failed to create error").Err()
	}
	return s.Err()
}
