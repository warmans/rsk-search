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
  set trivia(value: RskTrivia[]) {
    this._trivia = value;
    this.triviaTitles = value.map(v => this.simpleTrivia(v.description));
  }

  get trivia(): RskTrivia[] {
    return this._trivia;
  }

  private _trivia: RskTrivia[];

  triviaTitles: string[] = [];

  constructor() {
  }

  ngOnInit(): void {
  }

  simpleTrivia(original: string): string {
    const lines = (original.trim().split('\n') || []);
    return lines.length === 0 ? 'NA' : lines[0];
  }

}
