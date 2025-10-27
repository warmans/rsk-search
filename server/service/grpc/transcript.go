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
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/transcript"
	"github.com/warmans/rsk-search/pkg/util"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func NewTranscriptService(
	logger *zap.Logger,
	srvCfg config.SearchServiceConfig,
	persistentDB *rw.Conn,
	staticDB *ro.Conn,
	episodeCache *data.EpisodeCache,
	auth *jwt.Auth,
) *TranscriptService {
	return &TranscriptService{
		logger:       logger,
		srvCfg:       srvCfg,
		persistentDB: persistentDB,
		staticDB:     staticDB,
		episodeCache: episodeCache,
		auth:         auth,
	}
}

type TranscriptService struct {
	logger       *zap.Logger
	srvCfg       config.SearchServiceConfig
	persistentDB *rw.Conn
	staticDB     *ro.Conn
	auth         *jwt.Auth
	episodeCache *data.EpisodeCache
}

func (s *TranscriptService) RegisterGRPC(server *grpc.Server) {
	api.RegisterTranscriptServiceServer(server, s)
}

func (s *TranscriptService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterTranscriptServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *TranscriptService) GetTranscript(ctx context.Context, request *api.GetTranscriptRequest) (*api.Transcript, error) {
	ep, err := s.episodeCache.GetEpisode(request.Epid, true)
	if errors.Is(err, data.ErrNotFound) || ep == nil {
		return nil, ErrNotFound(request.Epid)
	}
	var rawTranscript string
	if request.WithRaw {
		var err error
		rawTranscript, err = transcript.Export(ep.Transcript, ep.Synopsis, ep.Trivia)
		if err != nil {
			return nil, ErrInternal(err)
		}
	}
	lockedEpsiodeIDs, err := s.lockedEpisodeIDs(ctx)
	if err != nil {
		return nil, ErrInternal(err)
	}
	_, locked := lockedEpsiodeIDs[ep.ID()]

	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		ratings, err := s.GetPendingTranscriptRatingScores(ctx, ep.ShortID())
		if err != nil {
			return err
		}
		if ep.Ratings.Scores == nil {
			ep.Ratings.Scores = make(map[string]float32)
		}
		for k, v := range ratings.Scores {
			ep.Ratings.Scores[k] = v
		}
		return nil
	})
	if err != nil {
		return nil, ErrInternal(err)
	}

	return ep.Proto(rawTranscript, locked), nil
}

func (s *TranscriptService) GetTranscriptDialog(ctx context.Context, request *api.GetTranscriptDialogRequest) (*api.TranscriptDialog, error) {
	ep, err := s.episodeCache.GetEpisode(request.Epid, true)
	if errors.Is(err, data.ErrNotFound) || ep == nil {
		return nil, ErrNotFound(request.Epid)
	}
	dialog := []*api.Dialog{}
	for _, d := range ep.Transcript {
		if d.Position >= int64(request.Range.Start) && d.Position <= int64(request.Range.End) {
			dialog = append(dialog, d.Proto(false))
		}
		if d.Position >= int64(request.Range.End) {
			break
		}
	}

	if len(dialog) == 0 {
		return nil, ErrNotFound(request.Epid)
	}
	return &api.TranscriptDialog{
		TranscriptMeta:    ep.ShortProto(false),
		Dialog:            dialog,
		MaxDialogPosition: int32(ep.Transcript[len(ep.Transcript)-1].Position),
	}, nil
}

