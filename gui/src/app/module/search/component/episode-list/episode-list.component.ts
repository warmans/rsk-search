import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {debounceTime, distinctUntilChanged, takeUntil} from 'rxjs/operators';
import {RskPublicationType, RskShortTranscript, RskTranscriptList} from 'src/app/lib/api-client/models';
import {UntypedFormControl} from '@angular/forms';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
import {KeyValue} from "@angular/common";
import {Eq} from "../../../../lib/filter-dsl/filter";
import {Str} from "../../../../lib/filter-dsl/value";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-episode-list',
  templateUrl: './episode-list.component.html',
  styleUrls: ['./episode-list.component.scss'],
})
export class EpisodeListComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  transcriptList: RskShortTranscript[] = [];

  subSections: { [index: string]: Array<string> } = {};

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
    if (this._activePublicationType !== value) {
      this._activePublicationType = value;
      this.resetEpisodeList(true);
    }
  }

  private _activeSubSection: string;

  set activeSubSection(value: string) {
    if (this._activeSubSection !== value) {
      this._activeSubSection = value;
      this.resetEpisodeList();
    }
  }

  get activeSubSection(): string {
    return this._activeSubSection;
  }

  private destroy$ = new EventEmitter<void>();

  constructor(private apiClient: SearchAPIClient, private router: Router, private route: ActivatedRoute) {
    // todo: bug - list is loaded twice if the pace is refreshed with something in the URL
    // this is a quick fix
    if (!this.route.snapshot.queryParamMap.get('publication_type')){
      this.listEpisodes();
    }
    route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      this.activePublicationType = params.get('publication_type') as RskPublicationType ?? RskPublicationType.PUBLICATION_TYPE_RADIO;
      this.activeSubSection = params.get('subsection') ?? (this.activePublicationType === RskPublicationType.PUBLICATION_TYPE_RADIO ? "xfm-S1" : undefined);
    });
  }

  ngOnInit(): void {

    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), distinctUntilChanged(), debounceTime(100)).subscribe((val) => {
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

  listEpisodes(): Subscription {
    this.transcriptList = [];
    this.filteredTranscriptList = [];

    this.loading.push(true);
    return this.apiClient.listTranscripts({filter: Eq("publication_type", Str(this.mapPublicationType(this._activePublicationType))).print()}).pipe(
      takeUntil(this.destroy$),
    ).subscribe((res: RskTranscriptList) => {
      this.transcriptList = res.episodes;
      this.updateFilteredTranscriptList();
      this.identifySubsections();
    }).add(() => {
      this.loading.pop();
    });
  }

  updateFilteredTranscriptList() {

    this.filteredTranscriptList = (this.transcriptList?.filter((t => {
      return t.publicationType === this.activePublicationType && (!this.activeSubSection || `${t.publication}-S${t.series}` === this.activeSubSection)
    })) || []).sort((v, k): number => {
      if (v.releaseDate) {
        return new Date(v.releaseDate).getTime() > new Date(k.releaseDate).getTime() ? 1 : -1
      }
      return v.series * 100 + v.episode > k.series * 100 + k.episode ? 1 : -1;
    });
  }

  identifySubsections() {
    let subsections: { [index: string]: Array<{ sub: string, publishDate: Date }> } = {};
    this.transcriptList.forEach((ts) => {
      const sub = `${ts.publication}-S${ts.series}`;
      if (subsections[ts.publicationType] == null) {
        subsections[ts.publicationType] = []
      }
      if (!subsections[ts.publicationType].find((v) => v.sub === sub)) {
        subsections[ts.publicationType].push({sub, publishDate: new Date(ts.releaseDate)});
      }
    })

    this.subSections = {};
    Object.keys(subsections).forEach((k) => {
      this.subSections[k] = subsections[k].sort((a, b) => a.publishDate > b.publishDate ? 1 : -1).map(s => s.sub);
    })
  }

  resetEpisodeList(reload: boolean = false) {
    if (reload) {
      this.listEpisodes().add(() => {
        this.searchInput.setValue('', {emitEvent: false});
        this.updateFilteredTranscriptList();
      })
    } else {
      this.searchInput.setValue('', {emitEvent: false});
      this.updateFilteredTranscriptList();
    }
  }

  loadPublicationTab(tab: RskPublicationType) {
    this.searchInput.setValue("");
    this.router.navigate(['/search'], {queryParams: {'publication_type': tab}});
  }

  loadSubsection(sub: string) {
    this.searchInput.setValue("");
    this.router.navigate(['/search'], {queryParams: {'subsection': sub}, queryParamsHandling: 'merge'});
  }

  originalOrder = (a: KeyValue<string, string>, b: KeyValue<string, string>): number => {
    return 0;
  }

  mapPublicationType(type: RskPublicationType): string {
    switch (type) {
      case RskPublicationType.PUBLICATION_TYPE_RADIO:
        return 'radio';
      case RskPublicationType.PUBLICATION_TYPE_OTHER:
        return 'other';
      case RskPublicationType.PUBLICATION_TYPE_PODCAST:
        return 'podcast';
      case RskPublicationType.PUBLICATION_TYPE_PROMO:
        return 'promo';
      case RskPublicationType.PUBLICATION_TYPE_TV:
        return 'tv';
      default:
        return 'unknown';
    }
  }

}
