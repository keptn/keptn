import {ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {DtSort, DtTableDataSource} from "@dynatrace/barista-components/table";
import {Observable, Subject} from "rxjs";
import {UniformRegistration} from "../../_models/uniform-registration";
import {takeUntil} from "rxjs/operators";
import {Secret} from "../../_models/secret";

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-secrets-list.component.html',
  styleUrls: ['./ktb-secrets-list.component.scss']
})
export class KtbSecretsListComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  private closeConfirmationDialogTimeout;

  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  public currentSecret: Secret;

  public deleteSecretDialogState: string | null;

  constructor(private dataService: DataService, private route: ActivatedRoute, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.dataService.getSecrets()
      .subscribe((secrets) => {
        this.tableEntries.data = secrets;
      });
  }

  triggerDeleteSecret(secret) {
    this.currentSecret = secret;
    clearTimeout(this.closeConfirmationDialogTimeout);
    this.deleteSecretDialogState = 'confirm';
  }

  deleteSecret(secret) {
    this.deleteSecretDialogState = 'deleting';
    this.dataService.deleteSecret(secret.name, secret.scope)
      .subscribe((result) => {
        this.deleteSecretDialogState = 'success';
        this.closeConfirmationDialogTimeout = setTimeout(() =>{
          this.closeConfirmationDialog();
        }, 2000);
        this.tableEntries.data = this.tableEntries.data.slice(this.tableEntries.data.indexOf(secret), 1);
      });
  }

  closeConfirmationDialog() {
    this.deleteSecretDialogState = null;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
