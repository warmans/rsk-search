package models

import "time"

type Author struct {
	ID        string
	Name      string
	Identity  string
	CreatedAt time.Time
	Banned    bool
	Approver  bool
}

type AuthorStats struct {
	ContributionsInLastHour int32
	PendingContributions    int32
	ApprovedContributions   int32
	RejectedContributions   int32
}
