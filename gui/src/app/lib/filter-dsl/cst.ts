import { Scan, Scanner, Tag, tagPrecedence, Tok } from './scanner';
import { Bool, Float, Int, Invalid, Str, Value, ValueKind } from './value';
import { trimChars } from '../util';
import { isBoolOp } from './filter';
import { Renderer2 } from '@angular/core';

export class ParseError {
  constructor(readonly reason: string, readonly cause: Tok = null) {
  }

  string() {
    return `${this.reason} (cause: ${this.cause.lexeme}`;
  }
}

enum NodeKind {
  BoolFilter = 'bool_filter',
  BoolOp = 'bool_op',
  CompFilter = 'comp_filter',
  CompOp = 'comp_op',
  Field = 'field',
  Value = 'value',
  Whitespace = 'whitespace',
  ParseError = 'parse_error',
  Unknown = 'unknown'
}

export class CSTNode {
  kind: NodeKind;
  children: CSTNode[] = [];

  constructor(kind: NodeKind, ...children: CSTNode[]) {
    this.kind = kind;
    if (children) {
      this.children.push(...children);
    }
  }

  valid(): boolean {
    return true;
  }

  appendChild(...ch: CSTNode[]) {
    this.children.push(...ch);
  }

  prependChild(...ch: CSTNode[]) {
    this.children.unshift(...ch);
  }

  string(): string {
    return this.children.map<string>((v) => v.string()).join('');
  }

  walk(fn: (v: CSTNode) => void): void {
    this.children.forEach((v) => {
      v.walk(fn);
    });
  }
}

export class TokenNode extends CSTNode {
  tok: Tok;

  constructor(kind: NodeKind, tok: Tok) {
    super(kind);
    this.tok = tok;
  }

  string(): string {
    return this.tok.lexeme;
  }

  walk(fn: (v: CSTNode) => void): void {
    fn(this);
  }
}

export class ValueNode extends CSTNode {
  constructor(public v: Value) {
    super(NodeKind.Value);
  }

  valid(): boolean {
    return this.v.kind !== ValueKind.InvalidValue;
  }

  string(): string {
    return this.v.token.lexeme;
  }

  walk(fn: (v: CSTNode) => void): void {
    fn(this);
  }
}

export function ParseCST(str: string) {
  return (new CSTParser(new Scanner(str, true))).parse();
}

export class CSTParser {
  private peeked: Tok = null;

  constructor(private s: Scanner) {
  }

  parse(): CSTNode {
    const node = this.parseOuter(1);
    this.requireTag(this.nextNonWhitespace(node), Tag.EOF);
    return node;
  }

  private parseOuter(minPrec: number): CSTNode {

    let innerNode = this.parseInner();

    while (true) {

      let outerNode: CSTNode = new CSTNode(NodeKind.BoolFilter);

      //some whitespace needs to be preserved before we know
      // which node will be returned.
      let whitespace = this.getWhitespace();

      const nextToken = this.peekNext();
      if (nextToken.tag == Tag.Error) {
        throw new ParseError(`scanner returned error`, nextToken);
      }
      console.log(nextToken);
      if (!isBoolOp(nextToken) || tagPrecedence(nextToken.tag) < minPrec) {
        if (whitespace.length > 0) {
          this.backtrackWhitespace();
        }
        //innerNode.appendChild(...whitespace);
        return innerNode;
      }

      // lhs
      outerNode.appendChild(innerNode);
      this.eatWhiteSpace(outerNode);

      // op
      let op = this.requireNext(Tag.And, Tag.Or)
      outerNode.appendChild(new TokenNode(NodeKind.BoolOp, op));
      this.eatWhiteSpace(outerNode);

      // rhs
      const rhs = this.parseOuter(tagPrecedence(op.tag) + 1);
      if (!rhs) {
        throw new ParseError("missing right hand statement", op)
      }
      outerNode.appendChild(rhs);
      this.eatWhiteSpace(outerNode);

      innerNode = outerNode;
    }
  }

