import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Selectable, ValueSource } from '../bar/bar.component';
import { Subscription } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

@Component({
  selector: 'app-value-list',
  templateUrl: './value-list.component.html',
  styleUrls: ['./value-list.component.scss'],
})
export class ValueListComponent implements OnInit {

  @Input()
  initialSelectedValues: string[] = [];

  @Input()
  valueSource: ValueSource;

  @Input()
  valueSourceFilters: Selectable[];

  @Input()
  valueSourcePaging: boolean;

  @Input()
  keyboardEvents: EventEmitter<string>;

  @Input()
  allowFreeInput: boolean;

  @Input()
  enableMultiselect: boolean;

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

  values: ValueListItem[] = [];

  page = 0;

  focusedValue = 0;

  private _filter: string;

  // use a observable for filter changes to allow debouncing
  private filterChange: EventEmitter<string> = new EventEmitter();

  private valueSub: Subscription;

  constructor() {
  }

  ngOnInit() {
    this.keyboardEvents.subscribe(key => this.onKeypress(key));

    this.fetchValuesFromSource(this._filter);
    this.filterChange.asObservable().pipe(debounceTime(100)).subscribe(value => {
      this.fetchValuesFromSource(value);
    });

    this.valuesSelected = (this.initialSelectedValues || []);
  }

  selectMulti() {
    if (this.enableMultiselect) {
      this.onValue.next(this.valuesSelected);
    }
  }

  select(value: ValueListItem) {
    this.toggleSelected(value.value);
    if (!this.enableMultiselect) {
      this.onValue.next([value.value]);
    }
  }

  onKeypress(key: string) {
    switch (key) {
      case 'ArrowDown':
        this.focusedValue = (this.focusedValue >= this.values.length - 1) ? 0 : this.focusedValue + 1;
        break;
      case 'ArrowUp':
        this.focusedValue = this.focusedValue === 0 ? this.values.length - 1 : this.focusedValue - 1;
        break;
      case 'Enter':
        if (this.values.length > 0) {
          this.select(this.values[this.focusedValue]);
        } else {
          if (this.allowFreeInput || !this.valueSource) {
            this.select({ value: this.filter, label: this.filter });
          }
        }
        break;
    }
  }

  fetchValuesFromSource(filter: string) {
    if (!this.valueSource) {
      return;
    }
    if (this.valueSub) {
      this.valueSub.unsubscribe();
    }
    this.valueSub = this.valueSource((this.valueSourceFilters || []), filter, this.page, this.pageSize).subscribe((values) => {
      this.values = values;
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

export interface ValueListItem {
  value: string;
  label?: string;
  helpText?: string;
}
