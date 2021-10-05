import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';
import { KtbPayloadViewerComponent } from './ktb-payload-viewer.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataServiceMock } from '../../_services/data.service.mock';
import { DataService } from '../../_services/data.service';
import { of } from 'rxjs';
import { Trace } from '../../_models/trace';
import { EvaluationTracesMock } from '../../_models/trace.mock';
import { TestUtils } from '../../_utils/test.utils';

describe('KtbPayloadViewerComponent', () => {
  let component: KtbPayloadViewerComponent;
  let fixture: ComponentFixture<KtbPayloadViewerComponent>;
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [DataServiceMock],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbPayloadViewerComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    fixture.detectChanges();
  });

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
    const payloadDialogCloseButton = document.querySelector('[uitestid=keptn-close-payload-dialog-button]');
    const payloadDialogMessage = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(payloadDialogCloseButton).toBeTruthy();
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim()).toBe(`This is the latest ${trace.type} event`);
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
    const payloadDialogCloseButton = document.querySelector('[uitestid=keptn-close-payload-dialog-button]');
    const payloadDialogMessage = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(payloadDialogCloseButton).toBeTruthy();
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim()).toBe(
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
    const payloadDialogCloseButton = document.querySelector('[uitestid=keptn-close-payload-dialog-button]');
    const payloadDialogMessage = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(payloadDialogCloseButton).toBeTruthy();
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim()).toBe(
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
    const payloadDialogCloseButton = document.querySelector('[uitestid=keptn-close-payload-dialog-button]');
    const payloadDialogMessage = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(payloadDialogCloseButton).toBeTruthy();
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim()).toBe(
      `This is the latest ${trace.type} event from stage ${trace.data.stage} for service ${trace.data.service}`
    );
  }));

  it('should show error message', fakeAsync(() => {
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
    const payloadDialogCloseButton = document.querySelector('[uitestid=keptn-close-payload-dialog-button]');
    const payloadDialogMessage = document.querySelector('[uitestid=keptn-payload-dialog-message]');

    expect(spy).toHaveBeenCalled();
    expect(payloadDialogCloseButton).toBeTruthy();
    expect(payloadDialogMessage).toBeTruthy();
    expect(payloadDialogMessage?.textContent?.trim()).toBe(`Could not load any ${trace.type} event`);
  }));
});
