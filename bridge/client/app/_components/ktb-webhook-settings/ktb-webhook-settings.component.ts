import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { AbstractControl, FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';
import { WebhookConfig } from '../../../../shared/models/webhook-config';

type ControlType = 'method' | 'url' | 'payload' | 'proxy' | 'header';

@Component({
  selector: 'ktb-webhook-settings',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss'],
})
export class KtbWebhookSettingsComponent implements OnInit {
  private _webhook: WebhookConfig = new WebhookConfig();
  public _projectName?: string;
  public _subscription?: UniformSubscription;
  public webhookConfigForm = new FormGroup({
    method: new FormControl('', [Validators.required]),
    url: new FormControl('', [Validators.required, FormUtils.urlValidator]),
    payload: new FormControl('', []),
    header: new FormArray([]),
    proxy: new FormControl('', [FormUtils.urlValidator]),
  });

  public webhookMethods: WebhookConfigMethod[] = ['GET', 'POST', 'PUT'];

  @Input()
  set webhook(webhookConfig: WebhookConfig | undefined) {
    if (webhookConfig && webhookConfig !== this._webhook) {
      this._webhook = webhookConfig;
      this.getFormControl('method').setValue(webhookConfig.method);
      this.getFormControl('url').setValue(webhookConfig.url);
      this.getFormControl('payload').setValue(webhookConfig.payload);
      this.getFormControl('proxy').setValue(webhookConfig.proxy);

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
    this._webhook.url =  this.getFormControl('url').value;
    this._webhook.payload = this.getFormControl('payload').value;
    this._webhook.proxy = this.getFormControl('proxy').value;
    this._webhook.header = this.getFormControl('header').value;
    this.webhookChange.emit(this._webhook);
  }

  public addHeader(name?: string, value?: string): void {
    this.header.push(new FormGroup({
      name: new FormControl(name, [Validators.required]),
      value: new FormControl(value, [Validators.required]),
    }));
  }

  public removeHeader(index: number): void {
    this.header.removeAt(index);
  }

  public getFormControl(controlName: ControlType): AbstractControl {
    return this.webhookConfigForm.get(controlName) as AbstractControl;
  }

}
