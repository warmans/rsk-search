import { Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { MetaService } from '../../../core/service/meta/meta.service';
import { RsksearchFieldMeta, RsksearchFieldValue, RsksearchFieldValueList } from '../../../../lib/api-client/models';
import { Observable, of, Subject } from 'rxjs';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { catchError, distinctUntilChanged, map, switchMap, takeUntil, tap } from 'rxjs/operators';

@Component({
  selector: 'app-gl-search-filter',
  templateUrl: './gl-search-filter.component.html',
  styleUrls: ['./gl-search-filter.component.scss']
})
export class GlSearchFilterComponent implements OnInit, OnDestroy {

  // field to search in
  public _field: RsksearchFieldMeta;

  @Input()
  set field(f: RsksearchFieldMeta) {
    this._field = f;
    this.updateOperators();
  }

  get field(): RsksearchFieldMeta {
    return this._field;
  }

  // operator
  public possibleOperators: string[];
  public operator: string = '=';

  // value with autocomplete
  public valueInput$: Subject<string> = new Subject<string>();
  public possibleValues$: Observable<RsksearchFieldValue[]>;
  public value: any = null;
  public valuesLoading: boolean = false;
  public searchPrefix: string = '';

  // e.g. and/or
  public connector: string = 'and';

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
        takeUntil(this.destroy$),
        distinctUntilChanged(),
        tap(() => this.valuesLoading = true),
        switchMap(term => this.apiClient.searchServiceListFieldValues({
          field: this._field.name,
          prefix: term,
        }).pipe(
          map((v: RsksearchFieldValueList): RsksearchFieldValue[] => v.values),
          catchError(() => of([])), // empty list on error
          tap(() => this.valuesLoading = false)
        ))
      );
  }

  updateOperators() {
    this.possibleOperators = this.meta.getOperatorsForType(this._field.kind);
  }

  trackByFn(item: RsksearchFieldValue) {
    return item.value;
  }

}
