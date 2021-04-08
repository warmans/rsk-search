package grpc

import (
	"context"
	"fmt"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SearchService) ListPendingRewards(ctx context.Context, empty *emptypb.Empty) (*api.PendingRewardList, error) {
	var rewards []*models.AuthorReward
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		rewards, err = s.ListPendingRewards(ctx)
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

func (s *SearchService) ClaimReward(ctx context.Context, request *api.ClaimRewardRequest) (*emptypb.Empty, error) {

	var donation *pledge.Donation
	err := s.persistentDB.WithStore(func(store *rw.Store) error {
		reward, err := store.GetRewardForUpdate(ctx, request.Id)
		if err != nil {
			return err
		}

		donationArgs := request.GetDonationArgs()
		if donationArgs == nil {
			return ErrInvalidRequestField("args", "exepcted donation details in args").Err()
		}
		//todo: fetch donations and check metadata for ID

		rewardValue := getRewardForThreshold(reward)
		donation, err = s.pledgeClient.CreateAnonymousDonation(pledge.AnonymousDonationRequest{
			OrganizationID: request.GetDonationArgs().Recipient,
			Amount:         fmt.Sprintf("%0.2f", rewardValue.Value),
			Metadata:       reward.ID,
		})
		if err != nil {
			s.logger.Error("Failed to claim reward. Pledge call failed", zap.Error(err))
			if err := store.FailReward(ctx, reward.ID, err.Error()); err != nil {
				s.logger.Error("Failed to mark reward as failed")
			}
			return ErrThirdParty("donation could not be completed").Err()
		}
		return store.ClaimReward(ctx, reward.ID)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return &emptypb.Empty{}, nil
}

func getRewardForThreshold(mod *models.AuthorReward) *api.Reward {
	switch mod.Threshold {
	case 5:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Man alive!"),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         1,
			ValueCurrency: "USD",
		}
	case 10:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Are you trying to turn my children into Communist revolutionaries?"),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         1,
			ValueCurrency: "USD",
		}
	case 15:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("In my opinion bronze is slightly better than gold."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         2,
			ValueCurrency: "USD",
		}
	case 20:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("I can't even begin to explain it."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         2,
			ValueCurrency: "USD",
		}
	case 25:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("There is a machine that can give you a tattoo."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         3,
			ValueCurrency: "USD",
		}
	case 30:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Kate Bush is on the phone!"),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         3,
			ValueCurrency: "USD",
		}
	default:
		return &api.Reward{
			Id:            mod.ID,
			Kind:          api.Reward_DONATION,
			Name:          fmt.Sprintf("Infinity sorty of, sorts it out for you."),
			Criteria:      fmt.Sprintf("Contribute %d transcription chunks.", mod.Threshold),
			Value:         1,
			ValueCurrency: "USD",
		}
	}
}
