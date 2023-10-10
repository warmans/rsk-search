import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OauthLoginComponent } from './oauth-login.component';

describe('RedditLoginComponent', () => {
  let component: OauthLoginComponent;
  let fixture: ComponentFixture<OauthLoginComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OauthLoginComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OauthLoginComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
