package grpc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/transcript"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewContribService(
	logger *zap.Logger,
	srvCfg config.SearchServiceConfig,
	persistentDB *rw.Conn,
	auth *jwt.Auth,
	pledgeClient *pledge.Client,
	episodeCache *data.EpisodeCache,
) *ContribService {

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

	return &ContribService{
		logger:       logger,
		srvCfg:       srvCfg,
		persistentDB: persistentDB,
		auth:         auth,
		pledgeClient: pledgeClient,
		episodeCache: episodeCache,
		rankCache:    rankCache,
	}
}

type ContribService struct {
	logger       *zap.Logger
	srvCfg       config.SearchServiceConfig
	persistentDB *rw.Conn
	auth         *jwt.Auth
	pledgeClient *pledge.Client
	episodeCache *data.EpisodeCache
	rankCache    models.Ranks
}

func (s *ContribService) RegisterGRPC(server *grpc.Server) {
	api.RegisterContribServiceServer(server, s)
}

func (s *ContribService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterContribServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *ContribService) ListTscripts(ctx context.Context, request *api.ListTscriptsRequest) (*api.TscriptList, error) {
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

func (s *ContribService) GetChunkStats(ctx context.Context, _ *emptypb.Empty) (*api.ChunkStats, error) {
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

func (s *ContribService) GetChunk(ctx context.Context, request *api.GetChunkRequest) (*api.Chunk, error) {
	var chunk *models.Chunk
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		chunk, err = s.GetChunk(ctx, request.Id)
		if err != nil {
			return err
		}
		if chunk == nil {
			return ErrNotFound(request.Id).Err()
		}
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return chunk.Proto(), nil
}

func (s *ContribService) ListChunks(ctx context.Context, request *api.ListChunksRequest) (*api.ChunkList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}
	if qm.Filter != nil {
		qm.Filter = filter.And(filter.Eq("tscript_id", filter.String(request.TscriptId)), qm.Filter)
	} else {
		qm.Filter = filter.Eq("tscript_id", filter.String(request.TscriptId))
	}
	if qm.Sorting == nil {
		qm.Sorting = &common.Sorting{Field: "start_second", Direction: common.SortAsc}
	}

	out := &api.ChunkList{
		Chunks: make([]*api.Chunk, 0),
	}
	if err := s.persistentDB.WithStore(func(store *rw.Store) error {
		chunks, err := store.ListChunks(ctx, qm)
		for _, v := range chunks {
			out.Chunks = append(out.Chunks, v.Proto())
		}
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return out, nil
}

func (s *ContribService) CreateChunkContribution(ctx context.Context, request *api.CreateChunkContributionRequest) (*api.ChunkContribution, error) {

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

	if err := transcript.Validate(bufio.NewScanner(bytes.NewBufferString(request.Transcript))); err != nil {
		return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
	}

	var contrib *models.ChunkContribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		contrib, err = s.CreateChunkContribution(ctx, &models.ContributionCreate{
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

func (s *ContribService) UpdateChunkContribution(ctx context.Context, request *api.UpdateChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.ChunkContribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetChunkContribution(ctx, request.ContributionId)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}

	// validate change is allowed
	if err := s.validateContributionStateUpdate(claims, contrib.Author.ID, contrib.State, request.State); err != nil {
		return nil, err
	}

	// allow invalid transcript while the contribution is still pending.
	if request.State != api.ContributionState_STATE_PENDING {
		if err := transcript.Validate(bufio.NewScanner(bytes.NewBufferString(request.Transcript))); err != nil {
			return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
		}
	}

	err = s.persistentDB.WithStore(func(tx *rw.Store) error {

		contrib.Transcription = request.Transcript
		contrib.State = models.ContributionStateFromProto(request.State)

		if err := tx.UpdateChunkContribution(ctx, &models.ContributionUpdate{
			ID:            contrib.ID,
			AuthorID:      contrib.Author.ID,
			Transcription: contrib.Transcription,
			State:         contrib.State,
		}); err != nil {
			return err
		}
		return tx.UpdateChunkActivity(ctx, contrib.ChunkID, rw.ActivityFromState(contrib.State))
	})
	if err != nil {
		return nil, ErrFromStore(err, contrib.ID).Err()
	}

	return contrib.Proto(), nil
}

func (s *ContribService) RequesChunktContributionState(ctx context.Context, request *api.RequestChunkContributionStateRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.ChunkContribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetChunkContribution(ctx, request.ContributionId)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	if err := s.validateContributionStateUpdate(claims, contrib.Author.ID, contrib.State, request.RequestState); err != nil {
		return nil, err
	}
	if request.Comment != "" && !claims.Approver {
		return nil, ErrPermissionDenied("Only an approver can set a state comment.").Err()
	}
	err = s.persistentDB.WithStore(func(tx *rw.Store) error {

		contrib.State = models.ContributionStateFromProto(request.RequestState)
		contrib.StateComment = request.Comment

		if err := tx.UpdateChunkContributionState(ctx, contrib.ID, contrib.State, contrib.StateComment); err != nil {
			return err
		}
		return tx.UpdateChunkActivity(ctx, contrib.ChunkID, rw.ActivityFromState(contrib.State))
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	return contrib.Proto(), nil
}

func (s *ContribService) DeleteChunkContribution(ctx context.Context, request *api.DeleteChunkContributionRequest) (*emptypb.Empty, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.ChunkContribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetChunkContribution(ctx, request.ContributionId)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	if claims.AuthorID != contrib.Author.ID {
		return nil, ErrPermissionDenied("you are not the author of this contribution").Err()
	}
	if contrib.State != models.ContributionStatePending {
		return nil, ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be deleted. Actual state was: %s", contrib.State)).Err()
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.DeleteContribution(ctx, request.ContributionId)
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *ContribService) ListChunkContributions(ctx context.Context, request *api.ListChunkContributionsRequest) (*api.ChunkContributionList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}

	if qm.Filter != nil {
		qm.Filter = filter.And(filter.Neq("state", filter.String("pending")), qm.Filter)
	} else {
		qm.Filter = filter.Neq("state", filter.String("pending"))
	}

	out := &api.ChunkContributionList{
		Contributions: make([]*api.ChunkContribution, 0),
	}
	if err := s.persistentDB.WithStore(func(store *rw.Store) error {
		contributions, err := store.ListChunkContributions(ctx, qm)
		for _, v := range contributions {
			out.Contributions = append(out.Contributions, v.Proto())
		}
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return out, nil
}

// GetAuthorLeaderboard deprecated
func (s *ContribService) GetAuthorLeaderboard(ctx context.Context, empty *emptypb.Empty) (*api.AuthorLeaderboard, error) {
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

func (s *ContribService) ListAuthorRanks(ctx context.Context, request *api.ListAuthorRanksRequest) (*api.AuthorRankList, error) {

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
		return nil, ErrFromStore(err, "").Err()
	}
	return out, err
}

func (s *ContribService) GetChunkContribution(ctx context.Context, request *api.GetChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var contrib *models.ChunkContribution
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetChunkContribution(ctx, request.ContributionId)
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

func (s *ContribService) ListPendingRewards(ctx context.Context, empty *emptypb.Empty) (*api.PendingRewardList, error) {

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

func (s *ContribService) ListClaimedRewards(ctx context.Context, empty *emptypb.Empty) (*api.ClaimedRewardList, error) {

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

func (s *ContribService) ClaimReward(ctx context.Context, request *api.ClaimRewardRequest) (*emptypb.Empty, error) {

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
		for _, v := range getDonationRecipients() {
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

func (s *ContribService) ListDonationRecipients(ctx context.Context, request *api.ListDonationRecipientsRequest) (*api.DonationRecipientList, error) {
	res := &api.DonationRecipientList{
		Organizations: getDonationRecipients(),
	}
	return res, nil
}

func (s *ContribService) ListTranscriptChanges(ctx context.Context, request *api.ListTranscriptChangesRequest) (*api.TranscriptChangeList, error) {

	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}

	out := &api.TranscriptChangeList{
		Changes: make([]*api.ShortTranscriptChange, 0),
	}
	if err := s.persistentDB.WithStore(func(store *rw.Store) error {
		contributions, err := store.ListTranscriptChanges(ctx, qm)
		for _, v := range contributions {
			// discard any pending contributions that are not owned by the author
			if v.State == models.ContributionStatePending && !s.isAuthor(ctx, v.Author.ID) {
				continue
			}
			out.Changes = append(out.Changes, v.ShortProto())
		}
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return out, nil
}

func (s *ContribService) GetTranscriptChange(ctx context.Context, request *api.GetTranscriptChangeRequest) (*api.TranscriptChange, error) {

	var change *models.TranscriptChange
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		change, err = s.GetTranscriptChange(ctx, request.Id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	if !s.isApprover(ctx) {
		if change.State == models.ContributionStatePending && !s.isAuthor(ctx, change.Author.ID) {
			return nil, ErrPermissionDenied("you cannot view another author's contribution when it is in the pending state").Err()
		}
	}

	return change.Proto(), nil
}

func (s *ContribService) GetTranscriptChangeDiff(ctx context.Context, request *api.GetTranscriptChangeDiffRequest) (*api.TranscriptChangeDiff, error) {

	var newTranscript *models.TranscriptChange
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		newTranscript, err = s.GetTranscriptChange(ctx, request.Id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	if !s.isApprover(ctx) {
		if newTranscript.State == models.ContributionStatePending && s.isAuthor(ctx, newTranscript.Author.ID) {
			return nil, ErrPermissionDenied("you cannot view another author's contribution diff when it is in the pending state").Err()
		}
	}

	oldTranscript, err := s.episodeCache.GetEpisode(newTranscript.EpID)
	if err != nil {
		return nil, ErrNotFound(newTranscript.EpID).Err()
	}

	oldRaw, err := transcript.Export(oldTranscript.Transcript, oldTranscript.Synopsis, oldTranscript.Trivia)
	if err != nil {
		return nil, err
	}

	edits := myers.ComputeEdits(span.URIFromPath(fmt.Sprintf("%s.txt", oldTranscript.ID())), oldRaw, newTranscript.Transcription)
	diff := fmt.Sprint(gotextdiff.ToUnified(fmt.Sprintf("%s.txt", oldTranscript.ID()), fmt.Sprintf("%s.txt", oldTranscript.ID()), oldRaw, edits))

	return &api.TranscriptChangeDiff{Diff: diff}, nil
}

func (s *ContribService) CreateTranscriptChange(ctx context.Context, request *api.CreateTranscriptChangeRequest) (*api.TranscriptChange, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	var change *models.TranscriptChange
	err = s.persistentDB.WithStore(func(s *rw.Store) error {

		// stop more than 1 pending change existing at once. Once the change is merged it can be ignored.
		changes, err := s.ListTranscriptChanges(
			ctx,
			common.Q(
				common.WithFilter(
					filter.And(
						filter.Eq("epid", filter.String(request.Epid)),
						filter.Eq("merged", filter.Bool(false)),
						filter.Neq("state", filter.String(string(models.ContributionStatePending))),
						filter.Neq("state", filter.String(string(models.ContributionStateRejected))),
					),
				),
			),
		)
		if err != nil {
			return err
		}
		if len(changes) > 0 {
			return ErrFailedPrecondition("multiple changes cannot exist at once. Try again once the current change has been processed.").Err()
		}
		change, err = s.CreateTranscriptChange(ctx, &models.TranscriptChangeCreate{
			AuthorID:      claims.AuthorID,
			EpID:          request.Epid,
			Summary:       "",
			Transcription: request.Transcript,
		})
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return change.Proto(), nil
}

func (s *ContribService) UpdateTranscriptChange(ctx context.Context, request *api.UpdateTranscriptChangeRequest) (*api.TranscriptChange, error) {

	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}

	var oldChange *models.TranscriptChange
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		oldChange, err = s.GetTranscriptChange(ctx, request.Id)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}

	// validate change is allowed
	if err := s.validateContributionStateUpdate(claims, oldChange.Author.ID, oldChange.State, request.State); err != nil {
		return nil, err
	}

	// allow invalid transcript while the contribution is still pending.
	if request.State != api.ContributionState_STATE_PENDING {
		if err := transcript.Validate(bufio.NewScanner(bytes.NewBufferString(request.Transcript))); err != nil {
			return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
		}
	}

	var updatedChange *models.TranscriptChange
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		updatedChange, err = s.UpdateTranscriptChange(ctx, &models.TranscriptChangeUpdate{
			ID:            request.Id,
			Summary:       "",
			Transcription: request.Transcript,
			State:         models.ContributionStateFromProto(request.State),
		}, request.PointsOnApprove)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return updatedChange.Proto(), nil
}

func (s *ContribService) DeleteTranscriptChange(ctx context.Context, request *api.DeleteTranscriptChangeRequest) (*emptypb.Empty, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	var change *models.TranscriptChange
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		change, err = s.GetTranscriptChange(ctx, request.Id)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	if change.Author.ID != claims.AuthorID {
		return nil, ErrNotFound(request.Id).Err()
	}
	if change.State != models.ContributionStatePending {
		return nil, ErrFailedPrecondition("change must be in pending state").Err()
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.DeleteTranscriptChange(ctx, request.Id)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *ContribService) RequestTranscriptChangeState(ctx context.Context, request *api.RequestTranscriptChangeStateRequest) (*emptypb.Empty, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	var oldChange *models.TranscriptChange
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		oldChange, err = s.GetTranscriptChange(ctx, request.Id)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	if err := s.validateContributionStateUpdate(claims, oldChange.Author.ID, oldChange.State, request.State); err != nil {
		return nil, err
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.UpdateTranscriptChangeState(ctx, request.Id, models.ContributionStateFromProto(request.State), request.PointsOnApprove)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *ContribService) ListAuthorContributions(ctx context.Context, request *api.ListAuthorContributionsRequest) (*api.AuthorContributionList, error) {

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
		return nil, ErrFromStore(err, "").Err()
	}
	return out, nil
}

func (s *ContribService) getClaims(ctx context.Context) (*jwt.Claims, error) {
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

func (s *ContribService) isAuthenticated(ctx context.Context) bool {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return false
	}
	if _, err := s.auth.VerifyToken(token); err == nil {
		return true
	}
	return false
}

func (s *ContribService) isAuthor(ctx context.Context, authorID string) bool {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return false
	}
	if claims, err := s.auth.VerifyToken(token); err == nil {
		return claims.AuthorID == authorID
	}
	return false
}

func (s *ContribService) isApprover(ctx context.Context) bool {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return false
	}
	if claims, err := s.auth.VerifyToken(token); err == nil {
		return claims.Approver
	}
	return false
}

func (s *ContribService) validateContributionStateUpdate(claims *jwt.Claims, currentAuthorID string, currentState models.ContributionState, requestedState api.ContributionState) error {
	if !claims.Approver {
		if currentAuthorID != claims.AuthorID {
			return ErrPermissionDenied("you are not the author of this contribution").Err()
		}
		if requestedState == api.ContributionState_STATE_APPROVED || requestedState == api.ContributionState_STATE_REJECTED {
			return ErrPermissionDenied("You are not allowed to approve/reject contributions.").Err()
		}
	}
	// if the contribution has been rejected allow the author to return it to pending.
	if currentState == models.ContributionStateRejected {
		if requestedState != api.ContributionState_STATE_PENDING {
			return ErrFailedPrecondition(fmt.Sprintf("Only rejected contributions can be reverted to pending. Actual state was: %s (requested: %s)", currentState, requestedState)).Err()
		}
	} else if currentState == models.ContributionStateApproved {

	} else {
		/// otherwise only allow it to be updated if it's in the pending or approval requested state.
		if currentState != models.ContributionStatePending && currentState != models.ContributionStateApprovalRequested {
			return ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be edited. Actual state was: %s", currentState)).Err()
		}
	}
	return nil
}

func getDonationRecipients() []*api.DonationRecipient {
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
		Name:          fmt.Sprintf("Here's some tat in a jiffy bag"),
		Criteria:      fmt.Sprintf("Earn %0.2f Points", mod.PointsSpent),
		Value:         2,
		ValueCurrency: "USD",
	}
}
