import { BoolFilter, CompFilter, Visitor } from 'src/app/lib/filter-dsl/filter';

export class FilterExtractor implements Visitor {

  filters: CompFilter[] = [];

  visitBoolFilter(f: BoolFilter): Visitor {
    f.lhs.accept(this);
    f.rhs.accept(this);
    return this;
  }

  visitCompFilter(f: CompFilter): Visitor {
    this.filters.push(f);
    return this;
  }
}
