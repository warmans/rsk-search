import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchBarCompatComponent } from './search-bar-compat.component';

describe('SearchBarCompatComponent', () => {
  let component: SearchBarCompatComponent;
  let fixture: ComponentFixture<SearchBarCompatComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SearchBarCompatComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchBarCompatComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
