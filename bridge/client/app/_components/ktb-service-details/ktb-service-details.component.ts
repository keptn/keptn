import {Location} from '@angular/common';
import {ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {Deployment} from '../../_models/deployment';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {takeUntil} from 'rxjs/operators';
import {Subject} from 'rxjs';
import {Sequence} from '../../_models/sequence';

@Component({
  selector: 'ktb-service-details',
  templateUrl: './ktb-service-details.component.html',
  styleUrls: ['./ktb-service-details.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbServiceDetailsComponent implements OnInit, OnDestroy{
  private _deployment: Deployment;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  public projectName: string;
  public selectedStage: string;
  public view: string;
  public remediations: {stage: string, remediations: Sequence[]};

  get deployment() {
    return this._deployment;
  }

  set deployment(deployment: Deployment) {
    if (this._deployment !== deployment) {
      const selectLast = !!this._deployment;
      this._deployment = deployment;
      if (!this._deployment.sequence) {
        this.loadSequence(selectLast);
      } else {
        this.selectLastStage();
      }
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, private router: Router, private location: Location) {

  }

  ngOnInit(): void {
    this.route.params.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(params => {
      this.projectName = params.projectName;
      this.selectedStage = params.stage;
      this.view = this.route.snapshot.url[this.route.snapshot.url.length - 1].path === 'remediation' ? 'remediation' : 'evaluation';
      this._changeDetectorRef.markForCheck();
    });
  }

  private loadSequence(selectLast: boolean) {
    this.dataService.getRoot(this.projectName, this.deployment.shkeptncontext).subscribe(sequence => {
      this.deployment.sequence = sequence;
      if (this.route.snapshot.url[this.route.snapshot.url.length - 1].path === 'remediation') {
        this.remediations = {stage: this.selectedStage, remediations: this.deployment.getStage(this.selectedStage).remediations};
      }
      if (selectLast || !this.selectedStage) {
        this.selectLastStage();
      }
      this._changeDetectorRef.markForCheck();
    });
  }

  private selectLastStage() {
    const stages = this.deployment.sequence.getStages();
    this.selectStage(stages[stages.length - 1]);
  }

  public selectStage(stageName: string) {
    this.selectedStage = stageName;
    const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', this.deployment.service, 'context', this.deployment.sequence.shkeptncontext, 'stage', stageName]);
    this.location.go(routeUrl.toString());
    this.view = 'evaluation';
    this._changeDetectorRef.markForCheck();
  }

  public selectRemediationView(stage: string) {
    this.selectedStage = stage;
    this.view = 'remediation';
    const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', this.deployment.service, 'context', this.deployment.sequence.shkeptncontext, 'stage', stage, 'remediation']);
    this.location.go(routeUrl.toString());
    this._changeDetectorRef.markForCheck();
  }

  public isUrl(value: string): boolean {
    try {
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
