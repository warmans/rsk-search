import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CommunityProjectsComponent } from './community-projects.component';

describe('CommunityProjectsComponent', () => {
  let component: CommunityProjectsComponent;
  let fixture: ComponentFixture<CommunityProjectsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommunityProjectsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(CommunityProjectsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
