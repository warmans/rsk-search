import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FindReplaceComponent } from './find-replace.component';

describe('FindReplaceComponent', () => {
  let component: FindReplaceComponent;
  let fixture: ComponentFixture<FindReplaceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ FindReplaceComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FindReplaceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
