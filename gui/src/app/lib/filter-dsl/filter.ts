import { Tag, Tok } from './scanner';
import { Value } from './value';

export enum CompOp {
  Eq = '=',
  Neq = '!=',
  Like = '~=',
  Lt = '<',
  Le = '<=',
  Gt = '>',
  Ge = '>=',
}


function compOpPrecedence(op: CompOp): number {
  return 3;
}

export enum BoolOp {
  And = 'and',
  Or = 'or',
}

function boolOpPrecedence(op: BoolOp): number {
  switch (op) {
    case BoolOp.And:
      return 2;
    case BoolOp.Or:
      return 1;
  }
  return 0;
}

export function isBoolOp(token: Tok) {
  return token.tag === Tag.And || token.tag === Tag.Or;
}

export interface Visitor {
  visitCompFilter(f: CompFilter): Visitor

  visitBoolFilter(f: BoolFilter): Visitor
}

export interface Filter {
  accept(v: Visitor): void

  precedence(): number
}

export class CompFilter implements Filter {

  constructor(public field: string, public op: CompOp, public value: Value) {
  }

  accept(v: Visitor): void {
    v.visitCompFilter(this);
  }

  precedence(): number {
    return compOpPrecedence(this.op);
  }
}

export class BoolFilter implements Filter {

  constructor(public lhs: Filter, public op: BoolOp, public rhs: Filter) {
  }

  accept(v: Visitor): void {
    v.visitBoolFilter(this);
  }

  precedence(): number {
    return boolOpPrecedence(this.op);
  }
}

export function NewCompFilter(field: string, op: CompOp, value: Value): CompFilter {
  return new CompFilter(field, op, value);
}

export function And(lhs: Filter, rhs: Filter, ...filters: Filter[]): Filter {
  let filter = new BoolFilter(lhs, BoolOp.And, rhs);
  filters.forEach((f: Filter) => {
    filter = new BoolFilter(filter, BoolOp.And, f);
  });
  return filter;
}

export function Or(lhs: Filter, rhs: Filter, ...filters: Filter[]): Filter {
  let filter = new BoolFilter(lhs, BoolOp.Or, rhs);
  filters.forEach((f: Filter) => {
    filter = new BoolFilter(filter, BoolOp.Or, f);
  });
  return filter;
}


export function Eq(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Eq, value);
}

export function Neq(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Neq, value);
}

export function Gt(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Gt, value);
}

export function Ge(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Ge, value);
}

export function Lt(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Lt, value);
}

export function Le(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Le, value);
}

export function Like(field: string, value: Value): Filter {
  return new CompFilter(field, CompOp.Like, value);
}

