import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AbstractControl } from '@angular/forms';
import { WebhookConfigMock } from '../../_services/_mockData/webhook-config.mock';
import { Secret } from '../../_models/secret';
import { SecretScopeDefault } from '../../../../shared/interfaces/secret-scope';
import { APIService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';

describe('KtbWebhookSettingsComponent', () => {
  let component: KtbWebhookSettingsComponent;
  let fixture: ComponentFixture<KtbWebhookSettingsComponent>;

  const secretDataSource = [
    {
      name: 'SecretA',
      keys: [
        {
          name: 'key1',
          path: '.secret.SecretA.key1',
        },
        {
          name: 'key2',
          path: '.secret.SecretA.key2',
        },
        {
          name: 'key3',
          path: '.secret.SecretA.key3',
        },
      ],
    },
    {
      name: 'SecretB',
      keys: [
        {
          name: 'key1',
          path: '.secret.SecretB.key1',
        },
        {
          name: 'key2',
          path: '.secret.SecretB.key2',
        },
        {
          name: 'key3',
          path: '.secret.SecretB.key3',
        },
      ],
    },
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: APIService, useClass: ApiServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbWebhookSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should be invalid URL when not starting with http(s)://', () => {
    const urlControl: AbstractControl = component.webhookConfigForm.get('url') as AbstractControl;
    const urls = [
      '://keptn.sh',
      'www.keptn.sh',
      'http:/www.keptn.sh',
      'http//www.keptn.sh',
      'htp://www.keptn.sh',
      'ftp://www.keptn.sh',
    ];

    for (const url of urls) {
      urlControl.setValue(url);
      expect(urlControl.valid).toEqual(false);
    }
  });

  it('should be invalid URL if it contains spaces', () => {
    const urlControl: AbstractControl = component.webhookConfigForm.get('url') as AbstractControl;
    const urls = [
      'http://my-jenkins-servcer.default.svc.cluster.local:8080/job/nodejs example app/build?token=',
      'http://my-jenkins-servcer.com/job/nodejs example app/build?token=',
    ];

    for (const url of urls) {
      urlControl.setValue(url);
      expect(urlControl.valid).toEqual(false);
    }
  });

  it('should be valid URL when it starts with http(s)://', () => {
    const urlControl: AbstractControl = component.webhookConfigForm.get('url') as AbstractControl;
    const urls = [
      'https://keptn.sh',
      'http://keptn.sh',
      'http://www.keptn.sh',
      'http://my-jenkins-servcer.default.svc.cluster.local:8080/job/nodejs%20example%20app/build?token=',
    ];

    for (const url of urls) {
      urlControl.setValue(url);
      expect(urlControl.valid).toEqual(true);
    }
  });

  it('should add headers', () => {
    // given
    const addHeaderButton = getAddHeaderButton();
    let headers = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    expect(headers.length).toEqual(0);

    // then
    for (let i = 1; i <= 2; ++i) {
      addHeaderButton.click();
      fixture.detectChanges();
      headers = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
      expect(headers.length).toEqual(i);
    }
  });

  it('should remove headers', () => {
    // given
    setParameters();
    fixture.detectChanges();

    // when
    let headerRows = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    const lengthBefore = headerRows.length;
    const buttons = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form>div>button');

    buttons[0].click();
    fixture.detectChanges();

    // then
    headerRows = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    expect(lengthBefore - 1).toEqual(headerRows.length);
  });

  it('should fill form fields with provided data', () => {
    // given
    setParameters();
    fixture.detectChanges();

    // then
    const headers = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    expect(component.getFormControl('url').value).toEqual(WebhookConfigMock.url);
    expect(component.getFormControl('method').value).toEqual(WebhookConfigMock.method);
    expect(component.getFormControl('payload').value).toEqual(WebhookConfigMock.payload);
    expect(component.getFormControl('proxy').value).toEqual(WebhookConfigMock.proxy);

    for (let i = 0; i < component.headerControls.length; ++i) {
      expect(WebhookConfigMock.header[i]).toEqual({
        name: component.headerControls[i].get('name')?.value,
        value: component.headerControls[i].get('value')?.value,
      });
    }
    expect(headers.length).toEqual(WebhookConfigMock.header.length);
  });

  function setParameters(): void {
    component.webhook = WebhookConfigMock;
  }

  it('should be invalid form if only URL is set', () => {
    // given
    const urlControl = component.getFormControl('url');
    urlControl.setValue('keptn.sh');
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be invalid form if proxy is invalid', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');
    const proxyControl = component.getFormControl('proxy');

    // when
    urlControl.setValue('keptn.sh');
    methodControl.setValue('POST');
    proxyControl.setValue('keptn');
    // then
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be invalid form if it has empty header configuration', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');

    // when
    urlControl.setValue('keptn.sh');
    methodControl.setValue('POST');

    const addHeaderButton = getAddHeaderButton();
    addHeaderButton.click();
    component.addHeader();

    // then
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be valid form if URL and method is set', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  it('should be valid form if valid proxy is set', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');
    const proxyControl = component.getFormControl('proxy');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');
    proxyControl.setValue('https://keptn.sh');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  it('should be valid form if valid header is added', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');

    // when
    component.addHeader('content-type', 'application/json');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  it('should map secrets to a tree when set', () => {
    // given, when
    const secrets = [new Secret(), new Secret()];
    secrets[0].name = 'SecretA';
    secrets[0].scope = SecretScopeDefault.WEBHOOK;
    secrets[0].keys = ['key1', 'key2', 'key3'];
    secrets[1].name = 'SecretB';
    secrets[1].scope = SecretScopeDefault.WEBHOOK;
    secrets[1].keys = ['key1', 'key2', 'key3'];

    // when
    component.secrets = secrets;

    // then
    expect(component.secretDataSource).toEqual(secretDataSource);
  });

  it('should emit the set values in the form', () => {
    // given
    const spy = jest.spyOn(component.webhookChange, 'emit');
    component.getFormControl('method').setValue(component.webhookMethods[0]);
    component.getFormControl('url').setValue('https://example.com');
    component.addHeader();
    component.headerControls[0].get('name')?.setValue('x-token');
    component.headerControls[0].get('value')?.setValue('token-value');
    component.getFormControl('payload').setValue('payload');
    component.getFormControl('proxy').setValue('https://proxy.com');

    // when
    component.onWebhookFormChange();

    // then
    expect(spy).toHaveBeenCalledTimes(1);
    expect(spy).toHaveBeenCalledWith({
      header: [{ name: 'x-token', value: 'token-value' }],
      method: 'GET',
      payload: 'payload',
      proxy: 'https://proxy.com',
      sendFinished: true,
      sendStarted: true,
      type: '',
      url: 'https://example.com',
    });
  });

  it('should set sendFinished and sendStarted true by default for triggered events', () => {
    // given
    component.eventType = 'triggered';
    fixture.detectChanges();

    // then
    expect(component.getFormControl('sendFinished').value).toEqual('true');
    expect(component.getFormControl('sendStarted').value).toEqual('true');
  });

  it('should set sendFinished and sendStarted null when eventType started', () => {
    // given
    component.eventType = 'started';
    fixture.detectChanges();

    // then
    expect(component.getFormControl('sendFinished').value).toEqual(null);
    expect(component.getFormControl('sendStarted').value).toEqual(null);
  });

  it('should set sendFinished and sendStarted null for finished events', () => {
    // given
    component.eventType = 'finished';
    fixture.detectChanges();

    // then
    expect(component.getFormControl('sendFinished').value).toEqual(null);
    expect(component.getFormControl('sendStarted').value).toEqual(null);
  });

  it('should have sendFinished and sendStarted set to true', () => {
    // given
    component.eventType = 'triggered';
    component.webhook = {
      header: [{ name: 'x-token', value: 'token-value' }],
      method: 'GET',
      payload: 'payload',
      proxy: 'https://proxy.com',
      sendFinished: true,
      sendStarted: true,
      filter: {
        projects: null,
        services: null,
        stages: null,
      },
      type: '',
      url: 'https://example.com',
    };

    // when

    // then
    expect(component.getFormControl('sendFinished').value).toEqual('true');
    expect(component.getFormControl('sendStarted').value).toEqual('true');
  });

  it('should have sendFinished and sendStarted set to false', () => {
    // given
    component.eventType = 'triggered';
    component.webhook = {
      header: [{ name: 'x-token', value: 'token-value' }],
      method: 'GET',
      payload: 'payload',
      proxy: 'https://proxy.com',
      sendFinished: false,
      sendStarted: false,
      filter: {
        projects: null,
        services: null,
        stages: null,
      },
      type: '',
      url: 'https://example.com',
    };

    // when

    // then
    expect(component.getFormControl('sendFinished').value).toEqual('false');
    expect(component.getFormControl('sendStarted').value).toEqual('false');
  });

  it('should correctly set event payload', () => {
    component.eventPayload = {
      data: {
        myCustomData: [
          [
            {
              myCustomKey: {
                myCustomElement: undefined,
              },
            },
            {
              myCustomKey2: {
                myCustomElement2: undefined,
              },
            },
          ],
        ],
        project: undefined,
      },
      id: undefined,
      keptnContext: undefined,
    };
    expect(component.eventDataSource).toEqual([
      {
        keys: [
          {
            keys: [
              {
                keys: [
                  {
                    keys: [
                      {
                        keys: [
                          {
                            name: 'myCustomElement',
                            path: '(index (index .data.myCustomData 0) 0).myCustomKey.myCustomElement',
                          },
                        ],
                        name: 'myCustomKey',
                      },
                    ],
                    name: '[0]',
                  },
                  {
                    keys: [
                      {
                        keys: [
                          {
                            name: 'myCustomElement2',
                            path: '(index (index .data.myCustomData 0) 1).myCustomKey2.myCustomElement2',
                          },
                        ],
                        name: 'myCustomKey2',
                      },
                    ],
                    name: '[1]',
                  },
                ],
                name: '[0]',
              },
            ],
            name: 'myCustomData',
          },
          {
            name: 'project',
            path: '.data.project',
          },
        ],
        name: 'data',
      },
      {
        name: 'id',
        path: '.id',
      },
      {
        name: 'keptnContext',
        path: '.keptnContext',
      },
    ]);
  });

  it('should have an error when payload contains one of these characters: $ | ; > & ` /var/run', () => {
    // given
    const payloadControl = component.getFormControl('payload');

    const chars = ['$', '$(', '|', ';', '>', '&', '&&', '`', '/var/run', '/VAR/RUN', '/vAr/RuN'];

    for (let i = 0; i < chars.length; i++) {
      const val = `{id: 12345${chars[i]}678}`;
      // when
      payloadControl.setValue(val);
      payloadControl.updateValueAndValidity();

      // then
      expect(payloadControl.valid).toEqual(false);
    }
  });

  it('should be an invalid form when payload contains a special character', () => {
    // given
    const payloadControl = component.getFormControl('payload');

    // when
    payloadControl.setValue('{id: 12345$$678}');
    payloadControl.updateValueAndValidity();

    // then
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  function getAddHeaderButton(): HTMLElement {
    return fixture.nativeElement.querySelector('[uitestid="ktb-webhook-settings-add-header-button"]');
  }
});
