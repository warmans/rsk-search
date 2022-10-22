import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { CompFilter, CompOp } from 'src/app/lib/filter-dsl/filter';
import { Str } from 'src/app/lib/filter-dsl/value';

@Component({
  selector: 'app-search-bar-advanced',
  templateUrl: './search-bar-advanced.component.html',
  styleUrls: ['./search-bar-advanced.component.scss']
})
export class SearchBarAdvancedComponent implements OnInit {

  @Output()
  searchModifiersUpdated: EventEmitter<CompFilter[]> = new EventEmitter<CompFilter[]>();

  @Input()
  modifiers: CompFilter[] = [];

  speakers: string[] = ['ricky', 'steve', 'karl', 'claire', 'camfield'];

  publications: string[] = ['xfm', 'guide', 'bbc2', 'nme', 'fame'];

  types: string[] = ['chat', 'song', 'none'];

  constructor() {
  }

  ngOnInit(): void {
  }

  selectString(field: string, value: string) {
    const idx = this.findModifierIndex(field, value);
    if (idx > -1) {
      this.modifiers.splice(idx, 1);
    } else {
      this.modifiers.push(new CompFilter(field, CompOp.Eq, Str(value)));
    }
    this.searchModifiersUpdated.next(this.modifiers);
  }

  findModifierIndex(field: string, value: string): number {
    return this.modifiers.findIndex((f: CompFilter) => (f.field === field && f.value.v === value));
  }

}
