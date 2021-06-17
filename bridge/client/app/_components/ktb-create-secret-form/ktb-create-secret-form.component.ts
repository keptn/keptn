import {Component, OnInit} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {Secret} from "../../_models/secret";

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-create-secret-form.component.html',
  styleUrls: ['./ktb-create-secret-form.component.scss']
})
export class KtbCreateSecretFormComponent implements OnInit {

  public isLoading: Boolean = false;
  public secret: Secret = null;

  public defaultFormControls: {} = {
    name: new FormControl('', [Validators.required])
  };
  public createSecretForm = new FormGroup(this.defaultFormControls);

  constructor(private dataService: DataService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.secret = new Secret();
    this.addPair();
  }

  createSecret() {
    if(this.createSecretForm.valid) {
      this.isLoading = true;
      this.dataService.addSecret(this.secret)
        .subscribe((result) => {
          this.isLoading = false;
          this.router.navigate(['../'], { relativeTo: this.route });
        });
    }
  }

  addPair() {
    this.secret.addData();
    this.createSecretForm.addControl('key'+this.secret.data.length, new FormControl('', [Validators.required]));
    this.createSecretForm.addControl('value'+this.secret.data.length, new FormControl('', [Validators.required]));
  }

  removePair(index) {
    this.secret.removeData(index);
    this.createSecretForm.removeControl('key'+(index+1));
    this.createSecretForm.removeControl('value'+(index+1));
  }

}
