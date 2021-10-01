import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { AbstractControl, FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { WebhookConfigMethod } from '../../../../shared/interfaces/webhook-config';
import { WebhookConfig } from '../../../../shared/models/webhook-config';
import { Secret } from '../../_models/secret';
import { SelectTreeNode, TreeListSelectOptions } from '../ktb-tree-list-select/ktb-tree-list-select.component';

type ControlType = 'method' | 'url' | 'payload' | 'proxy' | 'header';


@Component({
  selector: 'ktb-webhook-settings',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss'],
})
export class KtbWebhookSettingsComponent implements OnInit {
  private _webhook: WebhookConfig = new WebhookConfig();
  public webhookConfigForm = new FormGroup({
    method: new FormControl('', [Validators.required]),
    url: new FormControl('', [Validators.required, FormUtils.isUrlValidatorWithVariable, FormUtils.urlSpecialCharsWithVariablesValidator]),
    payload: new FormControl('', []),
    header: new FormArray([]),
    proxy: new FormControl('', [FormUtils.isUrlValidator, FormUtils.urlSpecialCharsValidator]),
  });
  public webhookMethods: WebhookConfigMethod[] = ['GET', 'POST', 'PUT'];
  public secretDataSource: SelectTreeNode[] = [];
  public secretOptions: TreeListSelectOptions = {headerText: 'selectSecret', emptyText: 'No secrets can be found.<p>Secrets can be configured under the menu entry "Secrets" in the Uniform.</p>'};

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

  @Input()
  set secrets(secrets: Secret[] | undefined) {
    if (secrets) {
      this.secretDataSource = secrets.map((secret: Secret) => this.mapSecret(secret));
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
    this.webhookChange.emit(this._webhook);
  }

  public addHeader(name?: string, value?: string): void {
    this.header.push(new FormGroup({
      name: new FormControl(name || '', [Validators.required]),
      value: new FormControl(value || '', [Validators.required]),
    }));
  }

  public removeHeader(index: number): void {
    this.header.removeAt(index);
  }

  public getFormControl(controlName: ControlType): AbstractControl {
    return this.webhookConfigForm.get(controlName) as AbstractControl;
  }

  public setSecret(secret: string, controlName: ControlType, selectionStart: number, controlIndex?: number): void {
    let control: AbstractControl;
    if (controlName === 'header' && controlIndex !== undefined) {
      const group = this.header.at(controlIndex) as FormGroup;
      control = group.controls.value;
    } else {
      control = this.getFormControl(controlName);
    }

    const secretVar = `{{.${secret}}}`;
    const firstPart = control.value.slice(0, selectionStart);
    const secondPart = control.value.slice(selectionStart, control.value.length);
    const finalString = firstPart + secretVar + secondPart;

    control.setValue(finalString);
    // Input event detection is not working reliable for adding secrets, so we have to call it to work properly
    this.onWebhookFormChange();
  }

  private mapSecret(secret: Secret): SelectTreeNode {
    const scrt: SelectTreeNode = {name: secret.name};
    if (secret.keys) {
      scrt.keys = secret.keys.map((key: string) => {
        return {name: key, path: `${secret.name}.${key}`};
      });
      scrt.keys.sort((a, b) => a.name.localeCompare(b.name));
    }
    return scrt;
  }
}
