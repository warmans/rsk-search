import {Component, Input} from '@angular/core';

@Component({
  selector: 'app-star-rating',
  standalone: false,
  templateUrl: './star-rating.component.html',
  styleUrl: './star-rating.component.scss'
})
export class StarRatingComponent {

  @Input()
  set score(value: number) {
    this._score = value || 0;
    if (this._score > 0) {
      this.scoreArr = Array(Math.floor(this._score)).fill('bi-star-fill');
      if (this._score > Math.floor(this._score) + 0.3) {
        this.scoreArr.push('bi-star-half')
      }
    }
    if (this.scoreArr.length < 5) {
      this.scoreArr.push(...Array(5 - Math.floor(this.scoreArr.length)).fill('bi-star'))
    }
  }
  get score(): number {
    return this._score || 0;
  }
  private _score: number = 0;

  scoreArr: Array<'bi-star-fill' | 'bi-star-half'> = [];

}
