import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-search-bar-mention',
  templateUrl: './search-bar-mention.component.html',
  styleUrls: ['./search-bar-mention.component.scss']
})
export class SearchBarMentionComponent implements OnInit {

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

  @Input()
  dataFn: (prefix: string) => Observable<string[]>;

  @Input()
  termTypeIdentifier: string;

  @Output()
  termSelected: EventEmitter<string> = new EventEmitter<string>();

  values: string[] = [];

  selectedAutoCompleteRow: number = -1;

  loading: boolean = false;

  prefixChanged$: Subject<string> = new Subject<string>();

  destroy$: Subject<void> = new Subject();

  constructor() {
    this.prefixChanged$.pipe(debounceTime(100), takeUntil(this.destroy$)).subscribe((termPrefix: string) => {
      this.loading = true;
      let replacer = new RegExp(`[${this.termTypeIdentifier}]`, 'g')
      this.dataFn(termPrefix.replace(replacer, '')).pipe(takeUntil(this.destroy$)).subscribe((res) => {
        this.values = res;
      });
    });
  }

  ngOnInit(): void {
    this.keyInput.pipe(takeUntil(this.destroy$)).subscribe((key) => {
      switch (key.code) {
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
          if (this.selectedAutoCompleteRow > (this.values || []).length) {
            this.selectedAutoCompleteRow = -1;
          }
          if (this.selectedAutoCompleteRow === -1) {
            this.selectTerm(this.termPrefix);
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
    this.termSelected.next(`${this.termTypeIdentifier}${line}`);
  }
}
