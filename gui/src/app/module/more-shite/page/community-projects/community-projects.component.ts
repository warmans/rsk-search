import {Component, OnDestroy, OnInit} from '@angular/core';
import {CommunityAPIClient} from "../../../../lib/api-client/services/community";
import {ActivatedRoute, ParamMap, Router} from "@angular/router";
import {UntypedFormControl} from "@angular/forms";
import {Subject} from "rxjs";
import {debounceTime, takeUntil} from "rxjs/operators";
import {Like} from "../../../../lib/filter-dsl/filter";
import {Str} from "../../../../lib/filter-dsl/value";
import {RskCommunityProject, RskCommunityProjectList} from "../../../../lib/api-client/models";

const PAGE_SIZE = 25;
const MAX_PAGINATION_LINKS = 10;

@Component({
  selector: 'app-community-projects',
  standalone: false,
  templateUrl: './community-projects.component.html',
  styleUrl: './community-projects.component.scss'
})
export class CommunityProjectsComponent implements OnInit, OnDestroy {

  searchInput: UntypedFormControl = new UntypedFormControl('');

  currentPage: number;
  pages: number[];
  morePages: boolean;
  maxPages: number = MAX_PAGINATION_LINKS;

  projects: Array<RskCommunityProject & {expand?: boolean}> = [];

  private destroy$: Subject<void> = new Subject<void>();

  constructor(private apiClient: CommunityAPIClient, private route: ActivatedRoute, private router: Router) {
  }

  ngOnInit(): void {
    this.route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((route: ParamMap): void => {
      this.currentPage = parseInt(route.get("page")) || 1;
      this.searchInput.setValue(route.get("q") || "", {emitEvent: false});
      this.updateResults();
    });
    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(500)).subscribe((searchTerm) => {
      this.router.navigate(['/more-shite', "community-projects"], {queryParams: {q: searchTerm.trim()}});
      this.currentPage = 1;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private updateResults() {
    this.apiClient.listCommunityProjects({
      filter: this.searchInput.value.trim() === "" ? "" : Like("name", Str(this.searchInput.value.trim())).print(),
      page: this.currentPage,
      pageSize: PAGE_SIZE,
    }).pipe(takeUntil(this.destroy$)).subscribe((res: RskCommunityProjectList): void => {
      this.projects = res.projects;
      let totalPages: number = Math.ceil(res.resultCount / PAGE_SIZE);
      this.pages = Array(Math.min(totalPages, MAX_PAGINATION_LINKS)).fill(0).map((x, i) => i + 1);
      this.morePages = totalPages > MAX_PAGINATION_LINKS;
    })
  }
}
