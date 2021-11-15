import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment';

@Component({
  selector: 'ktb-not-found',
  templateUrl: './not-found.component.html',
  styleUrls: ['./not-found.component.scss'],
})
export class NotFoundComponent {
  public routerUrl: string;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;

  constructor(private router: Router) {
    this.routerUrl = this.router.url;
  }
}
