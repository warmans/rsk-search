package models

import "github.com/warmans/rsk-search/gen/api"

type DonationRecipientStats []*DonationRecipientStat

func (drs DonationRecipientStats) Proto() *api.DonationStats {
	out := make([]*api.RecipientStats, len(drs))
	for k, v := range drs {
		out[k] = v.Proto()
	}
	return &api.DonationStats{Stats: out}
}

type DonationRecipientStat struct {
	RecipientName     string  `db:"recipient_name"`
	PointsSpent       float32 `db:"points_spent"`
	DonationAmountUSD float32 `db:"claim_value"`
}

func (dr *DonationRecipientStat) Proto() *api.RecipientStats {
	return &api.RecipientStats{
		DonationRecipient: dr.RecipientName,
		PointsSpent:       dr.PointsSpent,
		DonatedAmountUsd:  dr.DonationAmountUSD,
	}
}
