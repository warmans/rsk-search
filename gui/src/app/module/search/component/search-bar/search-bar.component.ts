import { Component, EventEmitter, HostListener, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Subject } from 'rxjs';
import { distinctUntilChanged, takeUntil } from 'rxjs/operators';
import { And, BoolFilter, CompFilter, CompOp, Filter, Visitor } from 'src/app/lib/filter-dsl/filter';
import { Str, Value } from 'src/app/lib/filter-dsl/value';
import { PrintPlainText } from 'src/app/lib/filter-dsl/printer';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { ParseAST } from 'src/app/lib/filter-dsl/ast';

export interface SearchModifier {
  field: string;
  value: Value;
  comparison: CompOp;
}

@Component({
  selector: 'app-search-bar',
  templateUrl: './search-bar.component.html',
  styleUrls: ['./search-bar.component.scss']
})
export class SearchBarComponent implements OnInit, OnDestroy {

  @Output()
  queryUpdated: EventEmitter<string> = new EventEmitter<string>();

  focusState: 'idle' | 'focus' | 'typing' = 'idle';

  searchDropdown: 'none' | 'advanced' | 'autocomplete' = 'none';

  termTextInput: FormControl = new FormControl();

  searchModifiers: CompFilter[] = [];

  destroy$: Subject<void> = new Subject();

  @ViewChild('componentRoot')
  componentRootEl: any;

  @HostListener('document:click', ['$event'])
  clickOut(event) {
    if (this.componentRootEl.nativeElement.contains(event.target)) {
      this.setStateFocussed();
      return;
    }
    this.setStateIdle();
  }

  constructor(route: ActivatedRoute) {
    route.queryParamMap.pipe(distinctUntilChanged(), takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      if (params.get('q') === null || params.get('q').trim() === '') {
        this.resetTerms();
        return;
      }
      this.parseQuery(params.get('q'));
    });
  }

  ngOnInit(): void {
    this.termTextInput.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((currentValue: string) => {
      if (currentValue === '') {
        this.setStateIdle();
      } else {
        this.setStateTyping();
      }
    });
  }

  onKeydown(key: KeyboardEvent): boolean {
    // todo: pass key presses to autocomplete
    this.setStateFocussed();
    switch (key.code) {
      case 'ArrowDown':
        break;
      case 'ArrowUp':
        break;
      case 'Enter':
        this.emitQuery();
        return true;
      case 'Escape':
        break;
      default:
        return true;
    }
    return false;
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  setStateIdle() {
    this.focusState = 'idle';
    this.searchDropdown = 'none';
  }

  setStateFocussed() {
    this.focusState = 'focus';
  }

  setStateTyping() {
    this.focusState = 'typing';
    this.searchDropdown = 'autocomplete';
  }

  resetTerms() {
    this.searchModifiers = [];
    this.termTextInput.setValue('');
  }

  emitQuery() {
    let term = this.termTextInput.value || '';

    // group terms by exact/non-exact and convert them into a single filter statement.
    let query: Filter = null;
    this.parseSearchTerm(term).forEach((comp: CompFilter) => {
      if (query == null) {
        query = comp;
        return;
      }
      query = And(query, comp);
    });
    this.searchModifiers.forEach((comp: CompFilter) => {
      if (query == null) {
        query = comp;
        return;
      }
      query = And(query, comp);
    });
    if (query !== null) {
      this.queryUpdated.emit(PrintPlainText(query));
    } else {
      this.queryUpdated.emit('');
    }
  }

  parseSearchTerm(term: string): CompFilter[] {
    let searchTerms: CompFilter[] = [];
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
            searchTerms.push(new CompFilter('content', CompOp.Eq, Str(currentTerm.trim())));
          }
          currentTermExact = false;
          currentTerm = '';
          continue;
        }
        // start of exact term (where previous term was vague)
        if (currentTerm.trim().length > 0) {
          searchTerms.push(new CompFilter('content', CompOp.FuzzyLike, Str(currentTerm.trim())));
        }
        currentTermExact = true;
        currentTerm = '';
        continue;
      }
      currentTerm += term[i];
    }
    if (currentTerm.trim().length > 0) {
      if (currentTermExact) {
        searchTerms.push(new CompFilter('content', CompOp.Eq, Str(currentTerm.replace('"', '').trim())));
      } else {
        searchTerms.push(new CompFilter('content', CompOp.FuzzyLike, Str(currentTerm.trim())));
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

    this.resetTerms();

    let termText = [];
    extractor.filters.forEach((compFilter: CompFilter) => {
      if (compFilter.value.v === '') {
        return;
      }
      if (compFilter.field === 'content') {
        if (compFilter.op === CompOp.Like || compFilter.op === CompOp.FuzzyLike) {
          termText.push(compFilter.value.v);
        } else {
          termText.push(`"${compFilter.value.v}"`);
        }
        return;
      }
      this.searchModifiers.push(compFilter);
    });
    this.termTextInput.setValue(termText.join(' '));
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
