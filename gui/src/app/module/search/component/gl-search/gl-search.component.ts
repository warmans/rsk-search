import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { MetaService } from '../../../core/service/meta/meta.service';
import { RsksearchFieldMeta } from '../../../../lib/api-client/models';
import { first } from 'rxjs/operators';
import { And, CompOp, Filter, NewCompFilter } from '../../../../lib/filter-dsl/filter';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { PrintPlainText } from '../../../../lib/filter-dsl/printer';
import { Str } from '../../../../lib/filter-dsl/value';

@Component({
  selector: 'app-gl-search',
  templateUrl: './gl-search.component.html',
  styleUrls: ['./gl-search.component.scss']
})
export class GlSearchComponent implements OnInit {

  @Output()
  queryUpdated: EventEmitter<string> = new EventEmitter<string>();

  searchForm = new FormGroup({
    term: new FormControl(null, [Validators.maxLength(1024)]),
  });

  fieldMeta: RsksearchFieldMeta[] = [];

  activeFilters: RsksearchFieldMeta[] = [];

  constructor(meta: MetaService) {

    // todo: parse query from route

    meta.getMeta().pipe(first()).subscribe((m) => {
      this.fieldMeta = m.fields;
    });
  }

  ngOnInit(): void {
  }

  addFilter(field: string) {
    this.activeFilters.push(this.fieldMeta.find((v) => v.name === field));
  }

  removeFilter(idx: number) {
    this.activeFilters.splice(idx, 1);
  }

  emitQuery() {
    let term = this.searchForm.get('term').value;
    if (term === null) {
      return;
    }

    // split terms into quoted and non-quoted strings
    const quotedStrings = term.match(/"([^\\"]+)"/g) || [];
    const exactMatchTerms: string[] = [];
    quotedStrings.forEach((v: string) => {
      term = term.replace(v, '', -1);
      exactMatchTerms.push(v.replace(/["]/g, ''));
    });
    if (quotedStrings.length === 0 && term.trim().length === 0) {
      return;
    }
    term = term.trim()

    console.log("ferm", term, exactMatchTerms);

    // assemble them into a query
    let query: Filter = term.length === 0 ?
      NewCompFilter('content', CompOp.Like, Str(exactMatchTerms[0])) :
      NewCompFilter('content', CompOp.Like, Str(term));

    exactMatchTerms.forEach((v: string, k: number) => {
      v = v.trim();
      if (term.length === 0 && k == 0) {
        return;
      }
      query = And(query, NewCompFilter('content', CompOp.Eq, Str(v)));
    });
    this.queryUpdated.emit(PrintPlainText(query));
  }
}
