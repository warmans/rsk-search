import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DslSearchComponent } from './dsl-search.component';

describe('DslSearchComponent', () => {
  let component: DslSearchComponent;
  let fixture: ComponentFixture<DslSearchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DslSearchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DslSearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
