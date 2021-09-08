import { Component, Input } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DataService } from '../../_services/data.service';
import { FormUtils } from '../../_utils/form.utils';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { UniformSubscriptionFilter } from '../../../../shared/interfaces/uniform-subscription';
import { AppUtils } from '../../_utils/app.utils';
import { WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';

type ControlType = 'method' | 'url' | 'payload' | 'proxy' | 'header';

@Component({
  selector: 'ktb-webhook-settings[subscriptionExists][projectName][subscription]',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss'],
})
export class KtbWebhookSettingsComponent {

  public _projectName?: string;
  public _subscription?: UniformSubscription;
  private _prevFilter?: UniformSubscriptionFilter;
  public webhookConfigForm = this.formBuilder.group({
    method: ['', [Validators.required]],
    url: ['', [Validators.required, Validators.pattern(FormUtils.URL_PATTERN)]],
    payload: ['', [Validators.required]],
    header: this.formBuilder.array([]),
    proxy: ['', [Validators.pattern(FormUtils.URL_PATTERN)]],
  });

  public webhookMethods: WebhookConfigMethod[] = ['POST', 'PUT'];
  private _subscriptionExists = false;
  public loading = false;

  @Input() set subscriptionExists(status: boolean) {
    if (this._subscriptionExists !== status) {
      this._subscriptionExists = status;
      this.getWebhook();
    }
  }

  @Input()
  get projectName(): string {
    return this._projectName || '';
  }

  set projectName(projectName: string) {
    if (this._projectName !== projectName) {
      this._projectName = projectName;
      this.getWebhook();
    }
  }

  @Input()
  get subscription(): UniformSubscription | undefined {
    return this._subscription;
  }

  set subscription(value: UniformSubscription | undefined) {
    if (this._subscription !== value) {
      this._subscription = value;
      this.prevFilter = this.subscription?.filter;
      this.getWebhook();
    }
  }

  get prevFilter(): UniformSubscriptionFilter | undefined {
    return this._prevFilter;
  }

  set prevFilter(value: UniformSubscriptionFilter | undefined) {
    if (this._prevFilter !== value) {
      this._prevFilter = AppUtils.copyObject(value);
    }
  }

  get header(): FormArray {
    return this.getFormControl('header') as FormArray;
  }

  get headerControls(): FormGroup[] {
    return this.header.controls as FormGroup[] || [];
  }

  constructor(private dataService: DataService, private formBuilder: FormBuilder) {
  }

  private getWebhook(): void {
    if (this._subscriptionExists && this._subscription && this._projectName) {
      this.loading = true;
      const stage: string | undefined = this._subscription.filter?.stages?.[0];
      const services: string | undefined = this._subscription.filter?.services?.[0];
      this.dataService.getWebhookConfig(this._projectName, stage, services)
        .subscribe(webhookConfig => {
          this.getFormControl('method').setValue(webhookConfig.method);
          this.getFormControl('url').setValue(webhookConfig.url);
          this.getFormControl('payload').setValue(webhookConfig.payload);
          this.getFormControl('proxy').setValue(webhookConfig.proxy);

          for (const header of webhookConfig.header || []) {
            this.addHeader(header.name, header.value);
          }

          this.loading = false;
        }, () => {
          this.loading = false;
        });
    }
  }

  public addHeader(name?: string, value?: string): void {
    this.header.push(this.formBuilder.group({
      name: [name, [Validators.required]],
      value: [value, [Validators.required]],
    }));
  }

  public removeHeader(index: number): void {
    this.header.removeAt(index);
  }

  public getFormControl(controlName: ControlType): AbstractControl {
    return this.webhookConfigForm.get(controlName) as AbstractControl;
  }

}
