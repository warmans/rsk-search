import { Tok } from './scanner';
import { FieldMetaKind, RskFieldMeta } from '../api-client/models';

export enum ValueKind {
  String,
  Int,
  Float,
  Bool,
  Regexp,
  InvalidValue,
}

export class Value {
  constructor(public kind: ValueKind, public v: any, public token: Tok = null) {
  }
}

export function Invalid(value: string, token: Tok = null): Value {
  return new Value(ValueKind.InvalidValue, value, token);
}

export function Regexp(value: string, token: Tok = null): Value {
  return new Value(ValueKind.Regexp, value, token);
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

export function ValueFromFieldMeta(m: RskFieldMeta, value: any): Value {
  switch (m.kind) {
    case FieldMetaKind.IDENTIFIER:
      return Str(value);
    case FieldMetaKind.KEYWORD:
      return Str(value);
    case FieldMetaKind.KEYWORD_LIST:
      console.error('not implemented');
      return Str(value);
    case FieldMetaKind.INT:
      return Int(value);
    case FieldMetaKind.FLOAT:
      return Float(value);
    case FieldMetaKind.DATE:
      return Str(value);
    case FieldMetaKind.UNKNOWN:
      console.error('unknown field kind, assuming text');
    default:
      console.error('unhandled field kind, assuming text');
    case FieldMetaKind.TEXT:
      return Str(value);
  }

}
