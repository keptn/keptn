import { ChangeDetectorRef, Component, Input } from '@angular/core';
import { DataService } from '../../../_services/data.service';
import { EndSessionData } from '../../../../../shared/interfaces/end-session-data';
import { Location } from '@angular/common';

@Component({
  selector: 'ktb-user[user]',
  templateUrl: './ktb-user.component.html',
  styleUrls: ['./ktb-user.component.scss'],
})
export class KtbUserComponent {
  @Input() user?: string;
  public logoutFormData: EndSessionData = {
    state: '',
    post_logout_redirect_uri: '',
    end_session_endpoint: '',
    id_token_hint: '',
  };

  constructor(
    private readonly dataService: DataService,
    private readonly _changeDetectorRef: ChangeDetectorRef,
    private readonly location: Location
  ) {}

  logout(submitEvent: { target: { submit: () => void } }): void {
    this.dataService.logout().subscribe((response) => {
      if (response) {
        this.logoutFormData = response;
        this._changeDetectorRef.detectChanges();
        submitEvent.target.submit();
      } else {
        window.location.assign(this.location.prepareExternalUrl('/logoutsession'));
      }
    });
  }
}
