import { Component, HostListener, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { Claims } from 'src/app/module/core/service/session/session.service';
import { Subject } from 'rxjs';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { NotificationKind, RskNotification } from 'src/app/lib/api-client/models';
import { debounceTime, takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-user-menu',
  templateUrl: './user-menu.component.html',
  styleUrls: ['./user-menu.component.scss']
})
export class UserMenuComponent implements OnInit, OnDestroy {

  @Input()
  loggedInUser: Claims;

  @ViewChild('componentRoot')
  componentRootEl: any;

  menuVisible: boolean = false;

  notifications: RskNotification[] = [];

  unreads: number = 0;

  destroy$: Subject<void> = new Subject<void>();

  notificationKinds = NotificationKind;

  markRead: Subject<void> = new Subject<void>();

  constructor(private apiClient: SearchAPIClient) {
  }

  @HostListener('document:click', ['$event'])
  clickOut(event) {
    if (!this.menuVisible || this.componentRootEl.nativeElement.contains(event.target)) {
      return;
    }
    this.hideMenu();
  }

  ngOnInit(): void {
    this
      .apiClient
      .listNotifications({ filter: '', sortField: 'created_at', sortDirection: 'DESC', page: 1, pageSize: 5 })
      .pipe(takeUntil(this.destroy$))
      .subscribe((res) => {
        this.notifications = res.notifications;
        this.unreads = (res.notifications.filter((n: RskNotification) => !n.readAt) || []).length;
      });

    this
      .markRead
      .pipe(takeUntil(this.destroy$), debounceTime(1000))
      .subscribe(() => {
        if (this.unreads > 0) {
          this.apiClient.markNotificationsRead().pipe(takeUntil(this.destroy$)).subscribe(() => {
            this.unreads = 0;
          });
        }
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  toggleMenu() {
    this.menuVisible = !this.menuVisible;
  }

  hideMenu() {
    this.menuVisible = false;
    this.markRead.next();
  }
}
