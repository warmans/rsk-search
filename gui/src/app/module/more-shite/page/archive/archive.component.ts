import {Component, OnDestroy, OnInit} from '@angular/core';
import {CommunityAPIClient} from "../../../../lib/api-client/services/community";
import {takeUntil} from "rxjs/operators";
import {RskArchive, RskArchiveList} from "../../../../lib/api-client/models";
import {Subject} from "rxjs";
import {ActivatedRoute, ParamMap, Router} from "@angular/router";

const PAGE_SIZE = 10;
const MAX_PAGINATION_LINKS = 10;

@Component({
  selector: 'app-archive',
  standalone: false,
  templateUrl: './archive.component.html',
  styleUrl: './archive.component.scss'
})
export class ArchiveComponent implements OnInit, OnDestroy {

  private destroy$: Subject<void> = new Subject<void>();

  archive: RskArchive[] = [];
  currentPage: number;
  pages: number[];
  morePages: boolean;
  maxPages: number = MAX_PAGINATION_LINKS;

  constructor(private apiClient: CommunityAPIClient, private route: ActivatedRoute) {
  }

  ngOnInit(): void {
    this.route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((route: ParamMap): void => {
      this.currentPage = parseInt(route.get("page")) || 1;
      this.load();
      window.scroll({
        top: 0,
        left: 0,
      });
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  load() {
    this.apiClient.listArchive({pageSize: PAGE_SIZE, page: this.currentPage}).pipe(takeUntil(this.destroy$)).subscribe((res: RskArchiveList): void => {
      this.archive = res.items
      let totalPages: number = Math.ceil(res.resultCount / PAGE_SIZE);
      this.pages = Array(Math.min(totalPages, MAX_PAGINATION_LINKS)).fill(0).map((x, i) => i + 1);
      this.morePages = totalPages > MAX_PAGINATION_LINKS;
    })
  }
}
