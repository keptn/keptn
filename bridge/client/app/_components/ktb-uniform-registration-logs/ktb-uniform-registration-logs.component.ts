import { Component, Input } from '@angular/core';
import { UniformRegistrationLog } from '../../../../server/interfaces/uniform-registration-log';
import { DataService } from '../../_services/data.service';

@Component({
  selector: 'ktb-uniform-registration-logs',
  templateUrl: './ktb-uniform-registration-logs.component.html',
  styleUrls: ['./ktb-uniform-registration-logs.component.scss']
})
export class KtbUniformRegistrationLogsComponent {
  @Input() logs: UniformRegistrationLog[] = [];
  @Input() projectName?: string;
  @Input() lastSeen?: string;

  public isUnread(time: string): boolean {
    return this.lastSeen ? new Date(this.lastSeen) < new Date(time) : true;
  }
}
