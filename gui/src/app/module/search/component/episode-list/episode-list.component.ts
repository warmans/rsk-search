import { Component, EventEmitter, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RskShortTranscript, RskTranscriptList } from '../../../../lib/api-client/models';
import { SelectableConfig, SelectableKind } from '../../../shared/component/filterbar/bar/bar.component';
import { of } from 'rxjs';
import { FormControl } from '@angular/forms';

@Component({
  selector: 'app-episode-list',
  templateUrl: './episode-list.component.html',
  styleUrls: ['./episode-list.component.scss']
})
export class EpisodeListComponent implements OnInit {

  loading: boolean[] = [];

  transcriptList: RskShortTranscript[] = [];

  filteredTranscriptList: RskShortTranscript[] = [];

  showDownloadDialog: boolean = false;

  // simple filter
  searchInput: FormControl = new FormControl('');

  // complex filtering WIP
  filterBarConfig: SelectableConfig[] = [{
    kind: SelectableKind.FREETEXT,
    field: 'shortId',
    label: 'ID',
    helpText: 'The episode ID',
  }, {
    kind: SelectableKind.FREETEXT,
    field: 'bar',
    label: 'Boo',
    helpText: 'Another value.',
    valueSourcePaging: true,
    valueSource: (filters, query, page, pagesize) => {
      const values = [];
      for (let i = 0; i <= 100; i++) {
        values.push({ label: 'Foo ' + i, value: 'foo' + i });
      }
      return of(values.filter((v) => v.label.indexOf(query) > -1));
    },
    multiSelect: true,
  }
  ];

  private destroy$ = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.listEpisodes();

    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(100)).subscribe((val) => {
      if (val !== '') {
        this.filteredTranscriptList = (this.transcriptList || []).filter((t: RskShortTranscript) => {
          return t.shortId.toLowerCase().indexOf(val.toLowerCase()) > 0 || t.name.toLowerCase().indexOf(val.toLowerCase()) > 0;
        });
      } else {
        this.filteredTranscriptList = this.transcriptList;
      }
    });
  }

  listEpisodes() {
    this.loading.push(true);
    this.apiClient.listTranscripts().pipe(
      takeUntil(this.destroy$),
    ).subscribe((res: RskTranscriptList) => {
      this.filteredTranscriptList = this.transcriptList = res.episodes;
    }).add(() => {
      this.loading.pop();
    });
  }
}
