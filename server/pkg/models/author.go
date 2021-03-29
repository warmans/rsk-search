package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"time"
)

type Author struct {
	ID        string
	Name      string
	Identity  string
	CreatedAt time.Time
	Banned    bool
	Approver  bool
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
