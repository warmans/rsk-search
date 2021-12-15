import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EpisodeSummaryComponent } from './episode-summary.component';

describe('EpisodeSummaryComponent', () => {
  let component: EpisodeSummaryComponent;
  let fixture: ComponentFixture<EpisodeSummaryComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EpisodeSummaryComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(EpisodeSummaryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
