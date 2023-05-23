import {BoolFilter, CompFilter, Filter, Visitor} from './filter';
import {Renderer2} from '@angular/core';
import {ValueKind} from './value';

export function PrintPlainText(f: Filter): string {
  const p = new PlainTextPrinter();
  f.accept(p);
  return p.result();
}

export class PlainTextPrinter implements Visitor {

  private buff: string[] = [];

  constructor() {
  }

  result(): string {
    return this.buff.join(' ');
  }

  visitBoolFilter(f: BoolFilter): Visitor {

    const needsLparen = f.lhs.precedence() < f.precedence();
    if (needsLparen) {
      this.buff.push('(');
    }
    f.lhs.accept(this);
    if (needsLparen) {
      this.buff.push(')');
    }
    this.buff.push(f.op);

    const needsRparen = f.rhs.precedence() < f.precedence();
    if (needsRparen) {
      this.buff.push('(');
    }
    f.rhs.accept(this);
    if (needsRparen) {
      this.buff.push(')');
    }
    return this;
  }

  visitCompFilter(f: CompFilter): Visitor {
    if (f.value.kind === ValueKind.String) {
      this.buff.push(f.field, f.op, `"${f.value.v}"`);
    } else if (f.value.kind === ValueKind.Regexp) {
      this.buff.push(f.field, f.op, `/${f.value.v}/`);
    } else {
      this.buff.push(f.field, f.op, '' + f.value.v);
    }
    return this;
  }
}

export function PrintHTML(renderer: Renderer2, f: Filter): HTMLElement {
  const p = new HTMLPrinter(renderer);
  f.accept(p);
  return p.el;
}

export class HTMLPrinter implements Visitor {

  readonly el: HTMLElement;

  constructor(private renderer: Renderer2) {
    this.el = renderer.createElement('span');
  }

  visitBoolFilter(f: BoolFilter): Visitor {

    const boolEl = this.span(['bool-filter'], '');

    const needsLparen = f.lhs.precedence() < f.precedence();

    if (needsLparen) {
      boolEl.appendChild(this.span(['paren'], '('));
    }

    const lhsPrinter = new HTMLPrinter(this.renderer);
    f.lhs.accept(lhsPrinter);
    boolEl.appendChild(lhsPrinter.el);

    if (needsLparen) {
      boolEl.appendChild(this.span(['paren'], ')'));
    }

    boolEl.appendChild(this.span(['bool-op'], ` ${f.op} `));

    if (f.rhs) {
      const needsRparen = f.rhs.precedence() < f.precedence();

      if (needsRparen) {
        boolEl.appendChild(this.span(['paren'], '('));
      }

      const rhsPrinter = new HTMLPrinter(this.renderer);
      f.rhs.accept(rhsPrinter);
      boolEl.appendChild(rhsPrinter.el);

      if (needsRparen) {
        boolEl.appendChild(this.span(['paren'], ')'));
      }
    }

    this.el.appendChild(boolEl);
    return this;
  }

  visitCompFilter(f: CompFilter): Visitor {

    let value: string;
    switch (f.value.kind) {
      case ValueKind.String:
        value = `"${f.value.v}"`;
        break;
      default:
        value = `${f.value.v}`;
    }

    const compEl = this.span(['comp-filter'], '');

    compEl.appendChild(this.span(['field'], `${f.field} `));
    compEl.appendChild(this.span(['comp-op'], `${f.op} `));
    compEl.appendChild(this.span(['value'], value));

    this.el.appendChild(compEl);
    return this;
  }

  private span(cl: string[], innerText: string): HTMLElement {
    const el = this.renderer.createElement('span');
    el.className = [...cl, 'filter-el'].join(' ');
    el.textContent = innerText;
    return el;
  }
}
