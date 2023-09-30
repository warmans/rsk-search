import {Component, EventEmitter} from '@angular/core';
import {RskAuthorContribution, RskTranscriptChange} from "../../../../lib/api-client/models";
import {And, Eq, Neq} from "../../../../lib/filter-dsl/filter";
import {Bool, Str} from "../../../../lib/filter-dsl/value";
import {takeUntil} from "rxjs/operators";
import {SearchAPIClient} from "../../../../lib/api-client/services/search";

@Component({
  selector: 'app-changes',
  templateUrl: './changes.component.html',
  styleUrls: ['./changes.component.scss']
})
export class ChangesComponent {

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  loading: boolean[] = [];

  pendingChanges: RskTranscriptChange[] = [];
  recentContributions: RskAuthorContribution[] = [];

  activeTab: 'pending' | 'recent' = 'recent';

  public constructor(private apiClient: SearchAPIClient) {
  }

  public ngOnInit() {
    this.getPendingChanges();
    this.getRecentContributions();
  }

  private getRecentContributions() {
    this.loading.push(true);
    this.apiClient.listAuthorContributions({
      pageSize: 10,
      sortField: 'created_at',
      sortDirection: 'DESC'
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((res) => {
      this.recentContributions = res.contributions;
      if (res.contributions.length > 0) {
        this.activeTab = 'pending';
      }
    }).add(() => {
      this.loading.pop();
    });
  }

  private getPendingChanges() {
    this.loading.push(true);
    this.apiClient.listTranscriptChanges({
      filter: And(
        Eq('merged', Bool(false)),
        Neq('state', Str('pending')),
        Neq('state', Str('rejected')),
      ).print()
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((res) => {
      this.pendingChanges = res.changes;
    }).add(() => {
      this.loading.pop();
    });
  }

}
