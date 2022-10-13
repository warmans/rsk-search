import { Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { MetaService } from '../../../core/service/meta/meta.service';
import { FieldMetaKind, RskFieldMeta, RskFieldValue, RskFieldValueList } from '../../../../lib/api-client/models';
import { Observable, of, Subject } from 'rxjs';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { catchError, distinctUntilChanged, filter, map, switchMap, takeUntil, tap } from 'rxjs/operators';
import { CompOp } from '../../../../lib/filter-dsl/filter';

export class SearchFilter {
  constructor(
    public meta: RskFieldMeta,
    public operator: CompOp = CompOp.Eq,
    public value: string = '',
  ) {
  }
}

@Component({
  selector: 'app-gl-search-filter',
  templateUrl: './gl-search-filter.component.html',
  styleUrls: ['./gl-search-filter.component.scss']
})
export class GlSearchFilterComponent implements OnInit, OnDestroy {

  // field to search in
  public _field: SearchFilter;

  @Input()
  set field(f: SearchFilter) {
    this._field = f;
    this.updateOperators();
  }

  get field(): SearchFilter {
    return this._field;
  }

  // operator
  public possibleOperators: string[];

  // value with autocomplete
  public valueInput$: Subject<string> = new Subject<string>();
  public possibleValues$: Observable<string[]>;
  public valuesLoading: boolean = false;

  public kinds = FieldMetaKind;

  private destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private meta: MetaService, private apiClient: SearchAPIClient) {

  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }

  ngOnInit(): void {
    this.possibleValues$ =
      this.valueInput$.pipe(
        distinctUntilChanged(),
        filter((v) => v !== null),
        tap(() => this.valuesLoading = true),
        switchMap(term => this.apiClient.listFieldValues({
          field: this._field.meta.name,
          prefix: term,
        }).pipe(
          map((v: RskFieldValueList): RskFieldValue[] => v.values),
          map((v: RskFieldValue[]): string[] => v.map((v: RskFieldValue) => v.value)),
          catchError(() => of([])), // empty list on error
          tap(() => this.valuesLoading = false)
        )),
        takeUntil(this.destroy$),
      );
  }

  updateOperators() {
    this.possibleOperators = this.meta.getOperatorsForType(this._field.meta.kind);
  }

  trackByFn(item: RskFieldValue) {
    return item.value;
  }

}
