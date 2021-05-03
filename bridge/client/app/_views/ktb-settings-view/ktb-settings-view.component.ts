import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subject} from "rxjs";
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../_services/data.service";
import {ActivatedRoute} from "@angular/router";
import {takeUntil} from "rxjs/operators";
import {DtToast} from "@dynatrace/barista-components/toast";

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss']
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private projectName: string;

  public gitUrlControl = new FormControl('', [Validators.required]);
  public gitUserControl = new FormControl('', [Validators.required]);
  public gitTokenControl = new FormControl('', [Validators.required]);
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl
  });
  public isGitUpstreamInProgress = false;

  constructor(private route: ActivatedRoute, private dataService: DataService, private toast: DtToast) {
  }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.dataService.getProject(params.projectName)
          .pipe(takeUntil(this.unsubscribe$))
          .subscribe(project => {
            this.projectName = project.projectName;
            if (project.gitRemoteURI) {
              this.gitUrlControl.setValue(project.gitRemoteURI);
            }
            if (project.gitUser) {
              this.gitUserControl.setValue(project.gitUser);
            }
            if (project.gitRemoteURI && project.gitUser) {
              this.gitTokenControl.setValue('***********************');
            }
          });
      })
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

  setGitUpstream() {
    this.isGitUpstreamInProgress = true;
    this.dataService.setGitUpstreamUrl(this.projectName, this.gitUrlControl.value, this.gitUserControl.value, this.gitTokenControl.value).subscribe(success => {
      this.isGitUpstreamInProgress = false;
      if (success) {
        this.toast.create('Git Upstream URL set successfully');
      } else {
        this.toast.create('Git Upstream URL could not be set');
      }
    });
  }

}
