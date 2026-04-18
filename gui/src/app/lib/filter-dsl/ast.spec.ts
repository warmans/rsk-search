import { describe, expect, it } from 'vitest';
import { ParseAST } from './ast';
import { And, Eq, FuzzyLike, Ge, Gt, Le, Like, Lt, Neq, Or } from './filter';
import { Bool, Float, Int, Str } from './value';

describe('ParseAST', () => {
  it.each([
    {
      input: 'field = "value"',
      expected: Eq('field', Str('value')),
    },
    {
      input: 'field != 123',
      expected: Neq('field', Int(123)),
    },
    {
      input: 'field > 1.5',
      expected: Gt('field', Float(1.5)),
    },
    {
      input: 'field >= 10',
      expected: Ge('field', Int(10)),
    },
    {
      input: 'field < 20',
      expected: Lt('field', Int(20)),
    },
    {
      input: 'field <= 5.5',
      expected: Le('field', Float(5.5)),
    },
    {
      input: 'field ~= "pattern"',
      expected: Like('field', Str('pattern')),
    },
    {
      input: 'field ~ "fuzzy"',
      expected: FuzzyLike('field', Str('fuzzy')),
    },
    {
      input: 'field = true',
      expected: Eq('field', Bool(true)),
    },
    {
      input: 'field = false',
      expected: Eq('field', Bool(false)),
    },
    {
      input: 'field = null',
      expected: Eq('field', null),
    },
    {
      input: 'a = 1 and b = 2',
      expected: And(Eq('a', Int(1)), Eq('b', Int(2))),
    },
    {
      input: 'a = 1 or b = 2',
      expected: Or(Eq('a', Int(1)), Eq('b', Int(2))),
    },
    {
      input: 'a = 1 and b = 2 or c = 3',
      expected: Or(And(Eq('a', Int(1)), Eq('b', Int(2))), Eq('c', Int(3))),
    },
    {
      input: 'a = 1 or b = 2 and c = 3',
      expected: Or(Eq('a', Int(1)), And(Eq('b', Int(2)), Eq('c', Int(3)))),
    },
    {
      input: '(a = 1 or b = 2) and c = 3',
      expected: And(Or(Eq('a', Int(1)), Eq('b', Int(2))), Eq('c', Int(3))),
    },
    {
      input: 'a = 1 and (b = 2 or c = 3)',
      expected: And(Eq('a', Int(1)), Or(Eq('b', Int(2)), Eq('c', Int(3)))),
    },
  ])('should parse "$input" correctly', ({ input, expected }) => {
    const actual = ParseAST(input);
    expect(actual).toEqual(expected);
  });

  it('should throw error on "field ="', () => {
    expect(() => ParseAST('field =')).toThrow();
  });

  it('should throw error on "field"', () => {
    expect(() => ParseAST('field')).toThrow();
  });

  it('should throw error on "= \"value\""', () => {
    expect(() => ParseAST('= "value"')).toThrow();
  });
});
