import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-transcript-change',
  templateUrl: './transcript-change.component.html',
  styleUrls: ['./transcript-change.component.scss']
})
export class TranscriptChangeComponent implements OnInit {

  episodeID: string;

  constructor() { }

  ngOnInit(): void {
  }

}