  private parseInner(): CSTNode {

    console.log("inner");

    const node = new CSTNode(NodeKind.CompFilter);
    this.eatWhiteSpace(node);

    const t = this.getNext();
    switch (t.tag) {
      case Tag.EOF:
        break;
      case Tag.LParen:
        const expr = this.parseOuter(0);
        // bracket not handled
        this.requireNext(Tag.RParen);
        node.appendChild(expr);
        return node;
      case Tag.Field:

        // field
        node.appendChild(new TokenNode(NodeKind.Field, t));

        // op
        let next = this.nextNonWhitespace(node);
        this.requireTag(next, Tag.Eq, Tag.Neq, Tag.Lt, Tag.Le, Tag.Gt, Tag.Ge, Tag.Like);
        node.appendChild(new TokenNode(NodeKind.CompOp, next));

        // value
        this.eatWhiteSpace(node);
        node.appendChild(new ValueNode(this.parseValue()));

        return node;
      default:
        throw new ParseError(`unexpected token`, t);
    }

    throw new ParseError(`unexpected EOF`, t);
  }

  private nextNonWhitespace(node: CSTNode): Tok {
    let next = this.getNext();
    while (next.tag === Tag.Whitespace) {
      node.appendChild(new TokenNode(NodeKind.Whitespace, next));
      next = this.getNext();
    }
    return next;
  }

  private eatWhiteSpace(node: CSTNode): void {
    let next = this.peekNext();
    while (next.tag === Tag.Whitespace) {
      node.appendChild(new TokenNode(NodeKind.Whitespace, this.getNext()));
      next = this.peekNext();
    }
  }

  private getWhitespace(): TokenNode[] {
    const sp: TokenNode[] = [];
    let next = this.peekNext();
    while (next.tag === Tag.Whitespace) {
      sp.push(new TokenNode(NodeKind.Whitespace, this.getNext()))
      next = this.peekNext();
    }
    return sp
  }

  private parseValue(): Value {

    let token = this.getNext();
    switch (token.tag) {
      case Tag.Null:
        return null;
      case Tag.Int:
        return Int(parseInt(token.lexeme), token);
      case Tag.Float:
        return Float(parseFloat(token.lexeme), token);
      case Tag.Bool:
        if (token.lexeme === 'true') {
          return Bool(true, token);
        }
        if (token.lexeme === 'false') {
          return Bool(false, token);
        }
        throw new ParseError(`could not parse bool from value`, token);
      case Tag.String:
        return Str(trimChars(token.lexeme, '"'), token);
      case Tag.Error:
        return Invalid(token.lexeme, token);
    }
    throw new ParseError(`statement was missing valid value`, token);
  }

  private backtrackWhitespace() {
    this.peeked = null;
    this.s.backtrackWhitespace();
  }

  private getNext(): Tok {
    if (this.peeked !== null) {
      const t = this.peeked;
      this.peeked = null;
      console.log("next", t.lexeme);
      return t;
    }
    let n = this.s.next()
    console.log("next", n.lexeme);
    return n;
  }

  private peekNext(): Tok {
    if (this.peeked === null) {
      this.peeked = this.getNext();
    }
    return this.peeked;
  }

  private requireNext(...oneOf: Tag[]): Tok {
    return this.requireTag(this.getNext(), ...oneOf);
  }

  private requireTag(t: Tok, ...oneOf: Tag[]): Tok {
    for (let i = 0; i < oneOf.length; i++) {
      if (oneOf[i] == t.tag) {
        return t;
      }
    }
    throw new ParseError(`expected one of [${oneOf.join(', ')}]`, t);
  }

}

export function renderCST(renderer: Renderer2, root: CSTNode): HTMLElement {
  const p = new CSTPrinter(renderer);
  p.print(root);
  return p.el;
}

class CSTPrinter {

  el: HTMLElement;

  constructor(private renderer: Renderer2) {
  }

  print(v: CSTNode) {

    this.el = this.span('', v.kind);

    v.children.forEach((v => {
      if (!v) {
        return;
      }
      if (v.kind === NodeKind.BoolFilter || v.kind === NodeKind.CompFilter) {
        this.el.appendChild(renderCST(this.renderer, v));
      } else {
        if (!v.valid()) {
          this.el.appendChild(this.span(v.string(), v.kind, 'invalid'));
        } else {
          this.el.appendChild(this.span(v.string(), v.kind));
        }
      }
    }));
  }

  private span(innerText: string, ...cl: string[]): HTMLElement {
    const el = this.renderer.createElement('span');
    el.className = [...cl, 'filter-el'].join(' ');
    el.textContent = innerText;
    return el;
  }
}



