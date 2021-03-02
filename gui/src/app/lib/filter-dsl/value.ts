import { Tok } from './scanner';

export enum ValueKind {
  String,
  Int,
  Float,
  Bool,
  InvalidValue,
}

export class Value {
  constructor(public kind: ValueKind, public v: any, public token: Tok = null) {
  }
}

export function Invalid(value: string, token: Tok = null): Value {
  return new Value(ValueKind.InvalidValue, value, token);
}

export function Str(value: string, token: Tok = null): Value {
  return new Value(ValueKind.String, value, token);
}

export function Int(value: number, token: Tok = null): Value {
  return new Value(ValueKind.Int, value, token);
}

export function Float(value: number, token: Tok = null): Value {
  return new Value(ValueKind.Float, value, token);
}

export function Bool(value: boolean, token: Tok = null): Value {
  return new Value(ValueKind.Bool, value, token);
}
