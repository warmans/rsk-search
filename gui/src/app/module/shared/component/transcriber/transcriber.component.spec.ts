import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TranscriberComponent } from './transcriber.component';

describe('TranscriberComponent', () => {
  let component: TranscriberComponent;
  let fixture: ComponentFixture<TranscriberComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TranscriberComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TranscriberComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
