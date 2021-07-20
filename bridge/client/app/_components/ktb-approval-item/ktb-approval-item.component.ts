import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { Trace } from '../../_models/trace';
import { Observable, of } from 'rxjs';
import {Project} from '../../_models/project';
import {DataService} from '../../_services/data.service';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';

@Component({
  selector: 'ktb-approval-item[event]',
  templateUrl: './ktb-approval-item.component.html',
  styleUrls: ['./ktb-approval-item.component.scss'],
})
export class KtbApprovalItemComponent implements OnInit{

  public project$: Observable<Project | undefined> = of(undefined);
  public _event?: Trace;
  public approvalResult?: boolean;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  @Input() isSequence = false;

  @Input()
  get event(): Trace | undefined {
    return this._event;
  }

  set event(value: Trace | undefined) {
    if (this._event !== value) {
      this._event = value;
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(private changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
  }

  ngOnInit(): void {
    if (this.event?.project) {
      this.project$ = this.dataService.getProject(this.event?.project);
    }
  }

  public handleApproval(approval: Trace, result: boolean) {
    this.dataService.sendApprovalEvent(approval, result);
    this.approvalResult = result;
  }
}
