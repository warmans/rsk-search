import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CausesComponent } from './causes.component';

describe('CausesComponent', () => {
  let component: CausesComponent;
  let fixture: ComponentFixture<CausesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CausesComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CausesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
