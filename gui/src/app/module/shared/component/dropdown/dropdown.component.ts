import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, Renderer2 } from '@angular/core';
import { Observable, Subscription } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-dropdown',
  templateUrl: './dropdown.component.html',
  styleUrls: ['./dropdown.component.scss'],
})
export class DropdownComponent implements OnInit, OnDestroy {

  private _valueSource: (fieldName: string, prefix: string) => Observable<string[]>;

  @Input()
  set valueSource(f: (fieldName: string, prefix: string) => Observable<string[]>) {
    if (this.valueSub) {
      this.valueSub.unsubscribe();
    }
    this._valueSource = f;
    if (this._filter !== null) {
      this.fetchValuesFromSource(this._filter);
    }
  }

  get valueSource(): (fieldName: string, prefix: string) => Observable<string[]> {
    return this._valueSource;
  }

  @Input()
  fieldName: string;

  @Input()
  keyboardEvents: EventEmitter<KeyboardEvent> = new EventEmitter<KeyboardEvent>();

  @Input()
  enableMultiselect: boolean = false;

  private _filter: string;

  @Input()
  set filter(f: string) {
    this._filter = f;
    this.page = 0;
    this.filterChange.next(f);
  }

  get filter(): string {
    return this._filter;
  }

  @Input()
  pageSize: 10;

  @Output()
  onValue: EventEmitter<string[]> = new EventEmitter();

  valuesSelected: string[] = [];

  values: string[] = [];

  page = 0;

  focusedValue = 0;

  loading: boolean = false;

  private $destroy: EventEmitter<any> = new EventEmitter<any>();

  // use a observable for filter changes to allow debouncing
  private filterChange: EventEmitter<string> = new EventEmitter();

  private valueSub: Subscription;

  constructor(private renderer: Renderer2) {
  }

  ngOnInit() {
    this.keyboardEvents.pipe(takeUntil(this.$destroy)).subscribe(key => this.onKeypress(key));

    this.fetchValuesFromSource(this._filter);
    this.filterChange.asObservable().pipe(takeUntil(this.$destroy), debounceTime(100)).subscribe(value => {
      this.fetchValuesFromSource(value);
    });
  }

  ngOnDestroy(): void {
    this.$destroy.emit(true);
    this.$destroy.complete();
  }

  selectMulti() {
    if (this.enableMultiselect) {
      this.onValue.next(this.valuesSelected);
    }
  }

  select(value: string) {
    if (value === undefined) {
      return;
    }
    this.toggleSelected(value);
    if (!this.enableMultiselect) {
      this.onValue.next([value]);
    }
  }

  onKeypress(key: KeyboardEvent) {
    switch (key.code) {
      case 'ArrowDown':
        this.focusedValue = (this.focusedValue >= this.values.length - 1) ? 0 : this.focusedValue + 1;
        break;
      case 'ArrowUp':
        this.focusedValue = this.focusedValue === 0 ? this.values.length - 1 : this.focusedValue - 1;
        break;
      case 'Enter':
        if (this.values.length > 0) {
          this.select(this.values[this.focusedValue]);
        }
        break;
    }
  }

  fetchValuesFromSource(filter: string) {
    if (!this._valueSource) {
      return;
    }
    this.loading = true;
    this.valueSub = this._valueSource(this.fieldName || '', filter || '').pipe(takeUntil(this.$destroy)).subscribe((vals: string[]) => {
      this.values = vals;
      this.loading = false;
    });
  }

  pageBack() {
    if (this.page > 0) {
      this.page = this.page - 1;
      this.fetchValuesFromSource(this._filter);
    }
  }

  pageForward() {
    // if there are less than the page size values we're probably on the last page.
    // or just unlucky.
    if (this.values.length === this.pageSize) {
      this.page++;
      this.fetchValuesFromSource(this._filter);
    }
  }

  isSelected(value: string) {
    return (this.valuesSelected.indexOf(value) !== -1);
  }

  toggleSelected(value: string) {
    if (this.isSelected(value)) {
      this.valuesSelected.splice(this.valuesSelected.indexOf(value), 1);
      return;
    }
    this.valuesSelected.push(value);
  }
}
