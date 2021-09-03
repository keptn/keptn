import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DataService } from '../../_services/data.service';
import { FormUtils } from '../../_utils/form.utils';
import { Project } from '../../_models/project';
import { UniformSubscription } from '../../_models/uniform-subscription';

@Component({
  selector: 'ktb-webhook-settings',
  templateUrl: './ktb-webhook-settings.component.html',
  styleUrls: ['./ktb-webhook-settings.component.scss']
})
export class KtbWebhookSettingsComponent implements OnInit {

  public _project?: Project;
  public _subscription?: UniformSubscription;
  private _prevFilter?: { projects: string[] | null; stages: string[] | null; services: string[] | null };
  public webhookConfigForm = this.formBuilder.group({
    method: ['', [Validators.required]],
    url: ['', [Validators.required, Validators.pattern(FormUtils.URL_PATTERN)]],
    payload: ['', [Validators.required]],
    header: this.formBuilder.array([]),
    proxy: ['', [Validators.pattern(FormUtils.URL_PATTERN)]],
  });

  public webhookMethods = ['POST', 'PUT'];

  public loading = false;

  @Input()
  get project(): Project | undefined {
    return this._project;
  }
  set project(value: Project | undefined) {
    if (this._project !== value) {
      this._project = value;
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
    }
  }

  get prevFilter(): { projects: string[] | null; stages: string[] | null; services: string[] | null } | undefined {
    return this._prevFilter;
  }
  set prevFilter(value: { projects: string[] | null; stages: string[] | null; services: string[] | null } | undefined) {
    if (this._prevFilter !== value) {
      this._prevFilter = JSON.parse(JSON.stringify(value));
    }
  }

  get projectName(): string {
    return this.project?.projectName || '';
  }

  get header(): FormArray | null {
    return this.webhookConfigForm.get('header') as FormArray;
  }

  get headerControls(): FormGroup[] {
    return this.header?.controls as FormGroup[] || [];
  }

  constructor(private dataService: DataService, private formBuilder: FormBuilder) {
  }

  ngOnInit(): void {
    this.loading = true;
    const stage: string | undefined = this.subscription?.filter?.stages?.length ? this.subscription?.filter?.stages[0] : undefined;
    const services: string | undefined = this.subscription?.filter?.services?.length ? this.subscription?.filter?.services[0] : undefined;
    this.dataService.getWebhookConfig(this.projectName, stage, services)
      .subscribe(webhookConfig => {
        this.webhookConfigForm?.get('method')?.setValue(webhookConfig.method);
        this.webhookConfigForm?.get('url')?.setValue(webhookConfig.url);
        this.webhookConfigForm?.get('payload')?.setValue(webhookConfig.payload);
        this.webhookConfigForm?.get('proxy')?.setValue(webhookConfig.proxy);

        for (const header of webhookConfig.header || []) {
          this.addHeader(header.name, header.value);
        }

        this.loading = false;
      }, err => {
        this.loading = false;
      });
  }

  public addHeader(name?: string, value?: string): void {
    this.header?.push(this.formBuilder.group({
      name: [name, [Validators.required]],
      value: [value, [Validators.required]]
    }));
  }

  public removeHeader(index: number): void {
    this.header?.removeAt(index);
  }

  public getFormControl(controlName: string): AbstractControl | null {
    return this.webhookConfigForm.get(controlName);
  }

}
