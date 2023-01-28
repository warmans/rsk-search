import { AfterViewInit, Component, ElementRef, EventEmitter, HostListener, OnDestroy, Output, Renderer2, ViewChild } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { distinctUntilChanged, map, takeUntil } from 'rxjs/operators';
import { And, BoolFilter, CompFilter, CompOp, Filter, Visitor } from 'src/app/lib/filter-dsl/filter';
import { Str } from 'src/app/lib/filter-dsl/value';
import { PrintPlainText } from 'src/app/lib/filter-dsl/printer';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { ParseAST } from 'src/app/lib/filter-dsl/ast';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';

@Component({
  selector: 'app-search-bar',
  templateUrl: './search-bar.component.html',
  styleUrls: ['./search-bar.component.scss']
})
export class SearchBarComponent implements OnDestroy, AfterViewInit {

  @Output()
  queryUpdated: EventEmitter<string> = new EventEmitter<string>();

  focusState: 'idle' | 'focus' | 'typing' = 'idle';

  caretContainer: 'term' | 'mention' | 'publication' = 'term';

  showHelp: boolean;

  // android chrome has really weird behavior for key up/down events that is hard to debug
  showDebugger: boolean = false;

  keyPress$: Subject<KeyboardEvent> = new Subject<KeyboardEvent>();

  debug: string[] = [];

  destroy$: Subject<void> = new Subject();

  lastActiveMentionElement: HTMLElement;
  lastActivePublicationElement: HTMLElement;

  // API for mentions
  mentionsDataFn: (prefix: string) => Observable<any> = (prefix: string) => this.apiClient.listFieldValues({
    field: 'actor',
    prefix: prefix
  }).pipe(map(res => res.values.map((v) => v.value)));

  publicationDataFn: (prefix: string) => Observable<any> = (prefix: string) => this.apiClient.listFieldValues({
    field: 'publication',
    prefix: prefix
  }).pipe(map(res => res.values.map((v) => v.value)));

  @ViewChild('componentRoot')
  componentRootEl: any;

  @ViewChild('termsInput')
  termsInput: ElementRef;

