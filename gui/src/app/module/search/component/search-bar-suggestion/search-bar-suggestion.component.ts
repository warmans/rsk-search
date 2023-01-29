import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { Filter } from 'src/app/lib/filter-dsl/filter';

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
  dataFn: (prefix: string, filter: Filter) => Observable<string[]>;

  @Output()
  termSelected: EventEmitter<string> = new EventEmitter<string>();

  @Output()
  emitQuery: EventEmitter<void> = new EventEmitter<void>();

  values: string[] = [];

  selectedAutoCompleteRow: number = -1;

  loading: boolean = false;

  prefixChanged$: BehaviorSubject<string> = new BehaviorSubject<string>('');

  destroy$: Subject<void> = new Subject();

  constructor() {
    this.prefixChanged$.pipe(debounceTime(100), takeUntil(this.destroy$)).subscribe((termPrefix: string) => {
      this.loading = true;
      this.dataFn(termPrefix, this.termFilters).pipe(takeUntil(this.destroy$)).subscribe((res: string[]) => {
        this.values = res;
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
