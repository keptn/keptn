import { Component, OnInit } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { NotificationType } from '../../_models/notification';
import { NotificationsService } from '../../_services/notifications.service';
import { BehaviorSubject, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { IServiceSecret } from '../../../../shared/interfaces/secret';
import { addData } from '../../_models/secret';

const secretNamePattern = '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*';
const secretKeyPattern = '[-._a-zA-Z0-9]+';

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-create-secret-form.component.html',
  styleUrls: ['./ktb-create-secret-form.component.scss'],
})
export class KtbCreateSecretFormComponent implements OnInit {
  FormGroupClass = FormGroup;
  scopeControl = new FormControl(undefined, [Validators.required]);

  nameControl = new FormControl('', [
    Validators.required,
    Validators.pattern(secretNamePattern),
    Validators.maxLength(253),
  ]);

  dataControl = this.fb.array([
    this.fb.group({
      key: ['', [Validators.required, Validators.pattern(secretKeyPattern), Validators.maxLength(253)]],
      value: ['', [Validators.required]],
    }),
  ]);

  createSecretForm = this.fb.group({
    name: this.nameControl,
    scope: this.scopeControl,
    data: this.dataControl,
  });

  _scopes = new BehaviorSubject<string[]>([]);
  scopes$ = this._scopes.asObservable();
  isUpdating = false;
  isLoading = false;

  constructor(
    private dataService: DataService,
    private router: Router,
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private notificationService: NotificationsService
  ) {}

  ngOnInit(): void {
    this.loadSecretScopes();
  }

  private setLoadingState(isLoading: boolean): void {
    this.isLoading = isLoading;
    if (isLoading) {
      this.scopeControl.disable();
    } else {
      this.scopeControl.enable();
    }
  }

  private loadSecretScopes(): void {
    this.setLoadingState(true);
    this.dataService
      .getSecretScopes()
      .pipe(finalize(() => this.setLoadingState(false)))
      .subscribe((scopes) => this._scopes.next(scopes));
  }

  public createSecret(): void {
    this.isUpdating = true;
    const secret: IServiceSecret = {
      name: this.nameControl.value,
      scope: this.scopeControl.value,
    };
    for (const dataGroup of this.dataControl.controls) {
      addData(secret, dataGroup.get('key')?.value, dataGroup.get('value')?.value);
    }

    this.dataService
      .addSecret(secret)
      .pipe(
        map(() => true),
        catchError((err) => {
          console.log(err);
          if (err.status !== 409) {
            return of(false);
          }
          const message = `A secret with the name ${secret.name} already exists. Please use another name for this secret to continue.`;
          this.notificationService.addNotification(NotificationType.ERROR, message);
          return of(false);
        }),
        finalize(() => (this.isUpdating = false))
      )
      .subscribe((success) => {
        if (success) {
          this.router.navigate(['../'], { relativeTo: this.route });
        }
      });
  }

  public addPair(): void {
    this.dataControl.push(
      this.fb.group({
        key: ['', [Validators.required, Validators.pattern(secretKeyPattern), Validators.maxLength(253)]],
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
}
