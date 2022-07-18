import { Component, Input } from '@angular/core';
import { Location } from '@angular/common';

@Component({
  selector: 'ktb-project-create-message',
  templateUrl: './ktb-project-create-message.component.html',
})
export class KtbProjectCreateMessageComponent {
  @Input() projectName?: string;
  @Input() routerLink?: string;

  constructor(public location: Location) {}
}
