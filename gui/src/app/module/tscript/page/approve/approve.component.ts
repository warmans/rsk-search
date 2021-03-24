import { Component, EventEmitter, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import {
  RsksearchChunkContribution,
  RsksearchContributionState,
  RsksearchTscriptChunkContributionList
} from '../../../../lib/api-client/models';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Data } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { parseTranscript, Tscript } from '../../../shared/lib/tscript';

@Component({
  selector: 'app-approve',
  templateUrl: './approve.component.html',
  styleUrls: ['./approve.component.scss']
})
export class ApproveComponent implements OnInit {

  tscriptID: string;

  groupedContributions: { [index: string]: RsksearchChunkContribution[] } = {};

  approvalList: RsksearchChunkContribution[] = [];

  states = RsksearchContributionState;

  private destroy$ = new EventEmitter<any>();

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute, private titleService: Title) {
    titleService.setTitle('contribute');

    route.paramMap.pipe(takeUntil(this.destroy$)).subscribe((d: Data) => {
      this.tscriptID = d.params['tscript_id'];
      if (d.params['tscript_id']) {
        this.tscriptID = d.params['tscript_id'];
        this.loadData();
      }
    });
  }

  ngOnInit(): void {

  }

  loadData() {
    this.apiClient.searchServiceListTscriptChunkContributions({
      tscriptId: this.tscriptID,
      page: 0
    }).pipe(takeUntil(this.destroy$)).subscribe((val: RsksearchTscriptChunkContributionList) => {
      val.contributions.forEach((c) => {
        if (this.groupedContributions[c.chunkId]) {
          this.groupedContributions[c.chunkId].push(c);
        } else {
          this.groupedContributions[c.chunkId] = [c];
        }
      });
      this.updateApprovalList();
    });
  }

  updateApprovalList() {
    let approvalMap: { [index: string]: RsksearchChunkContribution } = {};
    for (let chunkId in this.groupedContributions) {

      this.groupedContributions[chunkId].forEach((co: RsksearchChunkContribution) => {
        if (!approvalMap[chunkId]) {
          approvalMap[chunkId] = co;
          return;
        }
        if (approvalMap[chunkId] === RsksearchContributionState.STATE_REJECTED) {
          return;
        }
        if (approvalMap[chunkId] === RsksearchContributionState.STATE_REJECTED && co.state === RsksearchContributionState.STATE_REQUEST_APPROVAL) {
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

  updateState(co: RsksearchChunkContribution, state: RsksearchContributionState) {
    this.apiClient.searchServiceRequestChunkContributionState({
      chunkId: co.chunkId,
      contributionId: co.id,
      body: {
        chunkId: co.chunkId,
        contributionId: co.id,
        requestState: state,
      }
    }).subscribe((result) => {
      co.state = result.state;
      this.loadData();
    });
  }

  selectChunkContribution(oldID: string, ev: any) {
    const approvalListIndex = this.approvalList.findIndex((v) => v.id === oldID);
    const oldVal = this.approvalList[approvalListIndex];

    this.approvalList[approvalListIndex] = this.groupedContributions[oldVal.chunkId].find((v) => v.id === ev.target.value);
  }
}
