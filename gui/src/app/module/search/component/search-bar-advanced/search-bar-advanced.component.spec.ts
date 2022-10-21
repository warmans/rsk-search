import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchBarAdvancedComponent } from './search-bar-advanced.component';

describe('SearchBarAdvancedComponent', () => {
  let component: SearchBarAdvancedComponent;
  let fixture: ComponentFixture<SearchBarAdvancedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SearchBarAdvancedComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchBarAdvancedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
