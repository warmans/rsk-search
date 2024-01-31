import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {interval, Observable} from "rxjs";
import {takeUntil, timeInterval} from "rxjs/operators";
import {differenceInDays, differenceInSeconds, formatDistanceToNow} from "date-fns";

const PENCE_PER_SECOND: number = 26 / (24 * 60 * 60)
const INITIAL_PENCE_OWED: number = 620;

@Component({
  selector: 'app-catalog-warehouse',
  standalone: false,
  templateUrl: './catalog-warehouse.component.html',
  styleUrl: './catalog-warehouse.component.scss'
})
export class CatalogWarehouseComponent implements OnInit, OnDestroy {

  public valueInPence = 0;
  public timer;
  public inDebtDays = 0;

  destroy$: EventEmitter<void> = new EventEmitter<void>();

  ngOnInit(): void {

    let startDate = new Date(2003, 0, 25, 12, 10, 0, 0)
    this.valueInPence = INITIAL_PENCE_OWED + (differenceInSeconds(new Date(), startDate) * PENCE_PER_SECOND);
    this.inDebtDays = differenceInDays(new Date(), startDate);

    const seconds = interval(1000);
    seconds.pipe(timeInterval(), takeUntil(this.destroy$))
      .subscribe(
        () => {
          this.valueInPence += PENCE_PER_SECOND;
        }
      );
  }

  ngOnDestroy(): void {
    this.destroy$.emit();
    this.destroy$.complete();
  }
}
