import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { Filter } from 'src/app/lib/filter-dsl/filter';
import { RskPrediction } from 'src/app/lib/api-client/models';
import { highlightPrediction } from 'src/app/lib/util';


@Component({
  selector: 'app-search-bar-suggestion',
  templateUrl: './search-bar-suggestion.component.html',
  styleUrls: ['./search-bar-suggestion.component.scss']
})
export class SearchBarSuggestionComponent implements OnInit {

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
  termFilters: Filter;

  @Input()
  keyInput: Observable<KeyboardEvent> = new Observable<KeyboardEvent>();

  @Input()
  dataFn: (prefix: string, filter: Filter) => Observable<string[] | RskPrediction[]>;

  @Output()
  termSelected: EventEmitter<string> = new EventEmitter<string>();

  @Output()
  emitQuery: EventEmitter<void> = new EventEmitter<void>();

  values: string[] = [];

  highlightedValues: string[] = [];

  selectedAutoCompleteRow: number = -1;

  loading: boolean = false;

  prefixChanged$: BehaviorSubject<string> = new BehaviorSubject<string>('');

  destroy$: Subject<void> = new Subject();

  constructor() {
    this.prefixChanged$.pipe(debounceTime(100), takeUntil(this.destroy$)).subscribe((termPrefix: string) => {
      this.loading = true;
      this.dataFn(termPrefix.replace(/"/g, ''), this.termFilters).pipe(takeUntil(this.destroy$)).subscribe((res: string[] | RskPrediction[]) => {
        this.values = res.map((val: RskPrediction | string) => (typeof val === 'string') ? val : val.line)
        this.highlightedValues = res.map((val: RskPrediction | string) => (typeof val === 'string') ? val : highlightPrediction(val));
      }).add(() => {
        this.loading = false;
      });
    });
  }

  ngOnInit(): void {
    this.keyInput.pipe(takeUntil(this.destroy$)).subscribe((key) => {
      switch (key.key || key.code) {
        case 'ArrowDown':
          if (this.selectedAutoCompleteRow == this.values.length - 1) {
            this.selectedAutoCompleteRow = -1;
          } else {
            this.selectedAutoCompleteRow++;
          }
          break;
        case 'ArrowUp':
          if (this.selectedAutoCompleteRow > -1) {
            this.selectedAutoCompleteRow--;
          } else {
            this.selectedAutoCompleteRow = this.values ? this.values.length - 1 : -1;
          }
          break;
        case 'Enter':
        case 'Tab':
          if (this.selectedAutoCompleteRow > (this.values || []).length) {
            this.selectedAutoCompleteRow = -1;
          }
          if (this.selectedAutoCompleteRow === -1) {
            this.emitQuery.next();
            return;
          }
          this.selectTerm(this.values[this.selectedAutoCompleteRow]);
          return true;
        case 'Escape':
          break;
        default:
          return;
      }
    });
  }

  selectTerm(line: string) {
    this.termSelected.next(`${line}`);
  }
}
