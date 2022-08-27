import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AudioPlayerFixedComponent } from './audio-player-fixed.component';

describe('AudioPlayerFixedComponent', () => {
  let component: AudioPlayerFixedComponent;
  let fixture: ComponentFixture<AudioPlayerFixedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AudioPlayerFixedComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioPlayerFixedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
