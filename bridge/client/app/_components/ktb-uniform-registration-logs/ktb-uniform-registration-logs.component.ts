import {Component, Input, OnInit} from '@angular/core';
import {UniformRegistrationLog} from "../../_models/uniform-registration-log";

@Component({
  selector: 'ktb-uniform-registration-logs',
  templateUrl: './ktb-uniform-registration-logs.component.html',
  styleUrls: ['./ktb-uniform-registration-logs.component.scss']
})
export class KtbUniformRegistrationLogsComponent implements OnInit {

  @Input() logs: UniformRegistrationLog[];
  @Input() projectName: string;

  constructor() { }

  ngOnInit(): void {
  }

}
