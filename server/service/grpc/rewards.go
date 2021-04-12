package grpc

import (
	"context"
	"fmt"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/reward"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SearchService) ListPendingRewards(ctx context.Context, empty *emptypb.Empty) (*api.PendingRewardList, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	if s.srvCfg.RewardsDisabled {
		return &api.PendingRewardList{
			Rewards: make([]*api.Reward, 0),
		}, nil
	}

	var rewards []*models.AuthorReward

	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		rewards, err = s.ListPendingRewards(ctx, claims.AuthorID)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}

	result := &api.PendingRewardList{
		Rewards: []*api.Reward{},
	}
	for _, v := range rewards {
		result.Rewards = append(result.Rewards, getRewardForThreshold(v))
	}

	return result, nil
}

func (s *SearchService) ListClaimedRewards(ctx context.Context, empty *emptypb.Empty) (*api.ClaimedRewardList, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var rewards []*models.AuthorReward
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		rewards, err = s.ListClaimedRewards(ctx, claims.AuthorID)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}

	result := &api.ClaimedRewardList{
		Rewards: []*api.ClaimedReward{},
	}
	for _, v := range rewards {
		result.Rewards = append(result.Rewards, v.ClaimedProto())
	}
	return result, nil
}

func (s *SearchService) ClaimReward(ctx context.Context, request *api.ClaimRewardRequest) (*emptypb.Empty, error) {

	if s.srvCfg.RewardsDisabled {
		return nil, ErrFailedPrecondition("rewards are disabled temporarily").Err()
	}

	err := s.persistentDB.WithStore(func(store *rw.Store) error {

		reward, err := store.GetRewardForUpdate(ctx, request.Id)
		if err != nil {
			return err
		}

		donationArgs := request.GetDonationArgs()
		if donationArgs == nil {
			return ErrInvalidRequestField("args", "exepcted donation details in args").Err()
		}

		var recipient *api.DonationRecipient
		for _, v := range getDonationRecipients(reward.Threshold) {
			if v.Id == request.GetDonationArgs().Recipient {
				recipient = v
			}
		}
		if recipient == nil {
			return ErrInvalidRequestField("args", "unknown recipient").Err()
		}

		//todo: fetch donations and check metadata for ID

		rewardValue := getRewardForThreshold(reward)

		s.logger.Info(
			"creating donation",
			zap.String("reward_id", request.Id),
			zap.String("cause", recipient.Name),
			zap.String("cause_id", recipient.Id),
			zap.Float32("value", rewardValue.Value),
		)
		donation, err := s.pledgeClient.CreateAnonymousDonation(pledge.AnonymousDonationRequest{
			OrganizationID: recipient.Id,
			Amount:         fmt.Sprintf("%0.2f", rewardValue.Value),
			Metadata:       reward.ID,
		})
		if err != nil {
			s.logger.Error("Failed to claim reward. Pledge call failed", zap.Error(err))
			return ErrThirdParty("donation could not be completed").Err()
		}
		s.logger.Info(
			"donation OK",
			zap.String("id", request.Id),
			zap.String("cause", recipient.Name),
			zap.Float32("value", rewardValue.Value),
			zap.String("donation_id", donation.ID),
			zap.String("donation_status", donation.Status),
		)
		return store.ClaimReward(
			ctx,
			reward.ID,
			rewardValue.Kind.String(),
			rewardValue.Value,
			rewardValue.ValueCurrency,
			donation.ID,
			fmt.Sprintf("Donated %0.2f %s to %s", rewardValue.Value, rewardValue.ValueCurrency, recipient.Name),
		)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *SearchService) ListDonationRecipients(ctx context.Context, request *api.ListDonationRecipientsRequest) (*api.DonationRecipientList, error) {

	var reward *models.AuthorReward
	err := s.persistentDB.WithStore(func(store *rw.Store) error {
		var err error
		reward, err = store.GetRewardForUpdate(ctx, request.RewardId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.RewardId).Err()
	}

	//todo: vary results based on threshold

	res := &api.DonationRecipientList{
		Organizations: getDonationRecipients(reward.Threshold),
	}
	return res, nil
}

func getDonationRecipients(thresold int32) []*api.DonationRecipient {

	switch thresold {
	case 2:
		return []*api.DonationRecipient{
			{
				Id:      "5957dbb1-b979-4b33-b068-ad56aadbe3f8",
				Name:    "St. John's Ambulance",
				Mission: "We are the charity that steps forward in the moments that matter, to save lives and keep communities safe.",
				LogoUrl: "/assets/logo/43-1634280-0504257.png",
				NgoId:   "43-1634280-0504257",
				Url:     "https://www.sja.org.uk/",
				Quote:   "But seriously, All joking aside. I genuinely wanted to give some massive props - give some big-ups - to the St. John's people, because I genuinely, without any joking, and I  genuinely think they they do a brilliant job.",
			},
			{
				Id:      "11034875-b8d5-4653-8558-214ae12a81b7",
				Name:    "Dogs Trust",
				Mission: "Our mission is to bring about the day when all dogs can enjoy a happy life, free from the threat of unnecessary destruction.",
				LogoUrl: "/assets/logo/43-1634280-0279288.jpg",
				NgoId:   "43-1634280-0279288",
				Url:     "https://www.dogstrust.org.uk",
			},
			{
				Id:      "40ebb87d-62f4-4297-a808-c5f35ef3719f",
				Name:    "Rainforest Alliance",
				Mission: "The Rainforest Alliance works to conserve biodiversity and ensure sustainable livelihoods by transforming land-use practices, business practices and consumer behavior.\n\nWe envision a world where people can thrive and prosper in harmony with the land",
				LogoUrl: "/assets/logo/13-3377893.png",
				NgoId:   "13-3377893",
				Url:     "https://www.pledge.to/organizations/13-3377893/rainforest-alliance",
			},
		}
	default:
		return []*api.DonationRecipient{
			{
				Id:      "e349c52c-73aa-4123-83b2-6466d1aa2d54",
				Name:    "International Primate Protection League",
				Mission: "PPL is a grassroots nonprofit organization dedicated to protecting the world’s remaining primates, great and small. Since 1973 we have worked to expose primate abuse and battled international traffickers.",
				LogoUrl: "/assets/logo/51-0194013.png",
				NgoId:   "51-0194013",
				Url:     "https://www.pledge.to/organizations/51-0194013/international-primate-protection-league",
			},
			{
				Id:      "700f6e06-a00d-46fe-a76a-e8271585c2bb",
				Name:    "World Wildlife Fund",
				Mission: "As the world’s leading conservation organization, WWF works in nearly 100 countries. At every level, we collaborate with people around the world to develop and deliver innovative solutions that protect communities, wildlife, and the places in which they live.",
				LogoUrl: "/assets/logo/52-1693387.png",
				NgoId:   "52-1693387",
				Url:     "https://www.pledge.to/organizations/52-1693387/world-wildlife-fund",
			},
			{
				Id:      "27547c25-7b00-4cb1-9c21-2834acb37da3",
				Name:    "Rainforest Rescue",
				Mission: "Rainforest Rescue is a not-for-profit organisation that has been protecting and restoring rainforests in Australia and internationally since 1998 by providing opportunities for individuals and businesses to Protect Rainforests Forever.",
				LogoUrl: "/assets/logo/30-0108263-675.svg",
				NgoId:   "30-0108263-675",
				Url:     "https://www.pledge.to/organizations/30-0108263-675/rainforest-rescue",
			},
		}
	}
}

func getRewardForThreshold(mod *models.AuthorReward) *api.Reward {
	switch mod.Threshold {
	case 1:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Man alive!"),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         1,
			ValueCurrency: "USD",
		}
	case 2:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Are you trying to turn my children into Communist revolutionaries?"),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         1,
			ValueCurrency: "USD",
		}
	case 3:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("In my opinion bronze is slightly better than gold."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         2,
			ValueCurrency: "USD",
		}
	case 4:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("I can't even begin to explain it."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         2,
			ValueCurrency: "USD",
		}
	case 5:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("There is a machine that can give you a tattoo."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         3,
			ValueCurrency: "USD",
		}
	case 6:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Kate Bush is on the phone!"),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         3,
			ValueCurrency: "USD",
		}
	default:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Infinity sorty of, sorts it out for you."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold*reward.RewardSpacing),
			Value:         1,
			ValueCurrency: "USD",
		}
	}
}
