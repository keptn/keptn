import {Component, Input} from '@angular/core';
import {Location} from '@angular/common';

@Component({
  selector: 'ktb-user',
  templateUrl: './ktb-user.component.html',
  styleUrls: ['./ktb-user.component.scss']
})
export class KtbUserComponent {

  @Input() user;

  constructor(private readonly location: Location) {
  }

  logout(): void {
    window.location.href = this.location.prepareExternalUrl('/logout');
  }
}
