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
  episode: RskShortTranscript;

  constructor(private audioService: AudioService) {
  }

  ngOnInit(): void {
  }

  playEpisode(episode: RskShortTranscript) {
    this.audioService.setAudioSrc(episode.shortId, episode.audioUri);
    this.audioService.playAudio();
  }
}
