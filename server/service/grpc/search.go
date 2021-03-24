package grpc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/tscript"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
	"strings"
)

func NewSearchService(
	searchBackend *search.Search,
	store *ro.Conn,
	persistentDB *rw.Conn,
	csrfCache *oauth.CSRFTokenCache,
	auth *jwt.Auth,
) *SearchService {
	return &SearchService{
		searchBackend: searchBackend,
		staticDB:      store,
		persistentDB:  persistentDB,
		csrfCache:     csrfCache,
		auth:          auth,
	}
}

type SearchService struct {
	searchBackend *search.Search
	staticDB      *ro.Conn
	persistentDB  *rw.Conn
	csrfCache     *oauth.CSRFTokenCache
	auth          *jwt.Auth
}

func (s *SearchService) RegisterGRPC(server *grpc.Server) {
	api.RegisterSearchServiceServer(server, s)
}

func (s *SearchService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *SearchService) GetSearchMetadata(ctx context.Context, empty *emptypb.Empty) (*api.SearchMetadata, error) {
	return meta.GetSearchMeta().Proto(), nil
}

func (s *SearchService) ListFieldValues(ctx context.Context, request *api.ListFieldValuesRequest) (*api.FieldValueList, error) {
	vals, err := s.searchBackend.ListTerms(request.Field, request.Prefix)
	if err != nil {
		return nil, err
	}
	return &api.FieldValueList{Values: vals.Proto()}, nil
}

func (s *SearchService) Search(ctx context.Context, request *api.SearchRequest) (*api.SearchResultList, error) {
	f, err := filter.Parse(request.Query)
	if err != nil {
		return nil, ErrInvalidRequestField("query", err.Error()).Err()
	}
	if err := checkWhy(f); err != nil {
		return nil, err
	}
	return s.searchBackend.Search(ctx, f, request.Page)
}

func checkWhy(f filter.Filter) error {
	visitor := filter.NewExtractFilterVisitor(f)
	filters, err := visitor.ExtractCompFilters("content")
	if err != nil {
		return nil // don't fail because of this stupid feature
	}
	for _, v := range filters {
		if strings.TrimSpace(strings.Trim(v.Value.String(), `"?`)) == "why" {
			return ErrServerConfused().Err()
		}
	}
	return nil
}

