package models

import (
	"encoding/json"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/util"
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

func (a *Author) ShortAuthor() *ShortAuthor {
	sa := &ShortAuthor{
		ID:   a.ID,
		Name: a.Name,
	}
	if ident, err := a.DecodeIdentity(); err == nil {
		sa.IdentityIconImg = ident.Icon
	}
	return sa
}

func (a *Author) DecodeIdentity() (*oauth.Identity, error) {
	ident := &oauth.Identity{}
	if err := json.Unmarshal([]byte(a.Identity), ident); err != nil {
		return nil, err
	}
	return ident, nil
}

type ShortAuthor struct {
	ID              string `db:"id"`
	Name            string `db:"name"`
	IdentityIconImg string `db:"-"`
}

func (a *ShortAuthor) Proto() *api.Author {
	if a == nil {
		return nil
	}
	return &api.Author{
		Id:              a.ID,
		Name:            a.Name,
		IdentityIconImg: a.IdentityIconImg,
	}
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
	Author                *ShortAuthor
	Approver              bool
	AcceptedContributions int32
	AwardValue            float32
}

func (l *AuthorRanking) Proto() *api.AuthorRanking {
	if l == nil {
		return nil
	}
	return &api.AuthorRanking{
		Author:                l.Author.Proto(),
		Approver:              l.Approver,
		AcceptedContributions: l.AcceptedContributions,
		AwardValue:            l.AwardValue,
	}
}

type AuthorReward struct {
	ID       string `db:"id"`
	AuthorID string `db:"author_id"`

	// Threshold is not the number of contributions - it's the number of multiples of RewardSpacing
	// i.e. reward is issued at 5/10/15 which will be thresholds 1/2/3.
	Threshold int32 `db:"threshold"`

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

func (a *AuthorReward) ClaimedProto() *api.ClaimedReward {
	return &api.ClaimedReward{
		Id:               a.ID,
		ClaimKind:        util.PString(a.ClaimKind),
		ClaimValue:       util.PFloat32(a.ClaimValue),
		ClaimCurrency:    util.PString(a.ClaimValueCurrency),
		ClaimDescription: util.PString(a.ClaimDescription),
		ClaimAt:          util.PTime(a.ClaimAt).Format(time.RFC3339),
	}
}
