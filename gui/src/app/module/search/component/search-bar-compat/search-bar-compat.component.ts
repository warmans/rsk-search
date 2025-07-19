import {
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  OnDestroy,
  OnInit,
  Output,
  Renderer2,
  ViewChild
} from '@angular/core';
import {UntypedFormControl} from '@angular/forms';
import {Observable, Subject} from 'rxjs';
import {distinctUntilChanged, map, takeUntil} from 'rxjs/operators';
import {ParseTerms, Term} from 'src/app/lib/search-parser/parser';
import {PrintFilterString, PrintPlaintext} from 'src/app/lib/search-parser/printer';
import {getInputSelection} from 'src/app/lib/caret';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
import {CompFilter, CompOp, Filter} from 'src/app/lib/filter-dsl/filter';
import {ParseAST} from 'src/app/lib/filter-dsl/ast';
import {FilterExtractor} from 'src/app/lib/filter-dsl/util';
import {TermsToFilter} from 'src/app/lib/search-parser/util';
import {Suggestion, SuggestionType} from "../search-bar-suggestion/search-bar-suggestion.component";

@Component({
    selector: 'app-search-bar-compat',
    templateUrl: './search-bar-compat.component.html',
    styleUrls: ['./search-bar-compat.component.scss'],
    standalone: false
})
export class SearchBarCompatComponent implements OnInit, OnDestroy {

  @Output()
  queryUpdated: EventEmitter<string> = new EventEmitter<string>();

  focusState: 'idle' | 'focus' | 'typing' = 'idle';

  inputFormControl: UntypedFormControl = new UntypedFormControl('');

  destroy$: Subject<void> = new Subject<void>();

  keyPress$: Subject<KeyboardEvent> = new Subject<KeyboardEvent>();

  terms: Term[] = [];

  contentFilters: Filter;

  activeTerm: Term;

  query: string;

  showHelp: boolean;

  lastParsed: string;

  // API for mentions
  mentionsDataFn: (prefix: string, filter: Filter, exact: boolean) => Observable<any> = (prefix: string, filter: Filter, exact: boolean) => this.apiClient.listFieldValues({
    field: 'actor',
    prefix: prefix,
  }).pipe(map(res => res.values.map((v) => v.value)));

  publicationDataFn: (prefix: string, filter: Filter, exact: boolean) => Observable<any> = (prefix: string, filter: Filter, exact: boolean) => this.apiClient.listFieldValues({
    field: 'publication',
    prefix: prefix,
  }).pipe(map(res => res.values.map((v) => v.value)));

  contentDataFn: (prefix: string, filter: Filter, exact: boolean) => Observable<any> = (prefix: string, filter: Filter, exact: boolean) => this.apiClient.predictSearchTerm({
    prefix: prefix,
    maxPredictions: 10,
    query: filter ? filter.print() : '',
    exact: exact,
  }).pipe(map(res => res.predictions));

  @ViewChild('componentRoot')
  componentRootEl: any;

  @ViewChild('termsInput')
  termsInput: ElementRef;