func (s *SearchService) GetEpisode(ctx context.Context, request *api.GetEpisodeRequest) (*api.Episode, error) {
	var ep *models.Episode
	err := s.staticDB.WithStore(func(s *ro.Store) error {
		var err error
		ep, err = s.GetEpisode(ctx, request.Id)
		if err != nil {
			return err
		}
		if ep == nil {
			return ErrNotFound(request.Id).Err()
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return ep.Proto(), nil
}

func (s *SearchService) ListEpisodes(ctx context.Context, request *api.ListEpisodesRequest) (*api.EpisodeList, error) {
	el := &api.EpisodeList{
		Episodes: []*api.ShortEpisode{},
	}
	err := s.staticDB.WithStore(func(s *ro.Store) error {
		eps, err := s.ListEpisodes(ctx)
		if err != nil {
			return err
		}
		for _, e := range eps {
			el.Episodes = append(el.Episodes, e.ShortProto())
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return el, nil
}

func (s *SearchService) GetTscriptChunkStats(ctx context.Context, empty *emptypb.Empty) (*api.ChunkStats, error) {
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

func (s *SearchService) GetTscriptChunk(ctx context.Context, request *api.GetTscriptChunkRequest) (*api.TscriptChunk, error) {
	var chunk *models.Chunk
	var tscriptID string
	var contributionCount int32

	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		chunk, tscriptID, err = s.GetChunk(ctx, request.Id)
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
	return chunk.Proto(tscriptID, contributionCount), nil
}

func (s *SearchService) CreateChunkContribution(ctx context.Context, request *api.CreateChunkContributionRequest) (*api.ChunkContribution, error) {

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

	lines, _, err := tscript.Import(bufio.NewScanner(bytes.NewBufferString(request.Transcript)))
	if err != nil {
		return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
	}
	if len(lines) == 0 {
		return nil, ErrInvalidRequestField("transcript", "no valid lines parsed from transcript").Err()
	}

	contribution := &models.Contribution{
		AuthorID:      claims.AuthorID,
		ChunkID:       request.ChunkId,
		Transcription: request.Transcript,
		State:         models.ContributionStatePending,
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.CreateContribution(ctx, contribution)
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}

	return contribution.Proto(), nil
}

func (s *SearchService) UpdateChunkContribution(ctx context.Context, request *api.UpdateChunkContributionRequest) (*api.ChunkContribution, error) {

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
	if contrib.AuthorID != claims.AuthorID {
		return nil, ErrPermissionDenied("you are not the author of this contribution").Err()
	}
	if contrib.State != models.ContributionStatePending && contrib.State != models.ContributionStateApprovalRequested {
		return nil, ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be edited. Actual state was: %s", contrib.State)).Err()
	}
	lines, _, err := tscript.Import(bufio.NewScanner(bytes.NewBufferString(request.Transcript)))
	if err != nil {
		return nil, ErrInvalidRequestField("transcript", err.Error()).Err()
	}
	if len(lines) == 0 {
		return nil, ErrInvalidRequestField("transcript", "no valid lines parsed from transcript").Err()
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		mod := &models.Contribution{
			ID:            contrib.ID,
			AuthorID:      contrib.AuthorID,
			Transcription: request.Transcript,
			State:         contrib.State,
		}
		return s.UpdateContribution(ctx, mod)
	})
	if err != nil {
		return nil, ErrFromStore(err, contrib.ID).Err()
	}

	fmt.Print("TRANSCRIPT", request.Transcript)

	contrib.Transcription = request.Transcript
	return contrib.Proto(), nil
}

func (s *SearchService) RequestChunkContributionState(ctx context.Context, request *api.RequestChunkContributionStateRequest) (*api.ChunkContribution, error) {

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
	if contrib.AuthorID != claims.AuthorID {
		return nil, ErrPermissionDenied("you are not the author of this contribution").Err()
	}
	if contrib.State != models.ContributionStatePending && contrib.State != models.ContributionStateApprovalRequested {
		return nil, ErrFailedPrecondition(fmt.Sprintf("Only pending contributions can be edited. Actual state was: %s", contrib.State)).Err()
	}
	// ensure state is updated in response and db
	if request.RequestState == api.ContributionState_STATE_REQUEST_APPROVAL {
		contrib.State = models.ContributionStateApprovalRequested
	}
	if request.RequestState == api.ContributionState_STATE_PENDING {
		contrib.State = models.ContributionStatePending
	}
	if request.RequestState == api.ContributionState_STATE_APPROVED || request.RequestState == api.ContributionState_STATE_REJECTED {
		if claims.Approver == false {
			return nil, ErrPermissionDenied("you are not an approver").Err()
		}
		if request.RequestState == api.ContributionState_STATE_APPROVED {
			contrib.State = models.ContributionStateApproved
		}
		if request.RequestState == api.ContributionState_STATE_REJECTED {
			contrib.State = models.ContributionStateRejected
		}
	}

	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.UpdateContributionState(ctx, contrib.ID, contrib.State)
	})
	if err != nil {
		return nil, ErrFromStore(err, request.ContributionId).Err()
	}
	return contrib.Proto(), nil
}

func (s *SearchService) ListTscriptChunkContributions(ctx context.Context, request *api.ListTscriptChunkContributionsRequest) (*api.TscriptChunkContributionList, error) {
	var list []*models.Contribution
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		list, err = s.ListApprovableTscriptContributions(ctx, request.TscriptId, request.Page)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.TscriptId).Err()
	}
	out := &api.TscriptChunkContributionList{
		Contributions: make([]*api.ChunkContribution, len(list)),
	}
	for k, v := range list {
		out.Contributions[k] = v.Proto()
	}
	return out, nil
}

func (s *SearchService) ListAuthorContributions(ctx context.Context, request *api.ListAuthorContributionsRequest) (*api.ChunkContributionList, error) {

	var list []*models.Contribution
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		list, err = s.ListAuthorContributions(ctx, request.AuthorId, request.Page)
		return err
	})
	if err != nil {
		return nil, ErrFromStore(err, request.AuthorId).Err()
	}
	out := &api.ChunkContributionList{
		Contributions: make([]*api.ShortChunkContribution, len(list)),
	}
	for k, v := range list {
		out.Contributions[k] = v.ShortProto()
	}
	return out, nil
}

func (s *SearchService) GetChunkContribution(ctx context.Context, request *api.GetChunkContributionRequest) (*api.ChunkContribution, error) {
	var contrib *models.Contribution
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
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
	return contrib.Proto(), nil
}

func (s *SearchService) SubmitDialogCorrection(ctx context.Context, request *api.SubmitDialogCorrectionRequest) (*emptypb.Empty, error) {
	panic("implement me")

}

func (s *SearchService) GetRedditAuthURL(ctx context.Context, empty *emptypb.Empty) (*api.RedditAuthURL, error) {

	returnURL := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md["grpcgateway-referer"]) > 0 {
		// we don't want to keep the query or fragment
		if parsed, err := url.Parse(md["grpcgateway-referer"][0]); err == nil {
			parsed.RawQuery = ""
			parsed.RawFragment = ""
			returnURL = parsed.String()
		}
	}
	fmt.Println("RAW", md["grpcgateway-referer"])
	fmt.Println("RETURN", returnURL)

	return &api.RedditAuthURL{
		Url: fmt.Sprintf(
			"https://www.reddit.com/api/v1/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=%s&duration=temporary&scope=identity",
			oauth.RedditApplicationID,
			s.csrfCache.NewCSRFToken(returnURL),
			oauth.RedditReturnURI,
		),
	}, nil
}

func (s *SearchService) getClaims(ctx context.Context) (*jwt.Claims, error) {
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
