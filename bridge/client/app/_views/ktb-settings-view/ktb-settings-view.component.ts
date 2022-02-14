import { Component } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss'],
  providers: [],
})
export class KtbSettingsViewComponent {
  public hasUnreadLogs$: Observable<boolean>;

  constructor(dataService: DataService) {
    this.hasUnreadLogs$ = dataService.hasUnreadUniformRegistrationLogs;
  }
}
