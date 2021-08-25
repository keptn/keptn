import { Component } from '@angular/core';
import { NotificationsService } from '../../_services/notifications.service';

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: [],
  providers: [NotificationsService],
})
export class KtbSettingsViewComponent {
}
