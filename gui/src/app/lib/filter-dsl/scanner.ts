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

  IncompleteString = 'INCOMPLETE_STRING'
}

const keywords = {
  'and': Tag.And,
  'or': Tag.Or,
  'true': Tag.Bool,
  'false': Tag.Bool,
  'null': Tag.Null,
};

export class Tok {
  public tag: Tag;
  public lexeme: string;
  public start: number;
  public end: number;

  constructor(tag: Tag, lexeme: string, start: number, end: number) {
    this.tag = tag;
    this.lexeme = lexeme;
    this.start = start;
    this.end = end;
  }

  trimLexeme(cutset: string): Tok {
    this.lexeme = trimChars(this.lexeme, cutset);
    return this;
  }
}

export function Scan(str: string): Tok[] {
  const scanner = new Scanner(str);
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

  constructor(str: string) {
    this.input = str.split('');
  }

  next(): Tok {
    this.skipWhitespace();
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
        throw new Error('expected = after !');
      case '~':
        if (this.matchNextChar('=')) {
          return this.emit(Tag.Like);
        }
        throw new Error('expected = after ~');
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
        if (this.isValidInputChar(c)) {
          return this.scanField();
        }
        if (this.isStartOfNumber(c)) {
          return this.scanNumber();
        }
        throw new Error(`Unknown entity: ${c}`);
    }
  }

  private skipWhitespace() {
    while (!this.atEOF() && this.peekChar().match(/\s/) !== null) {
      this.nextChar();
    }
    this.offset = this.pos;
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

  scanString(): Tok {
    while (!this.matchNextChar('"')) {
      if (this.atEOF()) {
        // bit weird - allow an incomplete string in parsing to allow autocomplete even if
        // only a single quote was added
        return this.emit(Tag.IncompleteString);
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

  private atEOF(): boolean {
    return this.pos >= this.input.length;
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