  @HostListener('document:click', ['$event'])
  clickOut(event) {
    if (this.componentRootEl.nativeElement.contains(event.target)) {
      this.setStateFocussed();
      return;
    }
    this.setStateIdle();
    this.showHelp = false;
  }

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute, private renderer: Renderer2) {
    this.route.queryParamMap.pipe(distinctUntilChanged(), takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      if (params.get('debug') === '1') {
        this.showDebugger = true;
      }
    });
  }

  ngAfterViewInit() {
    this.route.queryParamMap.pipe(distinctUntilChanged(), takeUntil(this.destroy$)).subscribe((params: ParamMap) => {
      if (params.get('q') === null || params.get('q').trim() === '') {
        this.resetTerms();
        return;
      }
      this.populateSearchBarFromQuery(params.get('q'));
    });
  }

  getTermText(): string {
    return this.termsInput?.nativeElement?.innerText;
  }

  onKeyup(key: KeyboardEvent): boolean {
    this.debugData(`KEYUP code: ${key.code || 'EMPTY'} key: ${key.key} keyCode: ${key.keyCode}`);
    this.caretContainer = this.identifyCaretContainer();
    return true;
  }

  onKeydown(key: KeyboardEvent): boolean {
    this.debugData(`KEYDOWN code: ${key.code || 'EMPTY'} key: ${key.key} keyCode: ${key.keyCode}`);

    this.caretContainer = this.identifyCaretContainer();

    this.setStateFocussed();

    if (this.getTermText() === '') {
      this.setStateIdle();
    } else {
      if ((key.key || key.code) !== 'Enter') {
        this.handleTyping();
      }
    }

    // pass to child components
    this.keyPress$.next(key);

    switch (this.caretContainer) {
      case 'mention':
        this.lastActiveMentionElement = this.getAnchorNodeOfCaret().parentElement;
        break;
      case 'publication':
        this.lastActivePublicationElement = this.getAnchorNodeOfCaret().parentElement;
        break;
      case 'term':
        switch (key.key) {
          case '@':
            this.insertMention();
            this.caretContainer = this.identifyCaretContainer();
            return false;
          case '~':
            this.insertPublication();
            this.caretContainer = this.identifyCaretContainer();
            return false;
        }
    }

    switch ((key.key || key.code)) {
      case 'ArrowDown':
        this.handleTyping();
        break;
      case 'ArrowUp':
        this.handleTyping();
        break;
      case 'ArrowRight':
        return true;
      case 'Enter':
        if (this.focusState === 'typing') {
          this.setStateIdle();
          return false;
        }
        this.emitQuery();
        return false;
      case 'Escape':
        this.setStateIdle();
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
  }

  setStateFocussed() {
    this.focusState = 'focus';
  }

  handleTyping() {
    this.focusState = 'typing';
    this.showHelp = false;
  }

  resetTerms() {
    if (this.termsInput) {
      this.termsInput.nativeElement.innerText = '';
    }
  }

  setTermText(term: string) {
    if (!this.termsInput) {
      return;
    }
    // preserve non-terms before doing anything to the field.
    const modifiers = this.extractNonTermsFromSearch();

    this.termsInput.nativeElement.innerText = '';
    this.renderer.appendChild(this.termsInput.nativeElement, this.renderer.createText(term));
    // if the search contains a @mention or something, then re-add it.
    modifiers.forEach((node) => {
      this.renderer.appendChild(this.termsInput.nativeElement, node);
    });
    this.setCaretPositionToEnd(this.termsInput.nativeElement);
  }

  setTermAndEmit(term: string | null) {
    if (term !== null) {
      this.setTermText(term);
    }
    this.emitQuery();
  }

  applyMentionText(actor: string) {
    this.lastActiveMentionElement.innerText = actor;
    // if the span tag goes right to the end of the parent element you cannot escape from it with the arrow keys.
    // Adding a space lets the user move out of the span.
    this.termsInput.nativeElement.innerHTML = this.termsInput.nativeElement.innerHTML + '&#xa0;';
    this.setCaretPositionToEnd(this.termsInput.nativeElement);
  }

  applyPublicationText(name: string) {
    this.lastActivePublicationElement.innerText = name;
    // if the span tag goes right to the end of the parent element you cannot escape from it with the arrow keys.
    // Adding a space lets the user move out of the span.
    this.termsInput.nativeElement.innerHTML = this.termsInput.nativeElement.innerHTML + '\xa0';
    this.setCaretPositionToEnd(this.termsInput.nativeElement);
  }

  emitQuery() {

    // basically just go through the children of the search input and convert each node to a query fragment, depending
    // on if it's a html element (e.g. a mention span) or just a text node (a term).
    let query: Filter = null;
    (this.termsInput.nativeElement.childNodes || []).forEach((node) => {
      let comp: Filter = null;
      switch (node.className) {
        // this is a html node
        case 'mention':
          comp = new CompFilter('actor', CompOp.Eq, Str(node.innerText.replace(/@/g, '').trim()));
          break;
        case 'publication':
          comp = new CompFilter('publication', CompOp.Eq, Str(node.innerText.replace(/~/g, '').trim()));
          break;
        default:
          // this is a text node
          const text = (node.textContent || node.wholeText || node.innerHTML)?.replace(/&nbsp/g, ' ').trim();
          if (text === '' || text === undefined) {
            return;
          }
          // identify any quoted sections of the search
          this.extractTermGroups(text).forEach((line) => {
            let lineComp: CompFilter;
            if (line.length > 2 && line[0] === '"' && line[line.length - 1] === '"') {
              // quoted
              lineComp = new CompFilter('content', CompOp.Eq, Str(line.replace(/"/g, '').trim()));
            } else {
              // unquoted
              lineComp = new CompFilter('content', CompOp.FuzzyLike, Str(line.trim()));
            }
            if (comp) {
              comp = And(comp, lineComp);
            } else {
              comp = lineComp;
            }
          });
      }
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
    return;
  }

  extractTermGroups(nodeText: string): string[] {
    let termGroups: string[] = [];
    let currentTerm: string = '';
    let currentTermQuoted: boolean = false;
    for (let i = 0; i < nodeText.length; i++) {
      if (nodeText[i] === '"') {
        if (!currentTerm) {
          currentTermQuoted = true;
        } else {
          if (currentTermQuoted) {
            // this is a terminating quote
            termGroups.push(currentTerm + nodeText[i]);
            currentTerm = '';
            currentTermQuoted = false;
            continue;
          } else {
            // this is a new quote that terminates and un-quoted section
            termGroups.push(currentTerm);
            // start new quoted section
            currentTerm = nodeText[i];
            currentTermQuoted = true;
            continue;
          }
        }
      }
      currentTerm += nodeText[i];
    }
    if (currentTerm.length > 0) {
      termGroups.push(currentTerm);
    }
    return termGroups;
  }

  populateSearchBarFromQuery(query: string) {
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

    extractor.filters.forEach((compFilter: CompFilter) => {
      if (compFilter.value.v === '') {
        return;
      }
      if (this.termsInput.nativeElement.childNodes.length > 0) {
        // add a space if there are already elements in the search bar
        this.renderer.appendChild(this.termsInput.nativeElement, this.renderer.createText('\xa0'));
      }
      if (compFilter.field === 'content') {
        if (compFilter.op === CompOp.Like || compFilter.op === CompOp.FuzzyLike) {
          this.renderer.appendChild(this.termsInput.nativeElement, this.renderer.createText(compFilter.value.v));
        } else {
          this.renderer.appendChild(this.termsInput.nativeElement, this.renderer.createText(`"${compFilter.value.v}"`));
        }
        return;
      }
      if (compFilter.field === 'actor') {
        this.insertMention(compFilter.value.v, true);
        return;
      }
      if (compFilter.field === 'publication') {
        this.insertPublication(compFilter.value.v, true);
        return;
      }
    });
  }

  extractNonTermsFromSearch(): Node[] {
    const nodes: Node[] = [];
    this.termsInput.nativeElement.childNodes.forEach((node) => {
      switch (node.className) {
        // this is a html node
        case 'mention':
          nodes.push(node);
          return;
        case 'publication':
          nodes.push(node);
          return;
      }
    });
    return nodes;
  }

  insertMention(actor?: string, ignoreCaret?: boolean) {
    if (this.termsInput.nativeElement.textContent === '') {
      // clear any BR tags the browser has added
      this.termsInput.nativeElement.innerHTML = '';
    }
    let mention = this.renderer.createElement('span');
    mention.className = 'mention';
    mention.innerText = actor ? `@${actor}` : '@';

    this.renderer.appendChild(this.termsInput.nativeElement, mention);
    if (!ignoreCaret) {
      this.setCaretPositionToEnd(mention);
    }
    // add space so that the user can get out of the span
    this.renderer.appendChild(this.termsInput.nativeElement, this.renderer.createText('\xa0'));
    this.lastActiveMentionElement = mention;
  }

  insertPublication(name?: string, ignoreCaret?: boolean) {
    if (this.termsInput.nativeElement.textContent === '') {
      // clear any BR tags the browser has added
      this.termsInput.nativeElement.innerHTML = '';
    }
    let pub = this.renderer.createElement('span');
    pub.className = 'publication';
    pub.innerText = name ? `~${name}` : '~';

    this.renderer.appendChild(this.termsInput.nativeElement, pub);
    if (!ignoreCaret) {
      this.setCaretPositionToEnd(pub);
    }
    // add space so that the user can get out of the span
    this.renderer.appendChild(this.termsInput.nativeElement, this.renderer.createText('\xa0'));
    this.lastActiveMentionElement = pub;
  }

  private nodeIsChildOf(parent: HTMLElement, el: HTMLElement): boolean {
    if (el === parent) {
      return true;
    }
    if (el.parentElement !== null) {
      return this.nodeIsChildOf(parent, el.parentElement);
    }
    return false;
  }

  setCaretPositionToEnd(el: HTMLElement) {
    if (!el.childNodes?.length) {
      return;
    }
    try {
      let range = document.createRange();
      let sel = window.getSelection();
      const lastChild = el.childNodes[el.childNodes.length - 1];
      if (lastChild) {
        range.setStart(lastChild, lastChild.textContent.length);
      } else {
        range.setStart(el, 1);
      }
      range.collapse(true);
      sel.removeAllRanges();
      sel.addRange(range);
    } catch (e) {
      console.error('Failed to move cursor', e);
    }
  }

  getAnchorNodeOfCaret(): Node {
    let sel = document.getSelection();
    return sel.anchorNode;
  }

  // The container of the caret defines what sort of auto-complete would be relevant e.g. a mention, a term, or some other
  // search type.

  identifyCaretContainer(): 'mention' | 'term' | 'publication' {
    let node = this.getAnchorNodeOfCaret();
    let htmlNode = node as HTMLElement;
    let className = '';

    if (htmlNode?.className) {
      className = htmlNode?.className;
    } else {
      className = node?.parentElement?.className;
    }
    switch (className) {
      case 'mention':
        return 'mention';
      case 'publication':
        return 'publication';
      default:
        return 'term';
    }
  }

  toggleHelp() {
    this.showHelp = !this.showHelp;
  }

  debugData(msg: string) {
    if (this.showDebugger) {
      this.debug.unshift(msg);
      if (this.debug.length > 10) {
        this.debug.splice(10, this.debug.length - 10);
      }
    }
  }

  onInput($event: InputEvent) {
    this.debugData(`INPUT ${$event.data}`);
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
