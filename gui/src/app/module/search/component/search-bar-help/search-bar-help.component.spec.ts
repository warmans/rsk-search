import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchBarHelpComponent } from './search-bar-help.component';

describe('SearchBarHelpComponent', () => {
  let component: SearchBarHelpComponent;
  let fixture: ComponentFixture<SearchBarHelpComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SearchBarHelpComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchBarHelpComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
