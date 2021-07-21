import {Component, OnInit} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import {Secret} from '../../_models/secret';

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-create-secret-form.component.html',
  styleUrls: ['./ktb-create-secret-form.component.scss']
})
export class KtbCreateSecretFormComponent implements OnInit {
  public isLoading = false;
  public secret: Secret = new Secret();
  private secretNamePattern = '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*';
  private secretKeyPattern = '[-._a-zA-Z0-9]+';

  public defaultFormControls: {} = {
    name: new FormControl('', [Validators.required, Validators.pattern(this.secretNamePattern)])
  };
  public createSecretForm = new FormGroup(this.defaultFormControls);

  constructor(private dataService: DataService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.secret = new Secret();
    this.addPair();
  }

  public createSecret() {
    if (this.createSecretForm.valid) {
      this.isLoading = true;
      this.dataService.addSecret(this.secret)
        .subscribe((result) => {
          this.isLoading = false;
          this.router.navigate(['../'], { relativeTo: this.route });
        }, (err) => {
          this.isLoading = false;
        });
    }
  }

  public addPair() {
    this.secret.addData();
    this.createSecretForm.addControl('key' + this.secret.data.length,
                                    new FormControl('', [Validators.required, Validators.pattern(this.secretKeyPattern)]));
    this.createSecretForm.addControl('value' + this.secret.data.length,
                                    new FormControl('', [Validators.required]));
  }

  public removePair(index: number) {
    this.secret.removeData(index);
    this.createSecretForm.removeControl('key' + (index + 1));
    this.createSecretForm.removeControl('value' + (index + 1));
  }

  public getFormControl(controlName: string): AbstractControl | null {
    return this.createSecretForm.get(controlName);
  }

}
