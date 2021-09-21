import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { AbstractControl, FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';
import { WebhookConfig } from '../../../../shared/models/webhook-config';
import { combineLatest } from 'rxjs';

type ControlType = 'method' | 'url' | 'payload' | 'proxy' | 'header';

@Component({
  selector: 'ktb-webhook-settings',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss'],
})
export class KtbWebhookSettingsComponent implements OnInit {
  public _projectName?: string;
  public _subscription?: UniformSubscription;
  public webhookConfigForm = new FormGroup({
    method: new FormControl('', [Validators.required]),
    url: new FormControl('', [Validators.required, Validators.pattern(FormUtils.URL_PATTERN)]),
    payload: new FormControl('', []),
    header: new FormArray([]),
    proxy: new FormControl('', [Validators.pattern(FormUtils.URL_PATTERN)]),
  });

  public webhookMethods: WebhookConfigMethod[] = ['GET', 'POST', 'PUT'];
  private _webhook?: WebhookConfig;

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
    combineLatest([
      this.getFormControl('method').valueChanges,
      this.getFormControl('url').valueChanges,
      this.getFormControl('payload').valueChanges,
      this.getFormControl('proxy').valueChanges,
      this.getFormControl('header').valueChanges,
    ]).subscribe(([method, url, payload, proxy, header]) => {
      if (!this._webhook) {
        this._webhook = new WebhookConfig();
      }
      this._webhook.method = method;
      this._webhook.url = url;
      this._webhook.payload = payload;
      this._webhook.proxy = proxy;
      this._webhook.header = header;
      this.webhookChange.emit(this._webhook);
    });
  }

  public ngOnInit(): void {
    this.validityChanged.next(this.webhookConfigForm.valid);
  }

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
