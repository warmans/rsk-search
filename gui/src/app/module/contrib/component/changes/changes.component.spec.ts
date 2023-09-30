import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChangesComponent } from './changes.component';

describe('ChangesComponent', () => {
  let component: ChangesComponent;
  let fixture: ComponentFixture<ChangesComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ChangesComponent]
    });
    fixture = TestBed.createComponent(ChangesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