func (s *TranscriptService) ListTranscripts(ctx context.Context, req *api.ListTranscriptsRequest) (*api.TranscriptList, error) {
	qm, err := NewQueryModifiers(req)
	if err != nil {
		return nil, err
	}
	var episodeMeta []*models.EpisodeMeta
	err = s.staticDB.WithStore(func(s *ro.Store) error {
		var err error
		episodeMeta, err = s.ListEpisodes(ctx, qm)
		return err
	})
	if err != nil {
		return nil, ErrInternal(err)
	}

	pendingRatings := make(map[string]models.Ratings)
	if err := s.persistentDB.WithStore(func(s *rw.Store) error {
		pendingRatings, err = s.ListPendingRatings(ctx)
		return err
	}); err != nil {
		return nil, ErrInternal(err)
	}

	el := &api.TranscriptList{
		Episodes: []*api.ShortTranscript{},
	}
	for _, meta := range episodeMeta {
		ep, err := s.episodeCache.GetEpisode(meta.ID(), false)
		if err != nil {
			s.logger.Error("failed to get episode from cache", zap.Error(err))
			continue
		}
		// merge pending ratings if any
		if ep.Ratings.Scores == nil {
			ep.Ratings.Scores = make(map[string]float32)
		}
		if pending, ok := pendingRatings[ep.ShortID()]; ok {
			for k, v := range pending.Scores {
				ep.Ratings.Scores[k] = v
			}
		}
		el.Episodes = append(el.Episodes, ep.ShortProto(true))
	}
	return el, nil
}

func (s *TranscriptService) ListChunkedTranscripts(ctx context.Context, _ *emptypb.Empty) (*api.ChunkedTranscriptList, error) {
	el := &api.ChunkedTranscriptList{
		Chunked: []*api.ChunkedTranscriptStats{},
	}
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		eps, err := s.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, e := range eps {
			el.Chunked = append(el.Chunked, e.Proto())
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return el, nil
}

func (s *TranscriptService) GetChunkedTranscriptChunkStats(ctx context.Context, _ *emptypb.Empty) (*api.ChunkStats, error) {
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
		return nil, ErrFromStore(err, "")
	}
	return stats.Proto(), nil
}

func (s *TranscriptService) GetTranscriptChunk(ctx context.Context, request *api.GetTranscriptChunkRequest) (*api.Chunk, error) {
	var chunk *models.Chunk
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		chunk, err = s.GetChunk(ctx, request.Id)
		if err != nil {
			return err
		}
		if chunk == nil {
			return ErrNotFound(request.Id)
		}
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id)
	}
	return chunk.Proto(), nil
}

func (s *TranscriptService) ListTranscriptChunks(ctx context.Context, request *api.ListTranscriptChunksRequest) (*api.TranscriptChunkList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}
	if qm.Filter != nil {
		qm.Filter = filter.And(filter.Eq("tscript_id", filter.String(request.ChunkedTranscriptId)), qm.Filter)
	} else {
		qm.Filter = filter.Eq("tscript_id", filter.String(request.ChunkedTranscriptId))
	}
	if qm.Sorting == nil {
		qm.Sorting = &common.Sorting{Field: "start_second", Direction: common.SortAsc}
	}

	out := &api.TranscriptChunkList{
		Chunks: make([]*api.Chunk, 0),
	}
	if err := s.persistentDB.WithStore(func(store *rw.Store) error {
		chunks, err := store.ListChunks(ctx, qm)
		for _, v := range chunks {
			out.Chunks = append(out.Chunks, v.Proto())
		}
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "")
	}
	return out, nil
}

func (s *TranscriptService) ListChunkContributions(ctx context.Context, request *api.ListChunkContributionsRequest) (*api.ChunkContributionList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
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
		return nil, ErrFromStore(err, "")
	}
	return out, nil
}

func (s *TranscriptService) GetChunkContribution(ctx context.Context, request *api.GetChunkContributionRequest) (*api.ChunkContribution, error) {

	var contrib *models.ChunkContribution
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		contrib, err = s.GetChunkContribution(ctx, request.ContributionId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId)
	}
	if err := checkReadingAllowed(contrib.State, IsApprover(ctx, s.auth), IsAuthor(ctx, s.auth, contrib.Author.ID)); err != nil {
		return nil, err
	}
	return contrib.Proto(), nil
}

func (s *TranscriptService) CreateChunkContribution(ctx context.Context, request *api.CreateChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		stats, err := s.GetAuthorStats(ctx, claims.AuthorID)
		if err != nil {
			return err
		}
		if stats.ContributionsInLastHour > 10 {
			return ErrRateLimited()
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}

	if err := transcript.Validate(bufio.NewScanner(bytes.NewBufferString(request.Transcript))); err != nil {
		return nil, ErrInvalidRequestField("transcript", err)
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
		return nil, ErrFromStore(err, "")
	}
	return contrib.Proto(), nil
}

