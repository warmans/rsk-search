package models

import (
	"encoding/json"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/oauth"
	"time"
)

type Author struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Identity  string    `db:"identity"`
	CreatedAt time.Time `db:"created_at"`
	Banned    bool      `db:"banned"`
	Approver  bool      `db:"approver"`
}

func (a *Author) DecodeIdentity() (*oauth.Identity, error) {
	ident := &oauth.Identity{}
	if err := json.Unmarshal([]byte(a.Identity), ident); err != nil {
		return nil, err
	}
	return ident, nil
}

type AuthorStats struct {
	ContributionsInLastHour      int32
	PendingContributions         int32
	RequestApprovalContributions int32
	ApprovedContributions        int32
	RejectedContributions        int32
}

type AuthorLeaderboard struct {
	Authors []*AuthorRanking
}

func (l *AuthorLeaderboard) Proto() *api.AuthorLeaderboard {
	if l == nil {
		return nil
	}
	out := &api.AuthorLeaderboard{
		Authors: make([]*api.AuthorRanking, len(l.Authors)),
	}
	for k := range l.Authors {
		out.Authors[k] = l.Authors[k].Proto()
	}
	return out
}

type AuthorRanking struct {
	Name                  string
	AcceptedContributions int32
	Approver              bool
}

func (l *AuthorRanking) Proto() *api.AuthorRanking {
	if l == nil {
		return nil
	}
	return &api.AuthorRanking{
		AuthorName:            l.Name,
		AcceptedContributions: l.AcceptedContributions,
		Approver:              l.Approver,
	}
}

type AuthorReward struct {
	ID                    string     `db:"id"`
	AuthorID              string     `db:"author_id"`
	Threshold             int32      `db:"threshold"`
	CreatedAt             time.Time  `db:"created_at"`
	Claimed               bool       `db:"claimed"`
	ClaimKind             *string    `db:"claim_kind"`
	ClaimValue            *float32   `db:"claim_value"`
	ClaimValueCurrency    *string    `db:"claim_value_currency"`
	ClaimDescription      *string    `db:"claim_description"`
	ClaimAt               *time.Time `db:"claim_at"`
	ClaimConfirmationCode *string    `db:"claim_confirmation_code"`
	Error                 *string    `db:"error"`
}