  constructor(
    private renderer: Renderer2,
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.route.queryParamMap.pipe(distinctUntilChanged(), takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      if (params.get('q') === null || params.get('q').trim() === '') {
        this.reset();
        return;
      }
      this.createTermsFromFilter(params.get('q'));
      this.setStateIdle();
    });
  }

  @HostListener('document:click', ['$event'])
  clickOut(event) {
    if (this.componentRootEl.nativeElement.contains(event.target)) {
      this.setStateFocussed();
      return;
    }
    this.setStateIdle();
    this.showHelp = false;
  }

  ngOnInit(): void {

  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  parseAndApplyTermsString(val: string) {
    // prevent re-parsing the same string
    if (this.lastParsed === val) {
      return;
    }
    this.lastParsed = val;

    if (!val) {
      this.terms = [];
      this.query = '';
      return;
    }
    this.terms = ParseTerms(val).filter((term) => term.value.trim() != '');
    this.query = PrintFilterString(this.terms);
    this.activeTerm = undefined;
    const caretPos = this.getCaretPos();
    this.terms.forEach((term) => {
      if (caretPos >= term.tok.start && caretPos <= term.tok.end) {
        this.activeTerm = term;
      }
    });
    this.contentFilters = TermsToFilter(this.terms.filter((t: Term) => t.value !== (this.activeTerm?.value || '')));
  }

  emitQuery() {
    this.queryUpdated.next(this.query);
  }

  //
  // // this will not work properly on android chrome due to the android keyboard.
  // // so, basically anything in here can only enhance the UX, not implement any
  // // important aspect of it.
  onKeydown(key: KeyboardEvent): boolean {
    this.keyPress$.next(key);
    switch ((key.key || key.code)) {
      case 'ArrowDown':
        return false;
      case 'ArrowUp':
        return false;
      case 'Enter':
        return this.focusState === 'idle' || this.focusState === 'focus';
      case 'Escape':
        this.setStateIdle();
        break;
      default:
        this.setStateTyping();
        return true;
    }
    return true;
  }

  onKeyup(ev: KeyboardEvent) {
    if ((ev.key || ev.code) === 'Escape') {
      return;
    }
    this.parseAndApplyTermsString(this.inputFormControl.value);
  }

  setStateIdle() {
    this.focusState = 'idle';
  }

  setStateFocussed() {
    this.focusState = 'focus';
  }

  setStateTyping() {
    this.focusState = 'typing';
    this.showHelp = false;
  }

  toggleHelp() {
    this.showHelp = !this.showHelp;
  }

  getCaretPos(): number {
    if (!this.termsInput) {
      return 0;
    }
    return getInputSelection(this.termsInput.nativeElement).end;
  }

  applySuggestion(suggestion: Suggestion) {
    if (suggestion.type === SuggestionType.Dialog) {
      this.router.navigate(['/ep', suggestion.epid], {fragment: `pos-${suggestion.pos}`});
      return;
    }
    const hasWhitespace: boolean = (/\s/).test(suggestion.term);
    const withoutQuotes: string = suggestion.term.replace(/"/g, '');
    switch (this.activeTerm?.field) {
      case 'content':
        this.activeTerm = new Term(
          {start: 0, end: withoutQuotes.length},
          this.activeTerm.field,
          withoutQuotes,
          hasWhitespace ? CompOp.Eq : CompOp.Like,
        );
        this.terms = [...this.terms.filter((t) => t === this.activeTerm || t.field !== 'content'), this.activeTerm];
        break;
      default:
        this.activeTerm.value = hasWhitespace ? suggestion.term : withoutQuotes;
        break;
    }
    this.renderTerms();
    this.parseAndApplyTermsString(this.inputFormControl.value);
    this.emitQuery();
    this.setStateIdle();
  }

  renderTerms() {
    this.inputFormControl.setValue(PrintPlaintext(this.terms));
  }

  createTermsFromFilter(query: string) {
    if (!query || query.trim() === '') {
      return;
    }
    let filter: Filter;
    try {
      filter = ParseAST(query);
    } catch (err) {
      console.error('failed to parse query', query, err);
      return;
    }

    const extractor = new FilterExtractor();
    filter.accept(extractor);

    this.reset();

    let terms: string[] = [];

    extractor.filters.forEach((compFilter: CompFilter) => {
      if (compFilter.value.v === '') {
        return;
      }

      const hasWhitespace = (/\s/).test(compFilter.value.v);
      if (compFilter.field === 'content') {
        if (compFilter.op === CompOp.Like || compFilter.op === CompOp.FuzzyLike) {
          terms.push(`${compFilter.value.v}`);
        } else {
          terms.push(`"${compFilter.value.v}"`);
        }
        return;
      }
      if (compFilter.field === 'actor') {
        terms.push(hasWhitespace ? `@"${compFilter.value.v}"` : `@${compFilter.value.v}`);
        return;
      }
      if (compFilter.field === 'publication') {
        terms.push(hasWhitespace ? `~"${compFilter.value.v}"` : `~${compFilter.value.v}`);
        return;
      }
    });

    const termsStr = terms.join(' ');
    this.inputFormControl.setValue(terms.join(' '));
    this.parseAndApplyTermsString(termsStr);
  }

  reset() {
    this.terms = [];
    this.activeTerm = undefined;
    this.inputFormControl.reset();
    this.setStateIdle();
  }
}