func (s *TranscriptService) UpdateChunkContribution(ctx context.Context, request *api.UpdateChunkContributionRequest) (*api.ChunkContribution, error) {

	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, request.ContributionId)
	}

	// validate change is allowed
	if err := s.validateContributionStateUpdate(claims, contrib.Author.ID, contrib.State, request.State); err != nil {
		return nil, err
	}

	// allow invalid transcript while the contribution is still pending.
	if request.State != api.ContributionState_STATE_PENDING {
		if err := transcript.Validate(bufio.NewScanner(bytes.NewBufferString(request.Transcript))); err != nil {
			return nil, ErrInvalidRequestField("transcript", err)
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
		if err := s.createAuthorNotification(ctx, tx, contrib.Author.ID, request.State, "chunk contribution", ""); err != nil {
			return err
		}
		return tx.UpdateChunkActivity(ctx, contrib.ChunkID, rw.ActivityFromState(contrib.State))
	})
	if err != nil {
		return nil, ErrFromStore(err, contrib.ID)
	}

	return contrib.Proto(), nil
}

func (s *TranscriptService) RequestChunkContributionState(ctx context.Context, request *api.RequestChunkContributionStateRequest) (*api.ChunkContribution, error) {

	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, request.ContributionId)
	}
	if err := s.validateContributionStateUpdate(claims, contrib.Author.ID, contrib.State, request.RequestState); err != nil {
		return nil, err
	}
	if request.Comment != "" && !claims.Approver {
		return nil, ErrPermissionDenied("Only an approver can set a state comment.")
	}
	err = s.persistentDB.WithStore(func(tx *rw.Store) error {

		contrib.State = models.ContributionStateFromProto(request.RequestState)
		contrib.StateComment = request.Comment

		if err := tx.UpdateChunkContributionState(ctx, contrib.ID, contrib.State, contrib.StateComment); err != nil {
			return err
		}
		if err := s.createAuthorNotification(ctx, tx, contrib.Author.ID, request.RequestState, "chunk contribution", ""); err != nil {
			return err
		}
		return tx.UpdateChunkActivity(ctx, contrib.ChunkID, rw.ActivityFromState(contrib.State))
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId)
	}
	return contrib.Proto(), nil
}

func (s *TranscriptService) DeleteChunkContribution(ctx context.Context, request *api.DeleteChunkContributionRequest) (*emptypb.Empty, error) {
	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, request.ContributionId)
	}
	if claims.AuthorID != contrib.Author.ID {
		return nil, ErrPermissionDenied("you are not the author of this contribution")
	}
	if contrib.State != models.ContributionStatePending {
		return nil, ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be deleted. Actual state was: %s", contrib.State))
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.DeleteContribution(ctx, request.ContributionId)
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return &emptypb.Empty{}, nil
}

func (s *TranscriptService) ListTranscriptChanges(ctx context.Context, request *api.ListTranscriptChangesRequest) (*api.TranscriptChangeList, error) {

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
			if v.State == models.ContributionStatePending && !IsAuthor(ctx, s.auth, v.Author.ID) {
				continue
			}
			out.Changes = append(out.Changes, v.ShortProto())
		}
		return err
	}); err != nil {
		return nil, ErrFromStore(err, "")
	}
	return out, nil
}

func (s *TranscriptService) GetTranscriptChange(ctx context.Context, request *api.GetTranscriptChangeRequest) (*api.TranscriptChange, error) {

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
		return nil, ErrFromStore(err, request.Id)
	}
	if err := checkReadingAllowed(change.State, IsApprover(ctx, s.auth), IsAuthor(ctx, s.auth, change.Author.ID)); err != nil {
		return nil, err
	}

	return change.Proto(), nil
}

