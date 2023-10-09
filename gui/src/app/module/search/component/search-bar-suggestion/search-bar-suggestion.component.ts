import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {BehaviorSubject, Observable, Subject} from 'rxjs';
import {debounceTime, takeUntil} from 'rxjs/operators';
import {CompOp, Filter} from 'src/app/lib/filter-dsl/filter';
import {RskPrediction} from 'src/app/lib/api-client/models';
import {highlightPrediction} from 'src/app/lib/util';
import {Term} from 'src/app/lib/search-parser/parser';

export enum SuggestionType {
  Dialog = "dialog",
  Any = "any"
}

export interface Suggestion {
  type: SuggestionType;
  term: string;
  epid?: string;
  pos?: number;
  actor?: string;
}

@Component({
  selector: 'app-search-bar-suggestion',
  templateUrl: './search-bar-suggestion.component.html',
  styleUrls: ['./search-bar-suggestion.component.scss']
})
export class SearchBarSuggestionComponent implements OnInit {

  @Input()
  set term(value: Term) {
    this._term = value;
    this.termChanged$.next(value);
  }

  get term(): Term {
    return this._term;
  }

  private _term: Term;

  @Input()
  termFilters: Filter;

  @Input()
  keyInput: Observable<KeyboardEvent> = new Observable<KeyboardEvent>();

  @Input()
  dataFn: (prefix: string, filter: Filter, exact: boolean) => Observable<string[] | RskPrediction[]>;

  @Output()
  termSelected: EventEmitter<Suggestion> = new EventEmitter<Suggestion>();

  @Output()
  predicationSelected: EventEmitter<RskPrediction> = new EventEmitter<RskPrediction>();

  @Output()
  emitQuery: EventEmitter<void> = new EventEmitter<void>();

  values: Suggestion[] = [];

  highlightedValues: string[] = [];

  selectedAutoCompleteRow: number = -1;

  loading: boolean = false;

  termChanged$: BehaviorSubject<Term> = new BehaviorSubject<Term>(undefined);

  destroy$: Subject<void> = new Subject();

  constructor() {

    this.termChanged$.pipe(debounceTime(100), takeUntil(this.destroy$))
      .subscribe((term: Term) => {
        this.loading = true;
        this.dataFn(
          term.value.replace(/"/g, ''),
          this.termFilters,
          term.op === CompOp.Eq)
          .pipe(takeUntil(this.destroy$)).subscribe((res: string[] | RskPrediction[]) => {
          this.values = res.map((val: RskPrediction | string): Suggestion => {
            return (typeof val === 'string') ? {type: SuggestionType.Any, term: val} : {
              type: SuggestionType.Dialog,
              term: val.line,
              epid: val.epid,
              pos: val.pos,
              actor: val.actor,
            }
          });
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

  selectTerm(line: Suggestion) {
    this.termSelected.next({
      type: line.type,
      term: (/\s/).test(line.term) ? `"${line.term}"` : `${line.term}`,
      epid: line.epid,
      pos: line.pos
    });
  }
}
