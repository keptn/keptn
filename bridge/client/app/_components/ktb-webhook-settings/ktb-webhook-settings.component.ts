import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { AbstractControl, FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';
import { WebhookConfig } from '../../../../shared/models/webhook-config';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';
import { Secret } from '../../_models/secret';

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
  public sendFinishedOverlayConfig: DtOverlayConfig = {
    pinnable: true,
    originY: 'center',
  };
  public _eventType: string | undefined;

  @Input() public secrets: Secret[] | undefined;

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
}
