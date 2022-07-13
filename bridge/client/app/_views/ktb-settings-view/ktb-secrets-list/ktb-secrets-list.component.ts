import { Component, OnInit } from '@angular/core';
import { DataService } from '../../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { BehaviorSubject } from 'rxjs';
import { DeleteDialogState } from '../../../_components/_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';
import { IClientSecret } from '../../../../../shared/interfaces/secret';

@Component({
  selector: 'ktb-secrets-list',
  templateUrl: './ktb-secrets-list.component.html',
  styleUrls: [],
})
export class KtbSecretsListComponent implements OnInit {
  private _secrets = new BehaviorSubject<IClientSecret[]>([]);
  public secrets$ = this._secrets.asObservable();
  public currentSecret?: IClientSecret;
  public deleteState: DeleteDialogState = null;

  constructor(private dataService: DataService) {}

  ngOnInit(): void {
    this.dataService.getSecrets().subscribe((secrets) => {
      this._secrets.next(secrets);
    });
  }

  public createDataSource(secrets: IClientSecret[]): DtTableDataSource<IClientSecret> {
    return new DtTableDataSource(secrets);
  }

  public toSecret(value: unknown): IClientSecret {
    return <IClientSecret>value;
  }

  public triggerDeleteSecret(secret: IClientSecret): void {
    this.currentSecret = secret;
    this.deleteState = 'confirm';
  }

  public deleteSecret(secret?: IClientSecret): void {
    if (!secret) {
      return;
    }
    this.dataService.deleteSecret(secret.name, secret.scope).subscribe(() => {
      this.deleteState = 'success';
      this._secrets.next(this._secrets.getValue().filter((s) => s.name !== secret.name));
    });
  }
}
