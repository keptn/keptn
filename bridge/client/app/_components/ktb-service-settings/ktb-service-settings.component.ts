import { Component } from '@angular/core';
import { NotificationsService } from '../../_services/notifications.service';

@Component({
  selector: 'ktb-service-settings',
  templateUrl: './ktb-service-settings.component.html',
  providers: [NotificationsService],
})
export class KtbServiceSettingsComponent {
}
