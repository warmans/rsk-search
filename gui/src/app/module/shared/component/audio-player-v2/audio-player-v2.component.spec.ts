import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AudioPlayerV2Component } from './audio-player-v2.component';

describe('AudioPlayerV2Component', () => {
  let component: AudioPlayerV2Component;
  let fixture: ComponentFixture<AudioPlayerV2Component>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AudioPlayerV2Component ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioPlayerV2Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
