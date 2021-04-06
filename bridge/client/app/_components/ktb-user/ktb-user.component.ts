import {Component} from '@angular/core';

@Component({
  selector: 'ktb-user',
  templateUrl: './ktb-user.component.html',
  styleUrls: ['./ktb-user.component.scss']
})
export class KtbUserComponent {

  logout(): void {
    window.location.href = '/logout';
  }

}
