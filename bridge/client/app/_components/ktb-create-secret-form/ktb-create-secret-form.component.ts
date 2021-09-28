import { Component } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AbstractControl, FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Scopes, Secret } from '../../_models/secret';

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-create-secret-form.component.html',
  styleUrls: ['./ktb-create-secret-form.component.scss']
})
export class KtbCreateSecretFormComponent {

  public isLoading = false;

  private secretNamePattern = '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*';
  private secretKeyPattern = '[-._a-zA-Z0-9]+';

  public scopes = Object.values(Scopes);

  public createSecretForm = this.fb.group({
    name: ['', [Validators.required, Validators.pattern(this.secretNamePattern)]],
    scope: [this.scopes[0], [Validators.required]],
    data: this.fb.array([
      this.fb.group({
        key: ['', [Validators.required, Validators.pattern(this.secretKeyPattern)]],
        value: ['', [Validators.required]]
      })
    ])
  });

  constructor(private dataService: DataService, private router: Router, private route: ActivatedRoute, private fb: FormBuilder) {
  }

  get data(): FormArray | null {
    return this.createSecretForm.get('data') as FormArray;
  }

  get dataControls(): FormGroup[] {
    return this.data?.controls as FormGroup[] || [];
  }

  public createSecret(): void {
    if (this.createSecretForm.valid) {
      this.isLoading = true;

      const secret: Secret = new Secret();
      secret.setName(this.createSecretForm.get('name')?.value);
      secret.setScope(this.createSecretForm.get('scope')?.value);
      for (const dataGroup of this.data?.controls || []) {
        secret.addData(dataGroup.get('key')?.value, dataGroup.get('value')?.value);
      }

      this.dataService.addSecret(secret)
        .subscribe((result) => {
          this.isLoading = false;
          this.router.navigate(['../'], {relativeTo: this.route});
        }, (err) => {
          this.isLoading = false;
        });
    }
  }

  public addPair(): void {
    this.data?.push(this.fb.group({
      key: ['', [Validators.required, Validators.pattern(this.secretKeyPattern)]],
      value: ['', [Validators.required]]
    }));
  }

  public removePair(index: number): void {
    this.data?.removeAt(index);
  }

  public getFormControl(controlName: string): AbstractControl | null {
    return this.createSecretForm.get(controlName);
  }

}
