import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditorConfigComponent } from './editor-config.component';

describe('EditorConfigComponent', () => {
  let component: EditorConfigComponent;
  let fixture: ComponentFixture<EditorConfigComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditorConfigComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditorConfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
