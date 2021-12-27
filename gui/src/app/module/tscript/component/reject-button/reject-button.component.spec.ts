import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RejectButtonComponent } from './reject-button.component';

describe('RejectButtonComponent', () => {
  let component: RejectButtonComponent;
  let fixture: ComponentFixture<RejectButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RejectButtonComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RejectButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
