import {Component, OnDestroy, OnInit} from '@angular/core';
import {CommunityAPIClient} from "../../../../lib/api-client/services/community";
import {takeUntil} from "rxjs/operators";
import {RskArchive, RskArchiveList} from "../../../../lib/api-client/models";
import {Subject} from "rxjs";

@Component({
  selector: 'app-archive',
  standalone: false,
  templateUrl: './archive.component.html',
  styleUrl: './archive.component.scss'
})
export class ArchiveComponent implements OnInit, OnDestroy {

  private destroy$: Subject<void> = new Subject<void>();
  public archive: RskArchive[] = [];

  constructor(private apiClient: CommunityAPIClient,) {
  }

  ngOnInit(): void {
    this.apiClient.listArchive({}).pipe(takeUntil(this.destroy$)).subscribe((res: RskArchiveList): void => {
      this.archive = res.items
    })
    }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

}