func (s *TranscriptService) GetTranscriptChangeDiff(ctx context.Context, request *api.GetTranscriptChangeDiffRequest) (*api.TranscriptChangeDiff, error) {

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
		return nil, ErrFromStore(err, request.Id)
	}
	if err := checkReadingAllowed(newTranscript.State, IsApprover(ctx, s.auth), IsAuthor(ctx, s.auth, newTranscript.Author.ID)); err != nil {
		return nil, err
	}

	oldTranscript, err := s.episodeCache.GetEpisode(newTranscript.EpID, true)
	if err != nil {
		return nil, ErrNotFound(newTranscript.EpID)
	}

	oldTranscriptRaw, err := transcript.Export(oldTranscript.Transcript, oldTranscript.Synopsis, oldTranscript.Trivia)
	if err != nil {
		return nil, err
	}

	diffs := []string{}

	nameEdits := myers.ComputeEdits(
		span.URIFromPath("NAME"),
		oldTranscript.Name,
		newTranscript.Name,
	)
	if len(nameEdits) > 0 {
		nameDiff := fmt.Sprint(
			gotextdiff.ToUnified(
				"NAME",
				"NAME",
				oldTranscript.Name,
				nameEdits,
			),
		)
		diffs = append(diffs, nameDiff)
	}

	// summary diff
	summaryEdits := myers.ComputeEdits(
		span.URIFromPath("SUMMARY"),
		oldTranscript.Summary,
		newTranscript.Summary,
	)
	if len(summaryEdits) > 0 {
		summaryDiff := fmt.Sprint(
			gotextdiff.ToUnified(
				"SUMMARY",
				"SUMMARY",
				oldTranscript.Summary,
				summaryEdits,
			),
		)
		diffs = append(diffs, summaryDiff)
	}

	// transcript diff
	transcriptEdits := myers.ComputeEdits(
		span.URIFromPath("TRANSCRIPT"),
		oldTranscriptRaw,
		newTranscript.Transcription,
	)
	if len(transcriptEdits) > 0 {
		transcriptDiff := fmt.Sprint(
			gotextdiff.ToUnified(
				"TRANSCRIPT",
				"TRANSCRIPT",
				oldTranscriptRaw,
				transcriptEdits,
			),
		)
		diffs = append(diffs, transcriptDiff)
	}

	return &api.TranscriptChangeDiff{Diffs: diffs}, nil
}

func (s *TranscriptService) CreateTranscriptChange(ctx context.Context, request *api.CreateTranscriptChangeRequest) (*api.TranscriptChange, error) {
	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	if err := s.validateLockedState(ctx, request.Epid); err != nil {
		return nil, err
	}

	if len(request.Name) > 128 {
		return nil, ErrInvalidRequestField("name", nil, "Cannot be more than 128 characters")
	}
	if len(request.Summary) > 2048 {
		return nil, ErrInvalidRequestField("summary", nil, "Cannot be more than 2048 characters")
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
			return ErrFailedPrecondition("multiple changes cannot exist at once. Try again once the current change has been processed.")
		}
		change, err = s.CreateTranscriptChange(ctx, &models.TranscriptChangeCreate{
			AuthorID:          claims.AuthorID,
			EpID:              request.Epid,
			Summary:           request.Summary,
			Name:              request.Name,
			Transcription:     request.Transcript,
			TranscriptVersion: request.TranscriptVersion,
		})
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return change.Proto(), nil
}

