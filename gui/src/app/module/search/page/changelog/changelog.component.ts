import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { RskChangelog, RskChangelogList } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';

@Component({
    selector: 'app-changelog',
    templateUrl: './changelog.component.html',
    styleUrls: ['./changelog.component.scss'],
    standalone: false
})
export class ChangelogComponent implements OnInit, OnDestroy {

  changelogs: RskChangelog[] = [];

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient) { }

  ngOnInit(): void {

    this.apiClient.listChangelogs({pageSize: 10}).pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskChangelogList) => {
      this.changelogs = res.changelogs || [];
    });
  }
  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }


}
