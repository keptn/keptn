import {Component, Input} from '@angular/core';
import {UniformRegistrationLog} from "../../_models/uniform-registration-log";

@Component({
  selector: 'ktb-uniform-registration-logs',
  templateUrl: './ktb-uniform-registration-logs.component.html',
  styleUrls: ['./ktb-uniform-registration-logs.component.scss']
})
export class KtbUniformRegistrationLogsComponent {
  @Input() logs: UniformRegistrationLog[];
  @Input() projectName: string;
}
