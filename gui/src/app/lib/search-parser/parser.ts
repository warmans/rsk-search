import { CompOp, Filter, NewCompFilter } from 'src/app/lib/filter-dsl/filter';
import { Scanner, Tag, Tok } from 'src/app/lib/search-parser/scanner';
import { Str } from 'src/app/lib/filter-dsl/value';


export class Pos {
  start: number;
  end: number;
}

export class Term {
  constructor(public tok: Pos, public field: string, public value: string, public op?: CompOp) {
    if (!op) {
      this.op = CompOp.Eq;
    }
  }

  public toFilter(): Filter {
    return NewCompFilter(this.field, this.op, Str((this.value || '').replace(/"/g, '')));
  }
}

// Search term parsing is a sort of simplification of the query DSL that is intended to be easier
// for humans to write (also using familiar concepts like @mentions and #hashtags) e.g.
// `"foo bar" @ricky` would become `content="foo bar" and actor="ricky"`
// it is permissive to unclosed entities to allow real-time parsing of input, but this means
// the output may be invalid once converted to the actual filter DSL

export function ParseTerms(str: string): Term[] {
  return (new TermParser(new Scanner(str))).parse();
}

export class TermParser {

  private peeked: Tok = null;

  constructor(private s: Scanner) {
  }

  parse(): Term[] {
    const terms = this.parseOuter();
    this.requireNext(Tag.EOF);
    return terms;
  }

  private parseOuter(): Term[] {
    const terms: Term[] = [];
    let term = this.parseInner();
    while (term != null) {
      terms.push(term);
      term = this.parseInner();
    }
    return terms;
  }

  private parseInner(): Term | null {
    const t = this.getNext();
    switch (t.tag) {
      case Tag.EOF:
        return null;
      case Tag.QuotedString:
        return new Term(t, 'content', t.lexeme, CompOp.Eq);
      case Tag.Word:
        return new Term(t, 'content', t.lexeme, CompOp.FuzzyLike);
      case Tag.Mention:
        const mentionText = this.requireNext(Tag.QuotedString, Tag.Word, Tag.EOF);
        return new Term({start: t.start, end: mentionText.end}, 'actor', mentionText.lexeme, CompOp.Eq);
      case Tag.Publication:
        const pubText = this.requireNext(Tag.QuotedString, Tag.Word, Tag.EOF);
        return new Term({start: t.start, end: pubText.end}, 'publication', pubText.lexeme, CompOp.Eq);
      case Tag.Whitespace:
        return new Term(t, 'content', t.lexeme, CompOp.FuzzyLike);
      default:
        throw new Error(`unexpected token ${t.lexeme}`);
    }
  }

  private getNext(): Tok {
    if (this.peeked !== null) {
      const t = this.peeked;
      this.peeked = null;
      return t;
    }
    return this.s.next();
  }

  private requireNext(...oneOf: Tag[]): Tok {
    const t = this.getNext();
    for (let i = 0; i < oneOf.length; i++) {
      if (oneOf[i] == t.tag) {
        return t;
      }
    }
    throw new Error(`expected one of [${oneOf.join(', ')}], found ${t.tag} (${t.lexeme})`);
  }
}

