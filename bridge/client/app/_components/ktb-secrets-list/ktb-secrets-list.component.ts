import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute } from '@angular/router';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { Subject } from 'rxjs';
import { Secret } from '../../_models/secret';
import { DeleteDialogState } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';

@Component({
  selector: 'ktb-secrets-view',
  templateUrl: './ktb-secrets-list.component.html',
  styleUrls: ['./ktb-secrets-list.component.scss'],
})
export class KtbSecretsListComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public tableEntries: DtTableDataSource<Secret> = new DtTableDataSource();
  public currentSecret?: Secret;
  public SecretClass = Secret;
  public deleteSecretDialogState: DeleteDialogState = null;
  public deleteState: DeleteDialogState = null;

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private _changeDetectorRef: ChangeDetectorRef
  ) {
    this.deleteSecret.bind(this);
  }

  ngOnInit(): void {
    this.dataService.getSecrets().subscribe((secrets) => {
      this.tableEntries.data = secrets;
    });
  }

  public triggerDeleteSecret(secret: Secret): void {
    this.currentSecret = secret;
    this.deleteState = 'confirm';
  }

  public deleteSecret(secret?: Secret): void {
    if (secret) {
      this.dataService.deleteSecret(secret.name, secret.scope).subscribe((result) => {
        this.deleteState = 'success';
        const data: Secret[] = this.tableEntries.data;
        data.splice(
          data.findIndex((s: Secret) => s.name === secret.name),
          1
        );
        this.tableEntries = new DtTableDataSource(data);
      });
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
