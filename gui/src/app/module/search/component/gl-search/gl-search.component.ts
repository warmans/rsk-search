import { Component, OnInit } from '@angular/core';
import { MetaService } from '../../../core/service/meta/meta.service';
import { RsksearchFieldMeta } from '../../../../lib/api-client/models';
import { first } from 'rxjs/operators';

@Component({
  selector: 'app-gl-search',
  templateUrl: './gl-search.component.html',
  styleUrls: ['./gl-search.component.scss']
})
export class GlSearchComponent implements OnInit {

  fieldMeta: RsksearchFieldMeta[] = [];

  activeFilters: RsksearchFieldMeta[] = [];

  constructor(meta: MetaService) {
    meta.getMeta().pipe(first()).subscribe((m) => {
      this.fieldMeta = m.fields;
    });
  }

  ngOnInit(): void {
  }

  addFilter(field: string) {
    this.activeFilters.push(this.fieldMeta.find((v) => v.name === field));
  }

  removeFilter(idx: number) {
    this.activeFilters.splice(idx, 1);
  }
}


