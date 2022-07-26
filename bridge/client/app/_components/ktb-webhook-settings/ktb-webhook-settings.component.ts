import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { AbstractControl, FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { IWebhookConfigClient, WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';
import { SelectTreeNode } from '../ktb-tree-list-select/ktb-tree-list-select.component';
import { IClientSecret } from '../../../../shared/interfaces/secret';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';

type ControlType = 'method' | 'url' | 'payload' | 'proxy' | 'header' | 'sendFinished' | 'sendStarted';

@Component({
  selector: 'ktb-webhook-settings',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss'],
})
export class KtbWebhookSettingsComponent implements OnInit {
  private _webhook?: IWebhookConfigClient;
  public webhookConfigForm = new FormGroup({
    method: new FormControl('', [Validators.required]),
    url: new FormControl('', [Validators.required, FormUtils.isUrlOrSecretValidator]),
    payload: new FormControl('', [FormUtils.payloadSpecialCharValidator]),
    header: new FormArray([]),
    proxy: new FormControl('', [FormUtils.isUrlValidator]),
    sendFinished: new FormControl('true'),
    sendStarted: new FormControl('true'),
  });
  public webhookMethods: WebhookConfigMethod[] = ['GET', 'POST', 'PUT'];
  public secretDataSource: SelectTreeNode[] = [];
  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
    originY: 'center',
  };
  public projectName$ = this.route.paramMap.pipe(map((params) => params.get('projectName')));

  public _eventType?: string;
  public eventDataSource?: SelectTreeNode[];

  @Input()
  set eventType(eventType: string | undefined) {
    if (this._eventType != eventType) {
      this._eventType = eventType;
      this.setSendFinishedControl();
    }
  }

  @Input()
  set webhook(webhookConfig: IWebhookConfigClient | undefined) {
    if (webhookConfig && webhookConfig !== this._webhook) {
      this._webhook = webhookConfig;
      this.getFormControl('method').setValue(webhookConfig.method);
      this.getFormControl('url').setValue(webhookConfig.url);
      this.getFormControl('payload').setValue(webhookConfig.payload);
      this.getFormControl('proxy').setValue(webhookConfig.proxy);
      this.setSendFinishedControl();

      for (const header of webhookConfig.header || []) {
        this.addHeader(header.key, header.value);
      }

      for (const controlKey of Object.keys(this.webhookConfigForm.controls)) {
        this.webhookConfigForm.get(controlKey)?.markAsDirty();
      }
    }
  }

  @Input()
  set secrets(secrets: IClientSecret[] | undefined) {
    if (secrets) {
      this.secretDataSource = secrets.map((secret: IClientSecret) => this.mapSecret(secret));
    }
  }

  @Input()
  set eventPayload(event: Record<string, unknown> | undefined) {
    this.eventDataSource = event ? this.setObject(event) : undefined;
  }

  private setObject(data: Record<string, unknown>, path = ''): SelectTreeNode[] {
    const result: SelectTreeNode[] = [];
    for (const key of Object.keys(data)) {
      const newItem = this.generateNewTreeNode(data[key], key, `${path}.${key}`);
      result.push(newItem);
    }
    return result;
  }

  private generateNewTreeNode(property: unknown, itemName: string, itemPath: string): SelectTreeNode {
    const newItem: SelectTreeNode = {
      name: itemName,
    };
    if (property instanceof Array) {
      newItem.keys = this.setArray(property, itemPath);
    } else if (property && typeof property === 'object') {
      newItem.keys = this.setObject(property as Record<string, unknown>, itemPath);
    } else {
      newItem.path = itemPath;
    }
    return newItem;
  }

  private setArray(array: Array<unknown>, path: string): SelectTreeNode[] {
    const result: SelectTreeNode[] = [];
    const data = array;
    for (let i = 0; i < data.length; ++i) {
      const newItem = this.generateNewTreeNode(data[i], `[${i}]`, `(index ${path} ${i})`);
      result.push(newItem);
    }
    return result;
  }

  @Output() validityChanged: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() webhookChange: EventEmitter<IWebhookConfigClient> = new EventEmitter<IWebhookConfigClient>();
  @Output() webhookFormDirty: EventEmitter<boolean> = new EventEmitter<boolean>();

  get header(): FormArray {
    return this.getFormControl('header') as FormArray;
  }

  get headerControls(): FormGroup[] {
    return this.header.controls as FormGroup[];
  }

  constructor(private route: ActivatedRoute) {
    this.webhookConfigForm.statusChanges.subscribe((status: 'INVALID' | 'VALID') => {
      this.validityChanged.next(status === 'VALID');
    });
  }

  public ngOnInit(): void {
    this.validityChanged.next(this.webhookConfigForm.valid);
  }

  public onWebhookFormChange(): void {
    this._webhook = {
      method: this.getFormControl('method').value,
      url: this.getFormControl('url').value,
      payload: this.getFormControl('payload').value,
      proxy: this.getFormControl('proxy').value,
      header: this.getFormControl('header').value,
      sendFinished: this.getFormControl('sendFinished').value === 'true',
      sendStarted: this.getFormControl('sendStarted').value === 'true',
      type: this._webhook?.type ?? '',
    };
    this.webhookChange.emit(this._webhook);
    this.webhookFormDirty.emit(this.webhookConfigForm.dirty);
  }

  public addHeader(name?: string, value?: string): void {
    this.header.push(
      new FormGroup({
        name: new FormControl(name || '', [Validators.required]),
        value: new FormControl(value || '', [Validators.required]),
      })
    );
  }

  public removeHeader(index: number): void {
    this.header.removeAt(index);
  }

  public getFormControl(controlName: ControlType, controlIndex?: number): AbstractControl {
    if (controlName === 'header' && controlIndex !== undefined) {
      const group = this.header.at(controlIndex) as FormGroup;
      return group.controls.value;
    } else {
      return this.webhookConfigForm.get(controlName) as AbstractControl;
    }
  }

  private setSendFinishedControl(): void {
    if (this._eventType !== 'triggered' && this._eventType !== '>') {
      this.disableFormControl('sendStarted');
      this.disableFormControl('sendFinished');
    } else {
      this.enableFormControl('sendStarted', (this._webhook?.sendStarted ?? true).toString());
      this.enableFormControl('sendFinished', (this._webhook?.sendFinished ?? true).toString());
    }
  }

  private disableFormControl(controlName: ControlType): void {
    this.getFormControl(controlName).setValue(null);
    this.getFormControl(controlName).disable();
  }

  private enableFormControl(controlName: ControlType, initValue?: string): void {
    this.getFormControl(controlName).setValue(initValue);
    this.getFormControl(controlName).enable();
  }

  private mapSecret(secret: IClientSecret): SelectTreeNode {
    const scrt: SelectTreeNode = { name: secret.name };
    if (secret.keys) {
      scrt.keys = secret.keys.map((key: string) => {
        return { name: key, path: `.secret.${secret.name}.${key}` };
      });
      scrt.keys.sort((a, b) => a.name.localeCompare(b.name));
    }
    return scrt;
  }
}