func (s *TranscriptService) UpdateTranscriptChange(ctx context.Context, request *api.UpdateTranscriptChangeRequest) (*api.TranscriptChange, error) {

	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}

	if len(request.Name) > 128 {
		return nil, ErrInvalidRequestField("name", nil, "Cannot be more than 128 characters")
	}
	if len(request.Summary) > 2048 {
		return nil, ErrInvalidRequestField("summary", nil, "Cannot be more than 2048 characters")
	}

	var oldChange *models.TranscriptChange
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		oldChange, err = s.GetTranscriptChange(ctx, request.Id)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id)
	}

	if err := s.validateLockedState(ctx, oldChange.EpID); err != nil {
		return nil, err
	}

	// validate change is allowed
	if err := s.validateContributionStateUpdate(claims, oldChange.Author.ID, oldChange.State, request.State); err != nil {
		return nil, err
	}

	// allow invalid transcript while the contribution is still pending.
	if request.State != api.ContributionState_STATE_PENDING {
		if err := transcript.Validate(bufio.NewScanner(bytes.NewBufferString(request.Transcript))); err != nil {
			return nil, ErrInvalidRequestField("transcript", err)
		}
	}

	var updatedChange *models.TranscriptChange
	err = s.persistentDB.WithStore(func(tx *rw.Store) error {

		if err := s.createAuthorNotification(ctx, tx, oldChange.Author.ID, request.State, "transcript change", ""); err != nil {
			return err
		}
		var err error
		updatedChange, err = tx.UpdateTranscriptChange(ctx, &models.TranscriptChangeUpdate{
			ID:            request.Id,
			Name:          request.Name,
			Summary:       request.Summary,
			Transcription: request.Transcript,
			State:         models.ContributionStateFromProto(request.State),
		}, request.PointsOnApprove)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return updatedChange.Proto(), nil
}

func (s *TranscriptService) DeleteTranscriptChange(ctx context.Context, request *api.DeleteTranscriptChangeRequest) (*emptypb.Empty, error) {
	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, request.Id)
	}
	if change.Author.ID != claims.AuthorID {
		return nil, ErrNotFound(request.Id)
	}
	if change.State != models.ContributionStatePending {
		return nil, ErrFailedPrecondition("change must be in pending state")
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.DeleteTranscriptChange(ctx, request.Id)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id)
	}
	return &emptypb.Empty{}, nil
}

func (s *TranscriptService) RequestTranscriptChangeState(ctx context.Context, request *api.RequestTranscriptChangeStateRequest) (*emptypb.Empty, error) {
	claims, err := GetClaims(ctx, s.auth)
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
		return nil, ErrFromStore(err, request.Id)
	}

	if err := s.validateLockedState(ctx, oldChange.EpID); err != nil {
		return nil, err
	}

	if err := s.validateContributionStateUpdate(claims, oldChange.Author.ID, oldChange.State, request.State); err != nil {
		return nil, err
	}
	err = s.persistentDB.WithStore(func(tx *rw.Store) error {
		if err := s.createAuthorNotification(ctx, tx, oldChange.Author.ID, request.State, "transcript change", ""); err != nil {
			return err
		}
		return tx.UpdateTranscriptChangeState(ctx, request.Id, models.ContributionStateFromProto(request.State), request.PointsOnApprove)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id)
	}
	return &emptypb.Empty{}, nil
}

// if an episode is currently being transcribed, mark it as locked to prevent changes being submitted before
// all chunks have been completed.
func (s *TranscriptService) lockedEpisodeIDs(ctx context.Context) (map[string]struct{}, error) {
	inProgressTscriptIDs := map[string]struct{}{}
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		inProgressTscripts, err := s.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, v := range inProgressTscripts {
			inProgressTscriptIDs[models.EpID(v.Publication, v.Series, v.Episode)] = struct{}{}
		}
		return err
	})
	return inProgressTscriptIDs, err
}

// if an episode is currently being transcribed, mark it as locked to prevent changes being submitted before
// all chunks have been completed.
func (s *TranscriptService) validateLockedState(ctx context.Context, epID string) error {
	return s.persistentDB.WithStore(func(s *rw.Store) error {
		inProgressTscripts, err := s.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, v := range inProgressTscripts {
			if epID == models.EpID(v.Publication, v.Series, v.Episode) {
				return ErrFailedPrecondition("episode is locked")
			}
		}
		return err
	})
}

