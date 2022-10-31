import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RskShortTranscript, RskTranscriptList } from 'src/app/lib/api-client/models';
import { FormControl } from '@angular/forms';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';

type tabState = 'xfm' | 'guide' | 'special' | 'other' | 'karl' | 'preview';

@Component({
  selector: 'app-episode-list',
  templateUrl: './episode-list.component.html',
  styleUrls: ['./episode-list.component.scss'],
})
export class EpisodeListComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  transcriptList: RskShortTranscript[] = [];

  filteredTranscriptList: RskShortTranscript[] = [];

  showDownloadDialog: boolean = false;

  searchInput: FormControl = new FormControl('');

  private _activePublication: tabState = 'xfm';

  get activePublication(): tabState {
    return this._activePublication;
  }

  set activePublication(value: tabState) {
    this._activePublication = value;
    this.resetEpisodeList();
  }

  private destroy$ = new EventEmitter<void>();

  constructor(private apiClient: SearchAPIClient, private router: Router, route: ActivatedRoute) {
    route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      this.activePublication = params.get('publication') as tabState || 'xfm';
    });
  }

  ngOnInit(): void {
    this.listEpisodes();
    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(100)).subscribe((val) => {
      if (val !== '') {
        this.filteredTranscriptList = this.activePublicationTranscripts().filter((t: RskShortTranscript) => {
          return t.shortId.toLowerCase().indexOf(val.toLowerCase()) > 0 || t.name.toLowerCase().indexOf(val.toLowerCase()) > 0;
        });
      } else {
        this.resetEpisodeList();
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  listEpisodes() {
    this.loading.push(true);
    this.apiClient.listTranscripts().pipe(
      takeUntil(this.destroy$),
    ).subscribe((res: RskTranscriptList) => {
      this.transcriptList = res.episodes;
      this.filteredTranscriptList = this.activePublicationTranscripts();
    }).add(() => {
      this.loading.pop();
    });
  }

  activePublicationTranscripts(): RskShortTranscript[] {
    if (this.activePublication === 'special') {
      return this.transcriptList?.filter(t => t.special) || [];
    }
    return this.transcriptList?.filter((t => !t.special && t.publication === this.activePublication)) || [];
  }

  resetEpisodeList() {
    this.searchInput.setValue('');
    this.filteredTranscriptList = this.activePublicationTranscripts();
  }

  loadPublicationTab(tab: tabState) {
    this.router.navigate(['/search'], { queryParams: { 'publication': tab } }).finally();
  }
}
