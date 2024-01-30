import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CatalogWarehouseComponent } from './catalog-warehouse.component';

describe('CatalogWarehouseComponent', () => {
  let component: CatalogWarehouseComponent;
  let fixture: ComponentFixture<CatalogWarehouseComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CatalogWarehouseComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(CatalogWarehouseComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
