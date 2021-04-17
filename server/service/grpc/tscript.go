package grpc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/reward"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/tscript"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewTscriptService(
	logger *zap.Logger,
	srvCfg config.SearchServiceConfig,
	persistentDB *rw.Conn,
	auth *jwt.Auth,
	pledgeClient *pledge.Client,
) *TscriptService {
	return &TscriptService{
		logger:       logger,
		srvCfg:       srvCfg,
		persistentDB: persistentDB,
		auth:         auth,
		pledgeClient: pledgeClient,
	}
}

type TscriptService struct {
	logger       *zap.Logger
	srvCfg       config.SearchServiceConfig
	persistentDB *rw.Conn
	auth         *jwt.Auth
	pledgeClient *pledge.Client
}

func (s *TscriptService) RegisterGRPC(server *grpc.Server) {
	api.RegisterTscriptServiceServer(server, s)
}

func (s *TscriptService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterTscriptServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *TscriptService) ListTscripts(ctx context.Context, request *api.ListTscriptsRequest) (*api.TscriptList, error) {
	el := &api.TscriptList{
		Tscripts: []*api.TscriptStats{},
	}
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		eps, err := s.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, e := range eps {
			el.Tscripts = append(el.Tscripts, e.Proto())
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return el, nil
}

func (s *TscriptService) GetChunkStats(ctx context.Context, empty *emptypb.Empty) (*api.ChunkStats, error) {
	var stats *models.ChunkStats
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		stats, err = s.GetChunkStats(ctx)
		if stats == nil {
			stats = &models.ChunkStats{}
		}
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return stats.Proto(), nil
}

func (s *TscriptService) GetTscriptTimeline(ctx context.Context, request *api.GetTscriptTimelineRequest) (*api.TscriptTimeline, error) {
	result := &api.TscriptTimeline{
		Events: make([]*api.TscriptTimelineEvent, 0),
	}
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		events, err := s.ListTscriptTimelineEvents(ctx, request.TscriptId, int(request.Page))
		if err != nil {
			return err
		}
		for _, v := range events {
			result.Events = append(result.Events, v.Proto())
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return result, nil
}

func (s *TscriptService) GetChunk(ctx context.Context, request *api.GetChunkRequest) (*api.Chunk, error) {
	var chunk *models.Chunk
	var contributionCount int32

	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		chunk, err = s.GetChunk(ctx, request.Id)
		if err != nil {
			return err
		}
		if chunk == nil {
			return ErrNotFound(request.Id).Err()
		}
		contributionCount, err = s.GetChunkContributionCount(ctx, request.Id)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return chunk.Proto(contributionCount), nil
}

func (s *TscriptService) CreateChunkContribution(ctx context.Context, request *api.CreateChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		stats, err := s.GetAuthorStats(ctx, claims.AuthorID)
		if err != nil {
			return err
		}
		if stats.ContributionsInLastHour > 5 {
			return ErrRateLimited().Err()
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}

	lines, _, err := tscript.Import(bufio.NewScanner(bytes.NewBufferString(request.Transcript)), 0)
	if err != nil {
		return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
	}
	if len(lines) == 0 {
		return nil, ErrInvalidRequestField("transcript", "no valid lines parsed from transcript").Err()
	}

	var contrib *models.Contribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		contrib, err = s.CreateContribution(ctx, &models.ContributionCreate{
			AuthorID:      claims.AuthorID,
			ChunkID:       request.ChunkId,
			Transcription: request.Transcript,
			State:         models.ContributionStatePending,
		})
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return contrib.Proto(), nil
}

func (s *TscriptService) UpdateChunkContribution(ctx context.Context, request *api.UpdateChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.Contribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetContribution(ctx, request.ContributionId)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}

	// validate change is allowed
	if err := s.validateContributionStateUpdate(claims, contrib, request.State); err != nil {
		return nil, err
	}

	// allow invalid transcript while the contribution is still pending.
	if request.State != api.ContributionState_STATE_PENDING {
		lines, _, err := tscript.Import(bufio.NewScanner(bytes.NewBufferString(request.Transcript)), 0)
		if err != nil {
			return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
		}
		if len(lines) == 0 {
			return nil, ErrInvalidRequestField("transcript", "no valid lines parsed from transcript").Err()
		}
	}

	err = s.persistentDB.WithStore(func(tx *rw.Store) error {

		contrib.Transcription = request.Transcript
		contrib.State = models.ContributionStateFromProto(request.State)

		if err := s.createContributionActivity(tx, ctx, claims, contrib, ""); err != nil {
			return err
		}
		return tx.UpdateContribution(ctx, &models.ContributionUpdate{
			ID:            contrib.ID,
			AuthorID:      contrib.Author.ID,
			Transcription: contrib.Transcription,
			State:         contrib.State,
		})
	})
	if err != nil {
		return nil, ErrFromStore(err, contrib.ID).Err()
	}

	return contrib.Proto(), nil
}

func (s *TscriptService) RequestChunkContributionState(ctx context.Context, request *api.RequestChunkContributionStateRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.Contribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetContribution(ctx, request.ContributionId)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	if err := s.validateContributionStateUpdate(claims, contrib, request.RequestState); err != nil {
		return nil, err
	}
	if request.Comment != "" && claims.Approver {
		return nil, ErrPermissionDenied("Only an approver can set a state comment.").Err()
	}
	err = s.persistentDB.WithStore(func(tx *rw.Store) error {

		contrib.State = models.ContributionStateFromProto(request.RequestState)
		contrib.StateComment = request.Comment

		if err := s.createContributionActivity(tx, ctx, claims, contrib, contrib.StateComment); err != nil {
			return err
		}
		return tx.UpdateContributionState(ctx, contrib.ID, contrib.State, contrib.StateComment)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	return contrib.Proto(), nil
}

func (s *TscriptService) DiscardDraftContribution(ctx context.Context, request *api.DiscardDraftContributionRequest) (*emptypb.Empty, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.Contribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetContribution(ctx, request.ContributionId)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	if claims.AuthorID != contrib.Author.ID {
		return nil, ErrPermissionDenied("you are not the author of this contribution").Err()
	}
	if contrib.State != models.ContributionStatePending {
		return nil, ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be delete. Actual state was: %s", contrib.State)).Err()
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.DeleteContribution(ctx, request.ContributionId)
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *TscriptService) ListContributions(ctx context.Context, request *api.ListContributionsRequest) (*api.ContributionList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}
	out := &api.ContributionList{
		Contributions: make([]*api.Contribution, 0),
	}
	if err := s.persistentDB.WithStore(func(store *rw.Store) error {
		contributions, err := store.ListContributions(ctx, qm)
		for _, v := range contributions {
			out.Contributions = append(out.Contributions, v.TscriptContributionProto())
		}
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return out, nil
}

func (s *TscriptService) GetAuthorLeaderboard(ctx context.Context, empty *emptypb.Empty) (*api.AuthorLeaderboard, error) {
	var out *api.AuthorLeaderboard
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		lb, err := s.AuthorLeaderboard(ctx)
		if err != nil {
			return err
		}
		out = lb.Proto()
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return out, err
}

func (s *TscriptService) GetChunkContribution(ctx context.Context, request *api.GetChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.Contribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetContribution(ctx, request.ContributionId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	if claims.Approver == false {
		if contrib.State == models.ContributionStatePending && contrib.Author.ID != claims.AuthorID {
			return nil, ErrPermissionDenied("you cannot view another author's contribution when it is in the pending state").Err()
		}
	}
	return contrib.Proto(), nil
}

func (s *TscriptService) ListPendingRewards(ctx context.Context, empty *emptypb.Empty) (*api.PendingRewardList, error) {

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

func (s *TscriptService) ListClaimedRewards(ctx context.Context, empty *emptypb.Empty) (*api.ClaimedRewardList, error) {

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

func (s *TscriptService) ClaimReward(ctx context.Context, request *api.ClaimRewardRequest) (*emptypb.Empty, error) {

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

func (s *TscriptService) ListDonationRecipients(ctx context.Context, request *api.ListDonationRecipientsRequest) (*api.DonationRecipientList, error) {

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

func (s *TscriptService) getClaims(ctx context.Context) (*jwt.Claims, error) {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return nil, ErrUnauthorized("no token provided").Err()
	}
	claims, err := s.auth.VerifyToken(token)
	if err != nil {
		return nil, ErrUnauthorized(err.Error()).Err()
	}
	return claims, nil
}

func (s *TscriptService) validateContributionStateUpdate(claims *jwt.Claims, currentState *models.Contribution, requestedState api.ContributionState) error {
	if !claims.Approver {
		if currentState.Author.ID != claims.AuthorID {
			return ErrPermissionDenied("you are not the author of this contribution").Err()
		}
		if requestedState == api.ContributionState_STATE_APPROVED || requestedState == api.ContributionState_STATE_REJECTED {
			return ErrPermissionDenied("you are not an approver").Err()
		}
	}
	// if the contribution has been rejected allow the author to return it to pending.
	if currentState.State == models.ContributionStateRejected {
		if requestedState != api.ContributionState_STATE_PENDING {
			return ErrFailedPrecondition(fmt.Sprintf("Only rejected contributions can be reverted to pending. Actual state was: %s (requested: %s)", currentState.State, requestedState)).Err()
		}
	} else {
		/// otherwise only allow it to be updated if it's in the pending or approval requested state.
		if currentState.State != models.ContributionStatePending && currentState.State != models.ContributionStateApprovalRequested {
			return ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be edited. Actual state was: %s", currentState.State)).Err()
		}
	}
	return nil
}

func (s *TscriptService) createContributionActivity(tx *rw.Store, ctx context.Context, claims *jwt.Claims, contrib *models.Contribution, comment string) error {
	suffix := "."
	if comment != "" {
		suffix = fmt.Sprintf(" with comment '%s'.", comment)
	}
	switch contrib.State {
	case models.ContributionStateApprovalRequested:
		if err := tx.CreateTscriptTimelineEvent(ctx, contrib.ChunkID, claims.Identity.Name, fmt.Sprintf("Submitted contribution %s for approval%s", contrib.ID, suffix)); err != nil {
			return err
		}
	case models.ContributionStateApproved:
		if err := tx.CreateTscriptTimelineEvent(ctx, contrib.ChunkID, claims.Identity.Name, fmt.Sprintf("Approved contribution %s%s", contrib.ID, suffix)); err != nil {
			return err
		}
	case models.ContributionStateRejected:
		if err := tx.CreateTscriptTimelineEvent(ctx, contrib.ChunkID, claims.Identity.Name, fmt.Sprintf("Rejected contribution %s%s", contrib.ID, suffix)); err != nil {
			return err
		}
	}
	return nil
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
