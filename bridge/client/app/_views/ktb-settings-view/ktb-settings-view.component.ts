import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subject} from 'rxjs';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {filter, map, switchMap, takeUntil} from 'rxjs/operators';
import {DtToast} from '@dynatrace/barista-components/toast';
import { Project } from '../../_models/project';

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: []
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private projectName?: string;

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
    this.route.params.pipe(
      map(params => params.projectName),
      switchMap(projectName => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$),
      filter((project: Project | undefined): project is Project => !!project)
    ).subscribe(project => {
      this.projectName = project.projectName;
      this.gitUrlControl.setValue(project.gitRemoteURI);
      this.gitUserControl.setValue(project.gitUser);

      if (project.gitRemoteURI && project.gitUser) {
        this.gitTokenControl.setValue('***********************');
      } else {
        this.gitTokenControl.setValue('');
      }
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

  setGitUpstream() {
    if (this.projectName) {
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

}
