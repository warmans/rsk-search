import { isNumeric } from 'rxjs/internal-compatibility';
import { trimChars } from '../util';

export enum Tag {
  EOF = 'EOF',

  LParen = '(',
  RParen = ')',

  And = 'AND',
  Or = 'OR',

  Eq = '=',
  Neq = '!=',
  Like = '~=',
  Gt = '>',
  Ge = '>=',
  Le = '<=',
  Lt = '<',

  Field = 'FIELD',
  Int = 'INT',
  Float = 'FLOAT',
  Bool = 'BOOL',
  String = 'STRING',
  Null = 'NULL',

  Whitespace = 'WHITESPACE',
  Error = 'ERROR',
}

const keywords = {
  'and': Tag.And,
  'or': Tag.Or,
  'true': Tag.Bool,
  'false': Tag.Bool,
  'null': Tag.Null,
};

export class Tok {

  constructor(
    public tag: Tag,
    public lexeme: string,
    public start: number,
    public end: number,
    public error: string = null
  ) {
  }

  trimLexeme(cutset: string): Tok {
    this.lexeme = trimChars(this.lexeme, cutset);
    return this;
  }
}

export function Scan(str: string, cstMode: boolean = false): Tok[] {
  const scanner = new Scanner(str, cstMode);
  let tokens = [];
  while (true) {
    const tok = scanner.next();
    tokens.push(tok);
    if (tok.tag === Tag.EOF) {
      break;
    }
  }
  return tokens;
}

export class Scanner {

  input: string[] = [];
  pos: number = 0;
  offset: number = 0;
  cstMode: boolean = false; //concrete syntax trees require whitespace chars to be tokenized.

  constructor(str: string, cstMode: boolean = false) {
    this.input = str.split('');
    this.cstMode = cstMode;
  }

  next(): Tok {
    if (!this.cstMode) {
      this.skipWhitespace();
    }
    if (this.atEOF()) {
      return this.emit(Tag.EOF);
    }
    const c = this.nextChar();
    switch (c) {
      case '(':
        return this.emit(Tag.LParen);
      case ')':
        return this.emit(Tag.RParen);
      case '=':
        return this.emit(Tag.Eq);
      case '!':
        if (this.matchNextChar('=')) {
          return this.emit(Tag.Neq);
        }
        return this.emitError('expected = after !');
      case '~':
        if (this.matchNextChar('=')) {
          return this.emit(Tag.Like);
        }
        return this.emitError('expected = after ~');
      case '>':
        if (this.matchNextChar('=')) {
          return this.emit(Tag.Ge);
        }
        return this.emit(Tag.Gt);
      case '<':
        if (this.matchNextChar('=')) {
          return this.emit(Tag.Le);
        }
        return this.emit(Tag.Lt);
      case '"':
        return this.scanString();
      default:
        if (this.isWhitespace(c)) {
          return this.emit(Tag.Whitespace);
        }
        if (this.isValidInputChar(c)) {
          return this.scanField();
        }
        if (this.isStartOfNumber(c)) {
          return this.scanNumber();
        }
        return this.emitError(`unknown entity: ${c}`);
    }
  }

  backtrackWhitespace() {
    while (!this.atEOF() && this.isWhitespace(this.input[this.pos-1])) {
      this.pos--;
    }
    this.offset = this.pos;
  }

  private skipWhitespace() {
    while (!this.atEOF() && this.isWhitespace(this.peekChar())) {
      this.nextChar();
    }
    this.offset = this.pos;
  }

  private isWhitespace(char: string): boolean {
    return char.match(/[\s]/) !== null;
  }

  private scanField(): Tok {
    //todo: not sure about IsNumeric here
    while (!this.atEOF() && (this.isValidInputChar(this.peekChar()) || isNumeric(this.peekChar()))) {
      this.nextChar();
    }
    const tok = this.emit(Tag.Field);
    const tag = keywords[tok.lexeme];
    if (tag !== undefined) {
      tok.tag = tag;
    }
    return tok;
  }

  private scanNumber(): Tok {
    let hasDecimal: boolean = false;
    while (!this.atEOF() && (isNumeric(this.peekChar()) || (this.peekChar() == '.' && !hasDecimal))) {
      const c = this.nextChar();
      hasDecimal = hasDecimal || c === '.';
    }
    if (hasDecimal) {
      return this.emit(Tag.Float);
    }
    return this.emit(Tag.Int);
  }

  private scanString(): Tok {
    while (!this.matchNextChar('"')) {
      if (this.atEOF()) {
        return this.emitError('unclosed quote');
      }
      this.nextChar();
    }
    return this.emit(Tag.String);
  }

  private nextChar(): string {
    this.pos++;
    return this.input[this.pos - 1];
  }

  private matchNextChar(c: string): boolean {
    if (this.atEOF() || this.peekChar() != c) {
      return false;
    }
    this.nextChar();
    return true;
  }

  private peekChar(): string {
    return this.input[this.pos];
  }

  private emit(tag: Tag): Tok {
    const lexeme = this.input.slice(this.offset, this.pos).join('');
    const start = this.offset;
    const end = this.pos;
    // advance
    this.offset = this.pos;
    //emit
    return new Tok(tag, lexeme, start, end);
  }

  private emitError(reason: string): Tok {
    const lexeme = this.input.slice(this.offset, this.input.length).join('');
    const start = this.offset;
    const end = this.input.length;
    // advance
    this.offset = this.input.length;
    //emit
    return new Tok(Tag.Error, lexeme, start, end);
  }

  private atEOF(): boolean {
    return this.pos >= this.input.length || this.pos === -1;
  }

  private isValidInputChar(r: string): boolean {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'; // does this work in js?
  }

  private isStartOfNumber(c: string): boolean {
    return isNumeric(c) || c == '-';
  }
}

export function tagPrecedence(tag: Tag): number {
  switch (tag) {
    case Tag.And:
      return 2;
    case Tag.Or:
      return 1;
  }
  return 0;
}
