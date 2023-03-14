package grpc

import (
	"context"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/store/common"
	"strings"
)

func NewQueryModifiers(req interface{}) (*common.QueryModifier, error) {
	q := common.Q()
	if p, ok := req.(common.Pager); ok {
		q.Apply(common.WithPaging(p.GetPageSize(), p.GetPage()))
	}
	if p, ok := req.(common.Sorter); ok {
		if p.GetSortField() != "" {
			givenDirection := common.SortDirection(strings.ToUpper(p.GetSortDirection()))
			if givenDirection != common.SortAsc && givenDirection != common.SortDesc {
				return nil, ErrInvalidRequestField("sort_direction", errors.New("Must be 'asc' or 'desc'"))
			}
			q.Apply(common.WithSorting(p.GetSortField(), givenDirection))
		}
	}
	if p, ok := req.(common.Filterer); ok {
		if strings.TrimSpace(p.GetFilter()) != "" {
			fil, err := filter.Parse(p.GetFilter())
			if err != nil {
				return nil, ErrInvalidRequestField("filter", err)
			}
			q.Apply(common.WithFilter(fil))
		}
	}
	return q, nil
}

func GetClaims(ctx context.Context, auth *jwt.Auth) (*jwt.Claims, error) {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return nil, ErrUnauthorized("no token provided")
	}
	claims, err := auth.VerifyToken(token)
	if err != nil {
		return nil, ErrUnauthorized(err.Error())
	}
	return claims, nil
}

func IsAuthor(ctx context.Context, auth *jwt.Auth, authorID string) bool {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return false
	}
	if claims, err := auth.VerifyToken(token); err == nil {
		return claims.AuthorID == authorID
	}
	return false
}

func IsApprover(ctx context.Context, auth *jwt.Auth) bool {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return false
	}
	if claims, err := auth.VerifyToken(token); err == nil {
		return claims.Approver
	}
	return false
}
