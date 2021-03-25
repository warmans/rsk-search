import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WebhidDeviceComponent } from './webhid-device.component';

describe('WebhidDeviceComponent', () => {
  let component: WebhidDeviceComponent;
  let fixture: ComponentFixture<WebhidDeviceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WebhidDeviceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WebhidDeviceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
