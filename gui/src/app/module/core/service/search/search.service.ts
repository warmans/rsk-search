import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SearchService {

  public queries: BehaviorSubject<string> = new BehaviorSubject<string>(null);

  constructor() {
  }

  emitQuery(query: string) {
    this.queries.next(query);
  }

}
