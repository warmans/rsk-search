import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {SearchAPIClient} from "../../../../lib/api-client/services/search";
import {RskRandomQuote} from "../../../../lib/api-client/models";
import {takeUntil} from "rxjs/operators";

@Component({
    selector: 'app-random-quote',
    templateUrl: './random-quote.component.html',
    styleUrl: './random-quote.component.scss',
    standalone: false
})
export class RandomQuoteComponent implements OnInit, OnDestroy {

  quote: RskRandomQuote;
  quoteImage: string;

  destroy$: EventEmitter<void> = new EventEmitter<void>();

  constructor(
    private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.getQuote();
  }

  getQuote() {
    this.apiClient.getRandomQuote().pipe(takeUntil(this.destroy$)).subscribe((res: RskRandomQuote) => {
      this.quote = res
      switch (res.actor) {
        case "ricky":
          this.quoteImage = "ricky-seated.svg";
          break;
        case "steve":
          this.quoteImage = "steve-front.svg";
          break;
        case "karl":
          this.quoteImage = ["rgs-pointer.svg", "rgs-pointer-2.svg"][Math.floor(Math.random() * 2)];
          break;
        default:
          this.quoteImage = "rgs-pointer-2.svg";
          break;
      }
    })
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

}
