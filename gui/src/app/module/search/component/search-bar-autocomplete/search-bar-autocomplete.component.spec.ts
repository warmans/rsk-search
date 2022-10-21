import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchBarAutocompleteComponent } from './search-bar-autocomplete.component';

describe('SearchBarAutocompleteComponent', () => {
  let component: SearchBarAutocompleteComponent;
  let fixture: ComponentFixture<SearchBarAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SearchBarAutocompleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchBarAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
