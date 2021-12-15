import { Component, EventEmitter, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { RskShortTranscript, RskTranscriptList } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-episode-list',
  templateUrl: './episode-list.component.html',
  styleUrls: ['./episode-list.component.scss']
})
export class EpisodeListComponent implements OnInit {

  loading: boolean[] = [];

  transcriptList: RskShortTranscript[] = [];

  showSeries: number[] = [1, 2, 3, 4];

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>()

  constructor(private apiClient: SearchAPIClient) { }

  ngOnInit(): void {
    this.listEpisodes();
  }

  listEpisodes() {
    this.loading.push(true);
    this.apiClient.listTranscripts().pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe((res: RskTranscriptList) => {
      this.transcriptList = res.episodes;
    }).add(() => {
      this.loading.pop();
    });
  }

}
