import { Scanner, Tag, tagPrecedence, Tok } from './scanner';
import { And, Eq, Filter, Ge, Gt, isBoolOp, Le, Like, Lt, Neq, Or } from './filter';
import { trimChars } from '../util';
import { Bool, Float, Int, Str, Value } from './value';

export function ParseAST(str: string): Filter {
  return (new ASTParser(new Scanner(str, false))).parse();
}

export class ASTParser {

  private peeked: Tok = null

  constructor(private s: Scanner) {
  }

  parse(): Filter {
    const filter = this.parseOuter(1);
    this.requireNext(Tag.EOF);
    return filter;
  }

  private parseOuter(minPrec: number): Filter {
    let filter = this.parseInner();
    while (true) {
      const nextToken = this.peekNext();

      if (!isBoolOp(nextToken) || tagPrecedence(nextToken.tag) < minPrec) {
        break;
      }

      let op = this.getNext();
      let rhs = this.parseOuter(tagPrecedence(op.tag) + 1);

      if (op.tag === Tag.And) {
        filter = And(filter, rhs);
      } else if (op.tag === Tag.Or) {
        filter = Or(filter, rhs);
      } else {
        throw new Error(`unexpected token ${op.tag}`);
      }
    }
    return filter;
  }

  private parseInner(): Filter {
    const t = this.getNext();
    switch (t.tag) {
      case Tag.EOF:
        break;
      case Tag.LParen:
        const filter = this.parseOuter(0);
        this.requireNext(Tag.RParen);
        return filter;
      case Tag.Field:
        const op = this.requireNext(Tag.Eq, Tag.Neq, Tag.Lt, Tag.Le, Tag.Gt, Tag.Ge, Tag.Like);
        switch (op.tag) {
          case Tag.Eq:
            return Eq(t.lexeme, this.parseValue());
          case Tag.Neq:
            return Neq(t.lexeme, this.parseValue());
          case Tag.Like:
            return Like(t.lexeme, this.parseValue());
          case Tag.Gt:
            return Gt(t.lexeme, this.parseValue());
          case Tag.Ge:
            return Ge(t.lexeme, this.parseValue());
          case Tag.Le:
            return Le(t.lexeme, this.parseValue());
          case Tag.Lt:
            return Lt(t.lexeme, this.parseValue());
        }
        throw new Error(`unexpected field field ${op.tag}`);
      default:
        throw new Error(`unexpected token ${t.lexeme}`);
    }
  }

  private parseValue(): Value {
    let token = this.getNext();
    switch (token.tag) {
      case Tag.Null:
        return null;
      case Tag.Int:
        return Int(parseInt(token.lexeme));
      case Tag.Float:
        return Float(parseFloat(token.lexeme));
      case Tag.Bool:
        if (token.lexeme === 'true') {
          return Bool(true);
        }
        if (token.lexeme === 'false') {
          return Bool(false);
        }
        throw new Error(`Could not parse bool from value: ${token.lexeme}`);
      case Tag.String:
        return Str(trimChars(token.lexeme, '"'));
    }
    throw new Error(`Unexpected value ${token.lexeme}`);
  }

  private getNext(): Tok {
    if (this.peeked !== null) {
      const t = this.peeked;
      this.peeked = null;
      return t;
    }
    return this.s.next();
  }

  private peekNext(): Tok {
    if (this.peeked !== null) {
      return this.peeked;
    }
    const t = this.getNext();
    this.peeked = t;
    return t;
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

