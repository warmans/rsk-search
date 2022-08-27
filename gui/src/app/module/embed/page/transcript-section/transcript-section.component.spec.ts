import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TranscriptSectionComponent } from './transcript-section.component';

describe('TranscriptSectionComponent', () => {
  let component: TranscriptSectionComponent;
  let fixture: ComponentFixture<TranscriptSectionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TranscriptSectionComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TranscriptSectionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
