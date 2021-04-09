import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PendingRewardsComponent } from './pending-rewards.component';

describe('PendingRewardsComponent', () => {
  let component: PendingRewardsComponent;
  let fixture: ComponentFixture<PendingRewardsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PendingRewardsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PendingRewardsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
