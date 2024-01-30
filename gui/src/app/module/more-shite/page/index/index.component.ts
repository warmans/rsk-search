import {Component, OnInit} from '@angular/core';
import {Title} from "@angular/platform-browser";

@Component({
  selector: 'app-index',
  standalone: false,
  templateUrl: './index.component.html',
  styleUrl: './index.component.scss'
})
export class IndexComponent implements OnInit {

  constructor(private titleService: Title) {

  }

  ngOnInit(): void {
    this.titleService.setTitle('More Shoddy Shite');
  }

}
