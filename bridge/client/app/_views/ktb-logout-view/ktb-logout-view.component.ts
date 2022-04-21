import { Component } from '@angular/core';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'ktb-logout-view',
  templateUrl: './ktb-logout-view.component.html',
  styleUrls: ['./ktb-logout-view.component.scss'],
})
export class KtbLogoutViewComponent {
  public logoUrl = environment.config.logoInvertedUrl;
}
