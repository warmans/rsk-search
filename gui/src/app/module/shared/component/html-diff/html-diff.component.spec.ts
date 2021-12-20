import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HtmlDiffComponent } from './html-diff.component';

describe('HtmlDiffComponent', () => {
  let component: HtmlDiffComponent;
  let fixture: ComponentFixture<HtmlDiffComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HtmlDiffComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HtmlDiffComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
