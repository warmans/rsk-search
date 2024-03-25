import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RandomQuoteComponent } from './random-quote.component';

describe('RandomQuoteComponent', () => {
  let component: RandomQuoteComponent;
  let fixture: ComponentFixture<RandomQuoteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RandomQuoteComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RandomQuoteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
