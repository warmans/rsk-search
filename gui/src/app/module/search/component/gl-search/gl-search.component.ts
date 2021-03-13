import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { MetaService } from '../../../core/service/meta/meta.service';
import { FieldMetaKind, RsksearchFieldMeta } from '../../../../lib/api-client/models';
import { distinctUntilChanged, first } from 'rxjs/operators';
import { And, BoolFilter, CompFilter, CompOp, Filter, NewCompFilter, Visitor } from '../../../../lib/filter-dsl/filter';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { PrintPlainText } from '../../../../lib/filter-dsl/printer';
import { Str, ValueFromFieldMeta } from '../../../../lib/filter-dsl/value';
import { ActivatedRoute } from '@angular/router';
import { ParseAST } from '../../../../lib/filter-dsl/ast';
import { SearchFilter } from '../gl-search-filter/gl-search-filter.component';

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

  activeFilters: SearchFilter[] = [];

  filterOverlayOpen = false;

  fieldKinds = FieldMetaKind;

  constructor(meta: MetaService, route: ActivatedRoute) {
    meta.getMeta().pipe(first()).subscribe((m) => {
      this.fieldMeta = m.fields;
      route.queryParamMap.pipe(distinctUntilChanged()).subscribe((params) => {
        if (params.get('q') === null || params.get('q').trim() === '') {
          this.resetFilters();
          return;
        }
        this.parseQuery(params.get('q'));
      });
    });
  }

  ngOnInit(): void {
  }

  resetFilters() {
    this.activeFilters = [];
    this.searchForm.get('term').setValue('');
  }

  addFilter(field: string) {
    this.filterOverlayOpen = false;
    let meta = this.getMetaForField(field);
    if (!meta) {
      console.error('unknown field', field);
      return;
    }
    this.activeFilters.unshift(new SearchFilter(meta));
  }

  getMetaForField(field: string): RsksearchFieldMeta {
    return this.fieldMeta.find((v) => v.name === field);
  }

  removeFilter(idx: number) {
    this.activeFilters.splice(idx, 1);
  }

  emitQuery() {
    let term = this.searchForm.get('term').value || '';

    // group terms by exact/non-exact and convert them into a single filter statement.
    let query: Filter = null;
    this.parseSearchTerm(term).forEach((v: searchTerm) => {
      let comp = NewCompFilter('content', v.exact ? CompOp.Eq : CompOp.Like, Str(v.value));
      if (query == null) {
        query = comp;
        return;
      }
      query = And(query, comp);
    });
    this.activeFilters.forEach((f: SearchFilter) => {
      let comp = NewCompFilter(f.meta.name, f.operator, ValueFromFieldMeta(f.meta, f.value));
      if (query == null) {
        query = comp;
        return;
      }
      query = And(query, comp);
    });
    if (query !== null) {
      this.queryUpdated.emit(PrintPlainText(query));
    } else {
      this.queryUpdated.emit("");
    }
  }

  parseSearchTerm(term: string): searchTerm[] {
    let searchTerms: searchTerm[] = [];
    let currentTerm = '';
    let currentTermExact = false;

    for (let i = 0; i < term.length; i++) {
      if (term[i] === '"') {
        // new exact term
        if (currentTerm.length === 0) {
          currentTermExact = true;
          continue;
        }
        // end of exact term
        if (currentTermExact) {
          if (currentTerm.trim().length > 0) {
            searchTerms.push(new searchTerm(currentTerm.trim(), true));
          }
          currentTermExact = false;
          currentTerm = '';
          continue;
        }
        // start of exact term (where previous term was vague)
        if (currentTerm.trim().length > 0) {
          searchTerms.push(new searchTerm(currentTerm.trim(), false));
        }
        currentTermExact = true;
        currentTerm = '';
        continue;
      }
      currentTerm += term[i];
    }
    if (currentTerm.trim().length > 0) {
      if (currentTermExact) {
        searchTerms.push(new searchTerm(currentTerm.replace('"', '').trim(), true));
      } else {
        searchTerms.push(new searchTerm(currentTerm.trim(), false));
      }
    }
    return searchTerms;
  }

  parseQuery(query: string) {
    let filter: Filter;
    try {
      filter = ParseAST(query);
    } catch (err) {
      console.error('failed to parse query', query, err);
      return;
    }

    const extractor = new FilterExtractor();
    filter.accept(extractor);

    this.activeFilters = [];
    this.searchForm.get('term').setValue('');
    let termText = [];
    extractor.filters.forEach((compFilter: CompFilter) => {
      let meta = this.getMetaForField(compFilter.field);
      if (!meta) {
        console.error('unknown field', compFilter.field);
        return;
      }
      if (compFilter.value.v === '') {
        return;
      }
      if (compFilter.field === 'content') {
        if (compFilter.op === CompOp.Like) {
          termText.push(compFilter.value.v);
        } else {
          termText.push(`"${compFilter.value.v}"`);
        }
        return;
      }
      this.activeFilters.push(new SearchFilter(meta, compFilter.op, compFilter.value.v));
    });
    this.searchForm.get('term').setValue(termText.join(' '));
  }
}

class FilterExtractor implements Visitor {

  filters: CompFilter[] = [];

  visitBoolFilter(f: BoolFilter): Visitor {
    f.lhs.accept(this);
    f.rhs.accept(this);
    return this;
  }

  visitCompFilter(f: CompFilter): Visitor {
    this.filters.push(f);
    return this;
  }

}

class searchTerm {
  constructor(public value: string, public  exact: boolean) {
  }
}
