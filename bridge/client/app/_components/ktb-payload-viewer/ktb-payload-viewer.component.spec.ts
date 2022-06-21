import { ComponentFixture, fakeAsync, flush, TestBed } from '@angular/core/testing';
import { KtbPayloadViewerComponent } from './ktb-payload-viewer.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { of, throwError } from 'rxjs';
import { Trace } from '../../_models/trace';
import { EvaluationTracesMock } from '../../_services/_mockData/trace.mock';
import { TestUtils } from '../../_utils/test.utils';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbPayloadViewerModule } from './ktb-payload-viewer.module';

describe('KtbPayloadViewerComponent', () => {
  let component: KtbPayloadViewerComponent;
  let fixture: ComponentFixture<KtbPayloadViewerComponent>;
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbPayloadViewerModule, HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbPayloadViewerComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    fixture.detectChanges();
  });

  afterEach(fakeAsync(() => {
    const payloadDialogCloseButton: HTMLElement | null = document.querySelector(
      '[uitestid=keptn-close-payload-dialog-button]'
    );
    payloadDialogCloseButton?.click();
    TestUtils.updateDialog(fixture);
    flush();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show event payload dialog without stage and service', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    const spy = jest.spyOn(dataService, 'getEvent').mockReturnValue(of(trace));
    component.type = trace.type;
    component.project = trace.data.project;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    // then

    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(component.event).toBe(trace);
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim().toString()).toBe(`This is the latest ${trace.type} event`);
  }));

  it('should show event payload dialog with stage', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    const spy = jest.spyOn(dataService, 'getEvent').mockReturnValue(of(trace));
    component.type = trace.type;
    component.project = trace.data.project;
    component.stage = trace.data.stage;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    // then
    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(component.event).toBe(trace);
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim().toString()).toBe(
      `This is the latest ${trace.type} event from stage ${trace.data.stage}`
    );
  }));

  it('should show event payload dialog with service', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    const spy = jest.spyOn(dataService, 'getEvent').mockReturnValue(of(trace));
    component.type = trace.type;
    component.project = trace.data.project;
    component.service = trace.data.service;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    // then
    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(component.event).toBe(trace);
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim().toString()).toBe(
      `This is the latest ${trace.type} event for service ${trace.data.service}`
    );
  }));

  it('should show event payload dialog with stage and service', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    const spy = jest.spyOn(dataService, 'getEvent').mockReturnValue(of(trace));
    component.type = trace.type;
    component.project = trace.data.project;
    component.stage = trace.data.stage;
    component.service = trace.data.service;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    // then
    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(component.event).toBe(trace);
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim().toString()).toBe(
      `This is the latest ${trace.type} event from stage ${trace.data.stage} for service ${trace.data.service}`
    );
  }));

  it('should show empty message', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    const spy = jest.spyOn(dataService, 'getEvent').mockReturnValue(of(undefined));
    component.type = trace.type;
    component.project = trace.data.project;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    // then
    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(component.event).toBe(undefined);
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim().toString()).toBe(`Could not find any ${trace.type} event`);
  }));

  it('should show error message', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    const spy = jest.spyOn(dataService, 'getEvent').mockReturnValue(
      throwError({
        headers: {
          normalizedNames: {},
          lazyUpdate: null,
          lazyInit: null,
          headers: {},
        },
        status: 401,
        statusText: 'Unauthorized',
        url: `http://localhost:3000/api/mongodb-datastore/event?pageSize=1&type=${trace.type}&project=${trace.data.project}`,
        ok: false,
        name: 'HttpErrorResponse',
        message: `Http failure response for http://localhost:3000/api/mongodb-datastore/event?pageSize=1&type=${trace.type}&project=${trace.data.project}: 401 Unauthorized`,
        error: 'Request failed with status code 401',
      })
    );
    component.type = trace.type;
    component.project = trace.data.project;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    // then
    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(component.event).toBe(undefined);
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim().toString()).toBe(
      `Http failure response for http://localhost:3000/api/mongodb-datastore/event?pageSize=1&type=${trace.type}&project=${trace.data.project}: 401 Unauthorized`
    );
  }));

  it('should close dialog', fakeAsync(() => {
    // given
    const trace: Trace = Trace.fromJSON(EvaluationTracesMock[0]);
    jest.spyOn(dataService, 'getEvent').mockReturnValue(of(trace));
    component.type = trace.type;
    component.project = trace.data.project;
    fixture.detectChanges();

    // when
    const showDialogButton = fixture.nativeElement.querySelector('[uitestid=keptn-show-payload-dialog-button]');
    showDialogButton.click();
    TestUtils.updateDialog(fixture);

    const payloadDialogCloseButton: HTMLElement | null = document.querySelector(
      '[uitestid=keptn-close-payload-dialog-button]'
    );
    payloadDialogCloseButton?.click();
    TestUtils.updateDialog(fixture);

    // then
    const payloadDialogMessage: HTMLElement | null = document.querySelector('[uitestid=keptn-payload-dialog-message]');
    expect(payloadDialogMessage).toBeFalsy();

    flush();
  }));
});
