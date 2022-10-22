import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RskPrediction, RskSearchTermPredictions } from 'src/app/lib/api-client/models';
import { Observable, Subject } from 'rxjs';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { highlightPrediction } from 'src/app/lib/util';

@Component({
  selector: 'app-search-bar-autocomplete',
  templateUrl: './search-bar-autocomplete.component.html',
  styleUrls: ['./search-bar-autocomplete.component.scss']
})
export class SearchBarAutocompleteComponent implements OnInit {

  @Input()
  set termPrefix(value: string) {
    this._termPrefix = value;
    this.prefixChanged$.next(value);
  }

  get termPrefix(): string {
    return this._termPrefix;
  }

  private _termPrefix: string;

  @Input()
  keyInput: Observable<KeyboardEvent> = new Observable<KeyboardEvent>();

  @Output()
  temSelected: EventEmitter<string> = new EventEmitter<string>();

  autocompleteVals: RskPrediction[] = [];

  highlightedAutocompleteVals: string[] = [];

  selectedAutoCompleteRow: number = -1;

  prefixChanged$: Subject<string> = new Subject<string>();

  destroy$: Subject<void> = new Subject();

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.prefixChanged$.pipe(debounceTime(100), takeUntil(this.destroy$)).subscribe((termPrefix: string) => {
      this.apiClient.predictSearchTerm({ prefix: termPrefix, maxPredictions: 10 })
        .pipe(takeUntil(this.destroy$))
        .subscribe((res: RskSearchTermPredictions) => {
          this.autocompleteVals = res.predictions;
          this.highlightedAutocompleteVals = res.predictions.map((val: RskPrediction) => highlightPrediction(val));
        });
    });
  }

  selectTerm(line: string) {
    this.temSelected.next(line);
  }

}
