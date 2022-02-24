import { Component, OnDestroy } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Secret } from '../../_models/secret';
import { NotificationType } from '../../_models/notification';
import { NotificationsService } from '../../_services/notifications.service';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-create-secret-form.component.html',
  styleUrls: ['./ktb-create-secret-form.component.scss'],
})
export class KtbCreateSecretFormComponent implements OnDestroy {
  private secretNamePattern = '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*';
  private secretKeyPattern = '[-._a-zA-Z0-9]+';
  private _isLoading = false;
  private unsubscribe$ = new Subject<void>();
  public FormGroupClass = FormGroup;
  public scopeControl = new FormControl(undefined, [Validators.required]);
  public nameControl = new FormControl('', [
    Validators.required,
    Validators.pattern(this.secretNamePattern),
    Validators.maxLength(253),
  ]);
  public dataControl = this.fb.array([
    this.fb.group({
      key: ['', [Validators.required, Validators.pattern(this.secretKeyPattern), Validators.maxLength(253)]],
      value: ['', [Validators.required]],
    }),
  ]);

  public isUpdating = false;
  public scopes?: string[];
  public createSecretForm = this.fb.group({
    name: this.nameControl,
    scope: this.scopeControl,
    data: this.dataControl,
  });

  public get isLoading(): boolean {
    return this._isLoading;
  }
  public set isLoading(isLoading: boolean) {
    this._isLoading = isLoading;
    if (isLoading) {
      this.scopeControl.disable();
    } else {
      this.scopeControl.enable();
    }
  }

  constructor(
    private dataService: DataService,
    private router: Router,
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private notificationService: NotificationsService
  ) {
    this.getSecretScopes();
  }

  private getSecretScopes(): void {
    this.isLoading = true;
    this.dataService.getSecretScopes().subscribe(
      (scopes) => {
        this.scopes = scopes;
        this.isLoading = false;
      },
      () => {
        this.isLoading = false;
      }
    );
  }

  public createSecret(): void {
    if (this.createSecretForm.valid) {
      this.isUpdating = true;

      const secret: Secret = new Secret();
      secret.setName(this.nameControl.value);
      secret.setScope(this.scopeControl.value);
      for (const dataGroup of this.dataControl.controls) {
        secret.addData(dataGroup.get('key')?.value, dataGroup.get('value')?.value);
      }

      this.dataService.addSecret(secret).subscribe(
        () => {
          this.isUpdating = false;
          this.router.navigate(['../'], { relativeTo: this.route });
        },
        (err) => {
          if (err.status === 409) {
            this.notificationService.addNotification(
              NotificationType.ERROR,
              `A secret with the name ${secret.name} already exists. Please use another name for this secret to continue.`
            );
          }
          this.isUpdating = false;
        }
      );
    }
  }

  public addPair(): void {
    this.dataControl.push(
      this.fb.group({
        key: ['', [Validators.required, Validators.pattern(this.secretKeyPattern), Validators.maxLength(253)]],
        value: ['', [Validators.required]],
      })
    );
  }

  public removePair(index: number): void {
    this.dataControl.removeAt(index);
  }

  public isFormValid(): boolean {
    return this.createSecretForm.valid && !this.isUpdating && !this.isLoading;
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
