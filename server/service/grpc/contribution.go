package grpc

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/coffee"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

func NewContributionsService(
	logger *zap.Logger,
	srvCfg config.SearchServiceConfig,
	persistentDB *rw.Conn,
	auth *jwt.Auth,
	pledgeClient *pledge.Client,
	coffee *coffee.Client,
) *ContributionsService {

	var rankCache models.Ranks
	err := persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		rankCache, err = s.ListRanks(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create rank cache: %s", err.Error()))
	}

	return &ContributionsService{
		logger:       logger,
		srvCfg:       srvCfg,
		persistentDB: persistentDB,
		auth:         auth,
		pledgeClient: pledgeClient,
		rankCache:    rankCache,
		coffee:       coffee,
	}
}

type ContributionsService struct {
	logger       *zap.Logger
	srvCfg       config.SearchServiceConfig
	persistentDB *rw.Conn
	auth         *jwt.Auth
	pledgeClient *pledge.Client
	rankCache    models.Ranks
	coffee       *coffee.Client
}

func (s *ContributionsService) RegisterGRPC(server *grpc.Server) {
	api.RegisterContributionsServiceServer(server, s)
}

func (s *ContributionsService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterContributionsServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *ContributionsService) ListAuthorRanks(ctx context.Context, request *api.ListAuthorRanksRequest) (*api.AuthorRankList, error) {

	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}
	qm.Apply(common.WithDefaultSorting("points", common.SortDesc))

	out := &api.AuthorRankList{Rankings: make([]*api.AuthorRank, 0)}
	err = s.persistentDB.WithStore(func(st *rw.Store) error {
		lb, err := st.ListAuthorRankings(ctx, qm)
		if err != nil {
			return err
		}
		for _, v := range lb {
			out.Rankings = append(out.Rankings, v.Proto(s.rankCache))
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return out, err
}

func (s *ContributionsService) ListPendingRewards(ctx context.Context, empty *emptypb.Empty) (*api.PendingRewardList, error) {

	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, "")
	}

	result := &api.PendingRewardList{
		Rewards: []*api.Reward{},
	}
	for _, v := range rewards {
		result.Rewards = append(result.Rewards, getRewardForThreshold(v))
	}

	return result, nil
}

func (s *ContributionsService) ListClaimedRewards(ctx context.Context, empty *emptypb.Empty) (*api.ClaimedRewardList, error) {

	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, "")
	}

	result := &api.ClaimedRewardList{
		Rewards: []*api.ClaimedReward{},
	}
	for _, v := range rewards {
		result.Rewards = append(result.Rewards, v.ClaimedProto())
	}
	return result, nil
}

func (s *ContributionsService) ClaimReward(ctx context.Context, request *api.ClaimRewardRequest) (*emptypb.Empty, error) {

	if s.srvCfg.RewardsDisabled {
		return nil, ErrFailedPrecondition("rewards are disabled temporarily")
	}

	err := s.persistentDB.WithStore(func(store *rw.Store) error {

		reward, err := store.GetRewardForUpdate(ctx, request.Id)
		if err != nil {
			return err
		}

		donationArgs := request.GetDonationArgs()
		if donationArgs == nil {
			return ErrInvalidRequestField("args", errors.New("exepcted donation details in args"))
		}

		var recipient *api.DonationRecipient
		for _, v := range getDonationRecipients() {
			if v.Id == request.GetDonationArgs().Recipient {
				recipient = v
			}
		}
		if recipient == nil {
			return ErrInvalidRequestField("args", errors.New("unknown recipient"))
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
			return ErrThirdParty("donation could not be completed")
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
			recipient.Name,
		)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id)
	}
	return &emptypb.Empty{}, nil
}

func (s *ContributionsService) ListDonationRecipients(ctx context.Context, request *api.ListDonationRecipientsRequest) (*api.DonationRecipientList, error) {
	res := &api.DonationRecipientList{
		Organizations: getDonationRecipients(),
	}
	return res, nil
}

func (s *ContributionsService) ListAuthorContributions(ctx context.Context, request *api.ListAuthorContributionsRequest) (*api.AuthorContributionList, error) {

	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}

	out := &api.AuthorContributionList{Contributions: make([]*api.AuthorContribution, 0)}

	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		cont, err := s.ListAuthorContributions(ctx, qm)
		if err != nil {
			return err
		}
		for _, v := range cont {
			out.Contributions = append(out.Contributions, v.Proto())
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return out, nil
}

func (s *ContributionsService) GetDonationStats(ctx context.Context, empty *emptypb.Empty) (*api.DonationStats, error) {
	var stats models.DonationRecipientStats
	if err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		stats, err = s.GetDonationStats(ctx)
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "")
	}
	return stats.Proto(), nil
}

func (s *ContributionsService) ListIncomingDonations(ctx context.Context, _ *api.ListIncomingDonationsRequest) (*api.IncomingDonationList, error) {
	out := &api.IncomingDonationList{}

	if s.coffee == nil {
		return out, nil
	}
	sups, err := s.coffee.Supporters()
	if err != nil {
		return nil, ErrInternal(err)
	}
	for _, sup := range sups.Data {
		priceFloat, err := strconv.ParseFloat(sup.Price, 32)
		if err != nil {
			s.logger.Error("failed to parse float", zap.String("val", sup.Price), zap.Error(err))
			continue
		}
		out.Donations = append(out.Donations, &api.IncomingDonation{
			Name:           sup.Name,
			Amount:         float32(priceFloat) * float32(sup.Qty),
			AmountCurrency: sup.Currency,
			Note:           sup.Note,
		})
	}
	return out, nil
}

func getDonationRecipients() []*api.DonationRecipient {
	return []*api.DonationRecipient{
		{
			Id:      "b5096382-996f-463d-a531-75b29163a2b2",
			Name:    "UNICEF - Ukraine Relief",
			Mission: "UNICEF has been working nonstop in eastern Ukraine, delivering lifesaving programs for affected children and families as fighting has taken an increasingly heavy toll on the civilian population of 3.4 million people — including 510,000 children — living in the Donbas region.",
			LogoUrl: "/assets/logo/bf20017fbe112e29.jpeg",
			NgoId:   "13-1760110-006",
			Url:     "https://www.unicef.org/ukraine/en",
		},
		{
			Id:      "d694d55b-5888-4368-9a88-70acd55f33b0",
			Name:    "Revived Soldiers Ukraine",
			Mission: "Aiding with medical help to people of Ukraine, severely wounded soldiers and members of their families.",
			LogoUrl: "/assets/logo/rsu.jpg",
			NgoId:   "47-5315018",
			Url:     "https://www.rsukraine.org/",
		},
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

}

func getRewardForThreshold(mod *models.AuthorReward) *api.Reward {
	return &api.Reward{
		Id:            mod.ID,
		Kind:          api.Reward_DONATION,
		Name:          "Here's some tat in a jiffy bag",
		Criteria:      fmt.Sprintf("Earn %0.2f Points", mod.PointsSpent),
		Value:         2,
		ValueCurrency: "USD",
	}
}
