import { Component, Input, OnInit } from '@angular/core';
import { RskTrivia } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-trivia',
  templateUrl: './trivia.component.html',
  styleUrls: ['./trivia.component.scss']
})
export class TriviaComponent implements OnInit {

  @Input()
  episodeID: string;

  @Input()
  trivia: RskTrivia[];

  constructor() {
  }

  ngOnInit(): void {
  }

}
