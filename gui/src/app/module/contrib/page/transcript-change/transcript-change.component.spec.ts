import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TranscriptChangeComponent } from './transcript-change.component';

describe('TranscriptChangeComponent', () => {
  let component: TranscriptChangeComponent;
  let fixture: ComponentFixture<TranscriptChangeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TranscriptChangeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TranscriptChangeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
