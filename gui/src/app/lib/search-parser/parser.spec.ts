import { ParseTerms } from 'src/app/lib/search-parser/parser';
import { CompOp } from 'src/app/lib/filter-dsl/filter';
import { describe, expect, it } from 'vitest';

describe('ParseTerms', () => {
  it.each([
    {
      input: 'foo bar',
      expected: [{ field: 'content', value: 'foo bar', op: CompOp.FuzzyLike, start: 0, end: 7 }],
    },
    {
      input: '"foo bar"',
      expected: [{ field: 'content', value: 'foo bar', op: CompOp.Eq, start: 0, end: 9 }],
    },
    {
      input: '@ricky',
      expected: [{ field: 'actor', value: 'ricky', op: CompOp.Eq, start: 0, end: 6 }],
    },
    {
      input: '~xfm',
      expected: [{ field: 'publication', value: 'xfm', op: CompOp.Eq, start: 0, end: 4 }],
    },
    {
      input: 'foo bar @ricky',
      expected: [
        { field: 'content', value: 'foo bar', op: CompOp.FuzzyLike, start: 0, end: 8 },
        { field: 'actor', value: 'ricky', op: CompOp.Eq, start: 8, end: 14 },
      ],
    },
    {
      input: 'foo   bar',
      expected: [{ field: 'content', value: 'foo bar', op: CompOp.FuzzyLike, start: 0, end: 9 }],
    },
    {
      input: '@"ricky gervais"',
      expected: [{ field: 'actor', value: '"ricky gervais"', op: CompOp.Eq, start: 0, end: 16 }],
    },
    {
      input: '',
      expected: [],
    },
    {
      input: ' ',
      expected: [{ field: 'content', value: ' ', op: CompOp.FuzzyLike, start: 0, end: 1 }],
    },
    {
      input: '"unclosed',
      expected: [{ field: 'content', value: 'unclosed', op: CompOp.Eq, start: 0, end: 9 }],
    },
    {
      input: '@',
      expected: [{ field: 'actor', value: '', op: CompOp.Eq, start: 0, end: 1 }],
    },
    {
      input: '~"xfm show"',
      expected: [{ field: 'publication', value: '"xfm show"', op: CompOp.Eq, start: 0, end: 11 }],
    },
    {
      input: 'foo @bar ~baz',
      expected: [
        { field: 'content', value: 'foo', op: CompOp.FuzzyLike, start: 0, end: 4 },
        { field: 'actor', value: 'bar', op: CompOp.Eq, start: 4, end: 8 },
        { field: 'content', value: ' ', op: CompOp.FuzzyLike, start: 8, end: 9 },
        { field: 'publication', value: 'baz', op: CompOp.Eq, start: 9, end: 13 },
      ],
    },
    {
      input: '@ricky "gervais"',
      expected: [
        { field: 'actor', value: 'ricky', op: CompOp.Eq, start: 0, end: 6 },
        { field: 'content', value: ' ', op: CompOp.FuzzyLike, start: 6, end: 7 },
        { field: 'content', value: 'gervais', op: CompOp.Eq, start: 7, end: 16 },
      ],
    },
  ])('should parse "$input" correctly', ({ input, expected }) => {
    const actual = ParseTerms(input).map((t) => ({
      field: t.field,
      value: t.value,
      op: t.op,
      start: t.tok.start,
      end: t.tok.end,
    }));
    expect(actual).toEqual(expected);
  });
});
