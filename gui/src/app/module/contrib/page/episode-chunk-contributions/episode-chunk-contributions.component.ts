import {Component, EventEmitter, OnInit} from '@angular/core';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {RskChunk, RskChunkContribution, RskChunkContributionList, RskContributionState, RskTranscriptChunkList} from 'src/app/lib/api-client/models';
import {takeUntil} from 'rxjs/operators';
import {ActivatedRoute, Data} from '@angular/router';
import {Title} from '@angular/platform-browser';
import {parseTranscript, Tscript} from '../../../shared/lib/tscript';
import {SessionService} from '../../../core/service/session/session.service';
import {And, Eq} from 'src/app/lib/filter-dsl/filter';
import {Str} from 'src/app/lib/filter-dsl/value';

@Component({
    selector: 'app-episode-chunk-contributions',
    templateUrl: './episode-chunk-contributions.component.html',
    styleUrls: ['./episode-chunk-contributions.component.scss'],
    standalone: false
})
export class EpisodeChunkContributions implements OnInit {

  chunkedTranscriptID: string;

  chunks: Array<RskChunk & {selectedContribution?: number}> = [];

  groupedContributions: { [index: string]: RskChunkContribution[] } = {};

  approvalList: RskChunkContribution[] = [];

  states = RskContributionState;

  approver: boolean = false;

  loading: boolean[] = [];

  stateFilter: string | undefined;

  rejectCallback: (contributionId: string, comment: string) => void = (contributionId: string, comment: string) => {
    this.updateState(contributionId, RskContributionState.STATE_REJECTED, comment);
  };

  private destroy$ = new EventEmitter<any>();

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute, private titleService: Title, private session: SessionService) {

    session.onTokenChange.pipe(takeUntil(this.destroy$)).subscribe((token: string) => {
      this.approver = session.getClaims()?.approver || false;
    });

    route.paramMap.pipe(takeUntil(this.destroy$)).subscribe((d: Data) => {
      this.chunkedTranscriptID = d.params['tscript_id'];
      titleService.setTitle(`Contributions for ${this.chunkedTranscriptID}`);

      if (d.params['tscript_id']) {
        this.chunkedTranscriptID = d.params['tscript_id'];
        this.loadData();
      }
    });
  }

  ngOnInit(): void {

  }

  loadData() {

    let filter = Eq('tscript_id', Str(this.chunkedTranscriptID));
    if (this.stateFilter) {
      filter = And(filter, Eq('state', Str(this.stateFilter)));
    }


    this.loading.push(true);
    this.apiClient.listTranscriptChunks({
      chunkedTranscriptId: this.chunkedTranscriptID,
      sortField: 'start_second',
      sortDirection: 'asc',
      pageSize: 200
    }).pipe(takeUntil(this.destroy$)).subscribe((resp: RskTranscriptChunkList) => {
      this.chunks = resp.chunks;
    }).add(() => this.loading.pop());

    this.loading.push(true);
    this.apiClient.listChunkContributions({
      filter: filter.print(),
      pageSize: 200,
    }).pipe(takeUntil(this.destroy$)).subscribe((val: RskChunkContributionList) => {
      this.groupedContributions = {};
      val.contributions.forEach((c) => {
        if (c.state === RskContributionState.STATE_PENDING || c.state === RskContributionState.STATE_REJECTED) {
          return;
        }
        if (this.groupedContributions[c.chunkId]) {
          this.groupedContributions[c.chunkId].push(c);
        } else {
          this.groupedContributions[c.chunkId] = [c];
        }
      });
    }).add(() => this.loading.pop());
  }

  // unused - but might need to be re-implemented at some point
  updateApprovalList() {
    let approvalMap: { [index: string]: RskChunkContribution } = {};
    for (let chunkId in this.groupedContributions) {
      this.groupedContributions[chunkId].forEach((co: RskChunkContribution) => {
        if (!approvalMap[chunkId]) {
          approvalMap[chunkId] = co;
          return;
        }
        // do not replace the current value if it is already approved.
        if (approvalMap[chunkId].state === RskContributionState.STATE_APPROVED) {
          return;
        }
        if (approvalMap[chunkId].state === RskContributionState.STATE_REJECTED && co.state !== RskContributionState.STATE_REJECTED) {
          approvalMap[chunkId] = co;
        }
        if (co.state === RskContributionState.STATE_APPROVED) {
          approvalMap[chunkId] = co;
        }
      });
    }
    this.approvalList = [];
    for (let chunkId in approvalMap) {
      this.approvalList.push(approvalMap[chunkId]);
    }
  }

  parseTscript(raw: string): Tscript {
    return parseTranscript(raw);
  }

  updateState(contributionId: string, state: RskContributionState, comment?: string) {
    this.loading.push(true);
    this.apiClient.requestChunkContributionState({
        contributionId: contributionId,
        body: {
          requestState: state,
          comment: comment,
      }
    }).subscribe((result) => {
      this.loadData();
    }).add(() => this.loading.pop());
  }

  selectChunkContribution(oldID: string, ev: any) {
    const approvalListIndex = this.approvalList.findIndex((v) => v.id === oldID);
    const oldVal = this.approvalList[approvalListIndex];

    this.approvalList[approvalListIndex] = this.groupedContributions[oldVal.chunkId].find((v) => v.id === ev.target.value);
  }

  setStateFilter(state: string | undefined) {
    this.stateFilter = state === this.stateFilter ? undefined : state;
    this.loadData();
  }
}
