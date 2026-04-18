import { Component, OnInit } from '@angular/core';
import { AudioPlayerV2Component } from '../audio-player-v2/audio-player-v2.component';

@Component({
  selector: 'app-audio-player-fixed',
  templateUrl: './audio-player-fixed.component.html',
  styleUrls: ['./audio-player-fixed.component.scss'],
  imports: [AudioPlayerV2Component],
})
export class AudioPlayerFixedComponent implements OnInit {
  constructor() {}

  ngOnInit(): void {}
}
