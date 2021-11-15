import { Component } from '@angular/core';
import { environment } from '../../environments/environment';

@Component({
  templateUrl: './not-found.component.html',
  styleUrls: ['./not-found.component.scss'],
})
export class NotFoundComponent {
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
}
