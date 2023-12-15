import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {SearchAPIClient} from "../../../../lib/api-client/services/search";
import {RskRandomQuote} from "../../../../lib/api-client/models";
import {takeUntil} from "rxjs/operators";

@Component({
  selector: 'app-random-quote',
  templateUrl: './random-quote.component.html',
  styleUrl: './random-quote.component.scss'
})
export class RandomQuoteComponent implements OnInit, OnDestroy {

  quote: RskRandomQuote;

  destroy$: EventEmitter<void> = new EventEmitter<void>();

  constructor(
    private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.apiClient.getRandomQuote().pipe(takeUntil(this.destroy$)).subscribe((res: RskRandomQuote) => {
      this.quote = res
    })
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

}
