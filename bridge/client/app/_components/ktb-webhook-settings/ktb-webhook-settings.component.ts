import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { AbstractControl, FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';
import { WebhookConfig } from '../../../../shared/models/webhook-config';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';
import { Secret } from '../../_models/secret';
import { SelectTreeNode } from '../ktb-tree-list-select/ktb-tree-list-select.component';

type ControlType = 'method' | 'url' | 'payload' | 'proxy' | 'header' | 'sendFinished';

@Component({
  selector: 'ktb-webhook-settings',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss'],
})
export class KtbWebhookSettingsComponent implements OnInit {
  private _webhook: WebhookConfig = new WebhookConfig();
  public webhookConfigForm = new FormGroup({
    method: new FormControl('', [Validators.required]),
    url: new FormControl('', [Validators.required, FormUtils.isUrlValidator]),
    payload: new FormControl('', []),
    header: new FormArray([]),
    proxy: new FormControl('', [FormUtils.isUrlValidator]),
    sendFinished: new FormControl('true', []),
  });
  public webhookMethods: WebhookConfigMethod[] = ['GET', 'POST', 'PUT'];
  public secretDataSource: SelectTreeNode[] = [];
  public sendFinishedOverlayConfig: DtOverlayConfig = {
    pinnable: true,
    originY: 'center',
  };
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
  set webhook(webhookConfig: WebhookConfig | undefined) {
    if (webhookConfig && webhookConfig !== this._webhook) {
      this._webhook = webhookConfig;
      this.getFormControl('method').setValue(webhookConfig.method);
      this.getFormControl('url').setValue(webhookConfig.url);
      this.getFormControl('payload').setValue(webhookConfig.payload);
      this.getFormControl('proxy').setValue(webhookConfig.proxy);
      this.setSendFinishedControl();

      for (const header of webhookConfig.header || []) {
        this.addHeader(header.name, header.value);
      }

      for (const controlKey of Object.keys(this.webhookConfigForm.controls)) {
        this.webhookConfigForm.get(controlKey)?.markAsDirty();
      }
    }
  }

  @Input()
  set secrets(secrets: Secret[] | undefined) {
    if (secrets) {
      this.secretDataSource = secrets.map((secret: Secret) => this.mapSecret(secret));
    }
  }

  @Input()
  set eventPayload(event: Record<string, unknown> | undefined) {
    this.eventDataSource = event ? this.setObject(event) : undefined;
  }

  private setObject(data: Record<string, unknown>, path = '.event'): SelectTreeNode[] {
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
  @Output() webhookChange: EventEmitter<WebhookConfig> = new EventEmitter<WebhookConfig>();

  get header(): FormArray {
    return this.getFormControl('header') as FormArray;
  }

  get headerControls(): FormGroup[] {
    return this.header.controls as FormGroup[];
  }

  constructor() {
    this.webhookConfigForm.statusChanges.subscribe((status: 'INVALID' | 'VALID') => {
      this.validityChanged.next(status === 'VALID');
    });
  }

  public ngOnInit(): void {
    this.validityChanged.next(this.webhookConfigForm.valid);
  }

  public onWebhookFormChange(): void {
    this._webhook.method = this.getFormControl('method').value;
    this._webhook.url = this.getFormControl('url').value;
    this._webhook.payload = this.getFormControl('payload').value;
    this._webhook.proxy = this.getFormControl('proxy').value;
    this._webhook.header = this.getFormControl('header').value;
    this._webhook.sendFinished = this.getFormControl('sendFinished').value === 'true';
    this.webhookChange.emit(this._webhook);
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
      this.getFormControl('sendFinished').setValue(null);
      this.getFormControl('sendFinished').disable();
    } else {
      this.getFormControl('sendFinished').setValue(this._webhook.sendFinished.toString());
      this.getFormControl('sendFinished').enable();
    }
  }

  private mapSecret(secret: Secret): SelectTreeNode {
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
