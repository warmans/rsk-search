import { Component, Input, OnInit } from '@angular/core';
import { RskShortTranscript } from '../../../../lib/api-client/models';
import { AudioService } from '../../../core/service/audio/audio.service';

@Component({
  selector: 'app-episode-summary',
  templateUrl: './episode-summary.component.html',
  styleUrls: ['./episode-summary.component.scss']
})
export class EpisodeSummaryComponent implements OnInit {

  @Input()
  set episode(value: RskShortTranscript) {
    this._episode = value;
    this.episodeImage = value?.metadata["cover_art_url"] ? value?.metadata["cover_art_url"] : `/assets/cover/${value.publication}-s${value.series}.jpg`
  }
  get episode(): RskShortTranscript {
    return this._episode;
  }
  private _episode: RskShortTranscript;

  episodeImage: string;

  constructor(private audioService: AudioService) {
  }

  ngOnInit(): void {
  }

  playEpisode(episode: RskShortTranscript) {
    this.audioService.setAudioSrc(episode.shortId, episode.audioUri);
    this.audioService.playAudio();
  }
}
