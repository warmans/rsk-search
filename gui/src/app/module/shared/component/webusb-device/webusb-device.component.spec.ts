import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WebusbDeviceComponent } from './webusb-device.component';

describe('WebusbDeviceComponent', () => {
  let component: WebusbDeviceComponent;
  let fixture: ComponentFixture<WebusbDeviceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WebusbDeviceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WebusbDeviceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