func (s *TranscriptService) validateContributionStateUpdate(claims *jwt.Claims, currentAuthorID string, currentState models.ContributionState, requestedState api.ContributionState) error {
	if !claims.Approver {
		if currentAuthorID != claims.AuthorID {
			return ErrPermissionDenied("you are not the author of this contribution")
		}
		if requestedState == api.ContributionState_STATE_APPROVED || requestedState == api.ContributionState_STATE_REJECTED {
			return ErrPermissionDenied("You are not allowed to approve/reject contributions.")
		}
	}
	// if the contribution has been rejected, allow the author to return it to pending.
	switch currentState {
	case models.ContributionStateRejected:
		if requestedState != api.ContributionState_STATE_PENDING {
			return ErrFailedPrecondition(fmt.Sprintf("Only rejected contributions can be reverted to pending. Actual state was: %s (requested: %s)", currentState, requestedState))
		}
	case models.ContributionStateApproved:
		// do nothing
	default:
		// otherwise only allow it to be updated if it's in the pending or approval requested state.
		if currentState != models.ContributionStatePending && currentState != models.ContributionStateApprovalRequested {
			return ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be edited. Actual state was: %s", currentState))
		}
	}
	return nil
}

func (s *TranscriptService) createAuthorNotification(
	ctx context.Context,
	tx *rw.Store,
	authorID string,
	state api.ContributionState,
	entity string,
	comment string,
) error {

	var message string
	var kind string
	switch state {
	case api.ContributionState_STATE_REJECTED:
		var reason string
		if comment != "" {
			reason = fmt.Sprintf("The given reason was: %s", comment)
		} else {
			reason = "No reason was given."
		}
		message = fmt.Sprintf("Sorry, your %s was rejected. %s If you think was is a mistake you can edit/re-submit the change from your profile page.", entity, reason)
		kind = api.Notification_WARNING.String()
	case api.ContributionState_STATE_APPROVED:
		message = fmt.Sprintf("Great, your %s was accepted and will be merged soon.", entity)
		kind = api.Notification_CONFIRMATION.String()
	default:
		return nil
	}
	return tx.CreateAuthorNotification(ctx, models.AuthorNotificationCreate{
		AuthorID:       authorID,
		Kind:           kind,
		Message:        message,
		ClickThoughURL: util.StringP("/me"),
	})
}

func (s *TranscriptService) SetTranscriptRatingScore(ctx context.Context, request *api.SetTranscriptRatingScoreRequest) (*emptypb.Empty, error) {
	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.UpsertTranscriptRatingScore(ctx, request.Epid, claims.AuthorID, request.Score, request.Delete)
	})
	if err != nil {
		return nil, ErrInternal(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *TranscriptService) BulkSetTranscriptRatingScore(ctx context.Context, request *api.BulkSetTranscriptRatingScoreRequest) (*emptypb.Empty, error) {

	if err := s.authorizeSystemRequest(ctx); err != nil {
		return nil, err
	}

	if request.OauthSource == "" {
		return nil, fmt.Errorf("outh source is required to match the given scores to existing authors")
	}

	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		for author, rating := range request.Scores {
			id, err := s.GetOrCreateAuthorID(ctx, author, request.OauthSource)
			if err != nil {
				return err
			}
			if err := s.UpsertTranscriptRatingScore(ctx, request.Epid, id, rating, false); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, ErrInternal(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *TranscriptService) BulkSetTranscriptTags(ctx context.Context, request *api.BulkSetTranscriptTagsRequest) (*emptypb.Empty, error) {
	if err := s.authorizeSystemRequest(ctx); err != nil {
		return nil, err
	}
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		for _, v := range request.Tags {
			ts, err := time.ParseDuration(v.Timestamp)
			if err != nil {
				return fmt.Errorf("failed to pass timestamp for tag %s: %w", v.Name, err)
			}
			if err := s.UpsertTranscriptTag(ctx, request.Epid, v.Name, ts); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, ErrInternal(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *TranscriptService) authorizeSystemRequest(ctx context.Context) error {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return ErrUnauthorized("no token provided")
	}
	claims, err := s.auth.VerifyToken(token)
	if err != nil {
		return ErrUnauthorized(err.Error())
	}
	if !claims.System {
		return ErrUnauthorized("Not a system token")
	}
	return nil
}

func checkReadingAllowed(state models.ContributionState, isApprover bool, isAuthor bool) error {
	if !isApprover && (state == models.ContributionStatePending && !isAuthor) {
		return ErrPermissionDenied("you cannot view another author's contribution diff when it is in the pending state")
	}
	return nil
}
