import { Component, OnInit } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { BehaviorSubject } from 'rxjs';
import { Secret } from '../../_models/secret';
import { DeleteDialogState } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';

@Component({
  selector: 'ktb-secrets-list',
  templateUrl: './ktb-secrets-list.component.html',
  styleUrls: [],
})
export class KtbSecretsListComponent implements OnInit {
  private _secrets = new BehaviorSubject<Secret[]>([]);
  secrets$ = this._secrets.asObservable();
  public currentSecret?: Secret;
  public deleteState: DeleteDialogState = null;

  constructor(private dataService: DataService) {
    this.deleteSecret.bind(this);
  }

  ngOnInit(): void {
    this.dataService.getSecrets().subscribe((secrets) => {
      this._secrets.next(secrets);
    });
  }

  public createDataSource(secrets: Secret[]): DtTableDataSource<Secret> {
    return new DtTableDataSource(secrets);
  }

  public toSecret(value: unknown): Secret {
    return <Secret>value;
  }

  public triggerDeleteSecret(secret: Secret): void {
    this.currentSecret = secret;
    this.deleteState = 'confirm';
  }

  public deleteSecret(secret?: Secret): void {
    if (!secret) {
      return;
    }
    this.dataService.deleteSecret(secret.name, secret.scope).subscribe(() => {
      this.deleteState = 'success';
      this._secrets.next(this._secrets.getValue().filter((s) => s.name !== secret.name));
    });
  }
}
