import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QuotasComponent } from './quotas.component';

describe('QuotasComponent', () => {
  let component: QuotasComponent;
  let fixture: ComponentFixture<QuotasComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ QuotasComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QuotasComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
