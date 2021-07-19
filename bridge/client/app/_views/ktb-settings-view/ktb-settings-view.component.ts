import {Component, OnDestroy, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {Subject} from 'rxjs';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {filter, map, switchMap, take, takeUntil} from 'rxjs/operators';
import {DtToast} from '@dynatrace/barista-components/toast';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss']
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private projectName: string;

  @ViewChild('deleteProjectDialog')
  private deleteProjectDialog: TemplateRef<MatDialog>;

  public gitUrlControl = new FormControl('', [Validators.required]);
  public gitUserControl = new FormControl('', [Validators.required]);
  public gitTokenControl = new FormControl('', [Validators.required]);
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl
  });
  public isGitUpstreamInProgress = false;

  public deletionConfirmationControl = new FormControl('');
  public deletionConfirmationForm = new FormGroup({
    deletionConfirmation: this.deletionConfirmationControl
  });
  public deletionDialogRef: MatDialogRef<any>;
  public isDeleteProjectInProgress = false;
  public deletionError = '';

  constructor(private route: ActivatedRoute,
              private dataService: DataService,
              private toast: DtToast,
              private dialog: MatDialog,
              private router: Router) {
  }

  ngOnInit(): void {
    this.route.params.pipe(
      map(params => params.projectName),
      switchMap(projectName => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$),
      filter(project => !!project)
    ).subscribe(project => {
      this.projectName = project.projectName;
      this.gitUrlControl.setValue(project.gitRemoteURI);
      this.gitUserControl.setValue(project.gitUser);

      if (project.gitRemoteURI && project.gitUser) {
        this.gitTokenControl.setValue('***********************');
      } else {
        this.gitTokenControl.setValue('');
      }

      this.deletionConfirmationControl.setValidators([Validators.required, Validators.pattern(this.projectName)]);
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  setGitUpstream() {
    this.isGitUpstreamInProgress = true;
    this.dataService.setGitUpstreamUrl(this.projectName, this.gitUrlControl.value, this.gitUserControl.value, this.gitTokenControl.value)
      .subscribe(success => {
        this.isGitUpstreamInProgress = false;
        if (success) {
          this.toast.create('Git Upstream URL set successfully');
        } else {
          this.toast.create('Git Upstream URL could not be set');
        }
      });
  }

  public openProjectDeletionDialog() {
    this.deletionDialogRef = this.dialog.open(this.deleteProjectDialog, {
      data: {projectName: this.projectName},
      autoFocus: false
    });
    this.deletionDialogRef.afterClosed().subscribe(() => {
      this.deletionConfirmationControl.setValue('');
      this.deletionConfirmationForm.markAsUntouched();
      this.deletionConfirmationForm.updateValueAndValidity();
    });
  }

  public deleteProject() {
    this.isDeleteProjectInProgress = true;
    this.deletionError = '';
    this.dataService.projects
      .pipe(take(1))
      .subscribe(() => {
        this.router.navigate(['/', 'dashboard']);
      });

    this.dataService.deleteProject(this.projectName)
      .pipe(take(1))
      .subscribe(() => {
        this.deletionDialogRef.close();
        this.dataService.loadProjects();
      }, (err) => {
        this.isDeleteProjectInProgress = false;
        this.deletionError = 'Project could not be deleted: ' + err.message;
      });
  }
}
