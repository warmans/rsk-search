import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditorHelpComponent } from './editor-help.component';

describe('EditorHelpComponent', () => {
  let component: EditorHelpComponent;
  let fixture: ComponentFixture<EditorHelpComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditorHelpComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditorHelpComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
