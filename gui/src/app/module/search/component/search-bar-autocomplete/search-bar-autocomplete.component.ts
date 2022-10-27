import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RskPrediction, RskSearchTermPredictions } from 'src/app/lib/api-client/models';
import { Observable, Subject } from 'rxjs';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { highlightPrediction } from 'src/app/lib/util';
import { And, CompFilter, Filter } from 'src/app/lib/filter-dsl/filter';
import { PrintPlainText } from 'src/app/lib/filter-dsl/printer';

@Component({
  selector: 'app-search-bar-autocomplete',
  templateUrl: './search-bar-autocomplete.component.html',
  styleUrls: ['./search-bar-autocomplete.component.scss']
})
export class SearchBarAutocompleteComponent implements OnInit {

  @Input()
  searchModifiers: CompFilter[] = [];

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

  loading: boolean = false;

  prefixChanged$: Subject<string> = new Subject<string>();

  destroy$: Subject<void> = new Subject();

  constructor(private apiClient: SearchAPIClient) {
    this.prefixChanged$.pipe(debounceTime(100), takeUntil(this.destroy$)).subscribe((termPrefix: string) => {

      let query: Filter = null;
      if ((this.searchModifiers || []).length > 0) {
        this.searchModifiers.forEach((comp: CompFilter) => {
          if (query == null) {
            query = comp;
            return;
          }
          query = And(query, comp);
        });
      }

      this.loading = true;
      this.apiClient.predictSearchTerm({ prefix: termPrefix, maxPredictions: 10, query: query ? PrintPlainText(query) : '' })
        .pipe(takeUntil(this.destroy$))
        .subscribe((res: RskSearchTermPredictions) => {
          this.autocompleteVals = res.predictions;
          this.highlightedAutocompleteVals = res.predictions.map((val: RskPrediction) => highlightPrediction(val));
        }).add(() => {
        this.loading = false;
      });
    });
  }

  ngOnInit(): void {
    this.keyInput.pipe(takeUntil(this.destroy$)).subscribe((key) => {
      switch (key.code) {
        case 'ArrowDown':
          if (this.selectedAutoCompleteRow == this.autocompleteVals.length - 1) {
            this.selectedAutoCompleteRow = -1;
          } else {
            this.selectedAutoCompleteRow++;
          }
          break;
        case 'ArrowUp':
          if (this.selectedAutoCompleteRow > -1) {
            this.selectedAutoCompleteRow--;
          } else {
            this.selectedAutoCompleteRow = this.autocompleteVals ? this.autocompleteVals.length - 1 : -1;
          }
          break;
        case 'Enter':
          if (this.selectedAutoCompleteRow > (this.autocompleteVals || []).length) {
            this.selectedAutoCompleteRow = -1;
          }
          if (this.selectedAutoCompleteRow === -1) {
            this.selectTerm(this.termPrefix, false);
            return;
          }
          this.selectTerm(this.autocompleteVals[this.selectedAutoCompleteRow].line, true );
          return true;
        case 'Escape':
          break;
        default:
          return true;
      }
    });
  }

  selectTerm(line: string, exact: boolean) {
    this.temSelected.next(exact ? `"${line}"` : line);
  }

}
