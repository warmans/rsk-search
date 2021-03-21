import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AuthorContributionsComponent } from './author-contributions.component';

describe('AuthorContributionsComponent', () => {
  let component: AuthorContributionsComponent;
  let fixture: ComponentFixture<AuthorContributionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AuthorContributionsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthorContributionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
