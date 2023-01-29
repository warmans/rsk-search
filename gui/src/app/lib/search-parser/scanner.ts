export enum Tag {
  EOF = 'EOF',

  Mention = '@',
  Publication = '~',

  QuotedString = 'QUOTED_STRING',
  Word = 'WORD',

  Whitespace = 'WHITESPACE',
  Error = 'ERROR',
}

export class Tok {

  constructor(
    public tag: Tag,
    public lexeme: string,
    public start: number,
    public end: number,
    public error: string = null
  ) {
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
    //this.skipWhitespace();
    if (this.atEOF()) {
      return this.emit(Tag.EOF);
    }
    const c = this.nextChar();
    switch (c) {
      case '@':
        return this.emit(Tag.Mention);
      case '~':
        return this.emit(Tag.Publication);
      case '"':
        return this.scanQuotedString();
      default:
        if (this.isWhitespace(c)) {
          return this.emit(Tag.Whitespace);
        }
        if (this.isValidInputChar(c)) {
          return this.scanWord();
        }
        return this.emitError(`unknown entity: ${c}`);
    }
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

  private scanWord(): Tok {
    while (!this.atEOF() && (this.isValidInputChar(this.peekChar())) && !this.isWhitespace(this.peekChar())) {
      this.nextChar();
    }
    return this.emit(Tag.Word);
  }

  private scanQuotedString(): Tok {
    while (!this.matchNextChar('"')) {
      if (this.atEOF()) {
        // implicit quote close on EOF
        return this.emit(Tag.QuotedString);
      }
      this.nextChar();
    }
    return this.emit(Tag.QuotedString);
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
    return r != '@' && r != '~' && r != `"`;
  }
}
