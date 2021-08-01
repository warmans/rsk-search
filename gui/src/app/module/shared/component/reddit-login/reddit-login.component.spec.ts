import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RedditLoginComponent } from './reddit-login.component';

describe('RedditLoginComponent', () => {
  let component: RedditLoginComponent;
  let fixture: ComponentFixture<RedditLoginComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RedditLoginComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RedditLoginComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
