import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {DtTableDataSource} from '@dynatrace/barista-components/table';
import {Subject} from 'rxjs';
import {Secret} from '../../_models/secret';

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-secrets-list.component.html',
  styleUrls: ['./ktb-secrets-list.component.scss']
})
export class KtbSecretsListComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  private closeConfirmationDialogTimeout?: ReturnType<typeof setTimeout>;

  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  public currentSecret?: Secret;

  public deleteSecretDialogState: string | null = null;

  constructor(private dataService: DataService, private route: ActivatedRoute, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.dataService.getSecrets()
      .subscribe((secrets) => {
        this.tableEntries.data = secrets;
      });
  }

  public triggerDeleteSecret(secret: Secret) {
    this.currentSecret = secret;
    if (this.closeConfirmationDialogTimeout) {
      clearTimeout(this.closeConfirmationDialogTimeout);
    }
    this.deleteSecretDialogState = 'confirm';
  }

  public deleteSecret(secret: Secret) {
    this.deleteSecretDialogState = 'deleting';
    this.dataService.deleteSecret(secret.name, secret.scope)
      .subscribe((result) => {
        this.deleteSecretDialogState = 'success';
        this.closeConfirmationDialogTimeout = setTimeout(() => {
          this.closeConfirmationDialog();
        }, 2000);
        this.tableEntries.data = this.tableEntries.data.slice(this.tableEntries.data.indexOf(secret), 1);
      });
  }

  closeConfirmationDialog() {
    this.deleteSecretDialogState = null;
  }

  public toSecret(row: Secret): Secret {
    return row;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
