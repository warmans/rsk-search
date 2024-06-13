import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {debounceTime, takeUntil} from 'rxjs/operators';
import {RskPublicationType, RskShortTranscript, RskTranscriptList} from 'src/app/lib/api-client/models';
import {UntypedFormControl} from '@angular/forms';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
import {KeyValue} from "@angular/common";

@Component({
  selector: 'app-episode-list',
  templateUrl: './episode-list.component.html',
  styleUrls: ['./episode-list.component.scss'],
})
export class EpisodeListComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  transcriptList: RskShortTranscript[] = [];

  publicationCategories: { [index: string]: RskPublicationType } = {
    "Radio": RskPublicationType.PUBLICATION_TYPE_RADIO,
    "Podcast": RskPublicationType.PUBLICATION_TYPE_PODCAST,
    "Promo": RskPublicationType.PUBLICATION_TYPE_PROMO,
    //"TV": RskPublicationType.PUBLICATION_TYPE_TV,
    "Other": RskPublicationType.PUBLICATION_TYPE_OTHER,
  };
  filteredTranscriptList: RskShortTranscript[] = [];

  showDownloadDialog: boolean = false;

  searchInput: UntypedFormControl = new UntypedFormControl('');

  private _activePublicationType: RskPublicationType = RskPublicationType.PUBLICATION_TYPE_RADIO;


  get activePublicationType(): RskPublicationType {
    return this._activePublicationType;
  }

  set activePublicationType(value: RskPublicationType) {
    this._activePublicationType = value;
    this.resetEpisodeList();
  }

  private destroy$ = new EventEmitter<void>();

  constructor(private apiClient: SearchAPIClient, private router: Router, route: ActivatedRoute) {
    route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      this.activePublicationType = params.get('publication_type') as RskPublicationType || RskPublicationType.PUBLICATION_TYPE_RADIO;
    });
  }

  ngOnInit(): void {
    this.listEpisodes();
    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(100)).subscribe((val) => {
      val = val.trim().toLowerCase();
      if (val !== '') {
        this.filteredTranscriptList = this.transcriptList.filter((t: RskShortTranscript) => {
          return t.shortId.toLowerCase().indexOf(val) > -1 || t.name.toLowerCase().indexOf(val) > -1;
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
    return (this.transcriptList?.
    filter((t => t.publicationType === this.activePublicationType)) || []).
    sort((v, k): number => {
      if (v.releaseDate) {
        return new Date(v.releaseDate).getTime() > new Date(k.releaseDate).getTime() ? 1 : -1
      }
      return v.series * 100 + v.episode > k.series * 100 + k.episode ? 1 : -1;

    });
  }

  resetEpisodeList() {
    this.searchInput.setValue('');
    this.filteredTranscriptList = this.activePublicationTranscripts();
  }

  loadPublicationTab(tab: RskPublicationType) {
    this.searchInput.setValue("");
    this.router.navigate(['/search'], {queryParams: {'publication_type': tab}}).finally();
  }

  originalOrder = (a: KeyValue<string,string>, b: KeyValue<string,string>): number => {
    return 0;
  }
}
