export enum ValueKind {
  String,
  Int,
  Float,
  Bool,
  PartialString
}

export class Value {
  constructor(public kind: ValueKind, public v: any) {
  }
}

export function Str(value: string): Value {
  return new Value(ValueKind.String, value);
}

export function PartialStr(value: string): Value {
  return new Value(ValueKind.PartialString, value);
}

export function Int(value: number): Value {
  return new Value(ValueKind.Int, value);
}

export function Float(value: number): Value {
  return new Value(ValueKind.Float, value);
}

export function Bool(value: boolean): Value {
  return new Value(ValueKind.Bool, value);
}
