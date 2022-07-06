import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceTasksListComponent } from './ktb-sequence-tasks-list.component';
import { Trace } from '../../../_models/trace';
import { Location } from '@angular/common';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSequenceViewModule } from '../ktb-sequence-view.module';

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceTasksListComponent;
  let fixture: ComponentFixture<KtbSequenceTasksListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbSequenceViewModule, RouterTestingModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceTasksListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should select last sequence', () => {
    //given, when
    component.focusedEventId = undefined;
    component.tasks = getTasks();

    // then
    expect(component.focusedEventId).toBe('97b4a9a3-610a-4697-b37b-cec522a6b42f');
  });

  it('should select last sequence if given event does not exist', () => {
    //given, when
    component.focusedEventId = '97b4a9a3-610a-4697-b37b-cec522a6b42y';
    component.tasks = getTasks();

    // then
    expect(component.focusedEventId).toBe('97b4a9a3-610a-4697-b37b-cec522a6b42f');
  });

  it('should select event through input', () => {
    //given, when
    component.focusedEventId = '97b4a9a3-610a-4697-b37b-cec522a6b42a';
    component.tasks = getTasks();

    // then
    expect(component.focusedEventId).toBe('97b4a9a3-610a-4697-b37b-cec522a6b42a');
  });

  it('should be a focused task if the root-sequence is selected', () => {
    // given
    const traces = getTasks();
    component.tasks = traces;

    // when, then
    expect(component.isFocusedTask(traces[1])).toBe(true);
  });

  it('should not be a focused task', () => {
    // given
    const traces = getTasks();
    component.tasks = traces;

    // when, then
    expect(component.isFocusedTask(traces[0])).toBe(false);
  });

  it('should be a focused task if a child trace is selected', () => {
    // given
    const traces = getTasks();
    component.tasks = traces;
    component.focusedEventId = 'c0275a97-fc4f-4279-bc3e-c892482c8dbd';

    // when, then
    expect(component.isFocusedTask(traces[1].traces[0])).toBe(true);
  });

  it('should focus event', () => {
    // given
    const location = TestBed.inject(Location);
    const locationSpy = jest.spyOn(location, 'go');

    // when
    component.focusEvent(getTasks()[0]);

    // then
    expect(locationSpy).toHaveBeenCalledWith(
      `/project/sockshop/sequence/51b96d8c-ff67-4a75-859b-124bcdbb04ff/event/97b4a9a3-610a-4697-b37b-cec522a6b42a`
    );
  });

  it('should return trace identifier', () => {
    // given, when, then
    expect(component.identifyEvent(0, getTasks()[0])).toBe('97b4a9a3-610a-4697-b37b-cec522a6b42a');
  });

  function getTasks(): Trace[] {
    return [
      {
        data: {
          project: 'sockshop',
          service: 'carts',
          stage: 'dev',
        },
        id: '97b4a9a3-610a-4697-b37b-cec522a6b42a',
        shkeptncontext: '51b96d8c-ff67-4a75-859b-124bcdbb04ff',
        shkeptnspecversion: '0.2.4',
        source: 'bridge',
        specversion: '1.0',
        time: '2022-05-25T12:49:44.580Z',
        type: 'sh.keptn.event.dev.delivery.triggered',
        traces: [
          Trace.fromJSON({
            data: {
              project: 'sockshop',
              service: 'carts',
              stage: 'dev',
            },
            gitcommitid: '05f241d1c1f022d65ab6af60bd2ef1272c1fbfbb',
            id: 'c0275a97-fc4f-4279-bc3e-c892482c8dbq',
            shkeptncontext: '99a20ef4-d822-4185-bbee-0d7a364c213b',
            shkeptnspecversion: '0.2.3',
            source: 'shipyard-controller',
            specversion: '1.0',
            time: '2022-05-25T12:49:52.560Z',
            type: 'sh.keptn.event.delivery.triggered',
          }),
        ],
      },
      {
        data: {
          project: 'sockshop',
          service: 'carts',
          stage: 'production',
        },
        id: '97b4a9a3-610a-4697-b37b-cec522a6b42f',
        shkeptncontext: '51b96d8c-ff67-4a75-859b-124bcdbb04ff',
        shkeptnspecversion: '0.2.4',
        source: 'bridge',
        specversion: '1.0',
        time: '2022-05-25T12:50:44.580Z',
        type: 'sh.keptn.event.production.delivery.triggered',
        traces: [
          Trace.fromJSON({
            data: {
              approval: {
                pass: 'manual',
                warning: 'manual',
              },
              project: 'sockshop',
              service: 'carts',
              stage: 'production',
            },
            gitcommitid: '05f241d1c1f022d65ab6af60bd2ef1272c1fbfbb',
            id: 'c0275a97-fc4f-4279-bc3e-c892482c8dbd',
            shkeptncontext: '99a20ef4-d822-4185-bbee-0d7a364c213b',
            shkeptnspecversion: '0.2.3',
            source: 'shipyard-controller',
            specversion: '1.0',
            time: '2022-05-25T12:50:52.560Z',
            type: 'sh.keptn.event.approval.triggered',
          }),
        ],
      },
    ].map((trace) => Trace.fromJSON(trace));
  }
});
