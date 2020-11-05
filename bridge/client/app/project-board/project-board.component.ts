import {ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {filter, map, startWith, switchMap, takeUntil} from "rxjs/operators";
import {Observable, Subject, Subscription, timer} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";
import {Location} from "@angular/common";

import * as moment from 'moment';

import {Root} from "../_models/root";
import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";
import {ApiService} from "../_services/api.service";
import DateUtil from "../_utils/date.utils";
import {Service} from "../_models/service";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {DtCheckboxChange} from "@dynatrace/barista-components/checkbox";
import {EVENT_LABELS} from "../_models/event-labels";
import {DtOverlayConfig} from "@dynatrace/barista-components/overlay";
import {DtToggleButtonItem} from "@dynatrace/barista-components/toggle-button-group";
import {ClipboardService} from "../_services/clipboard.service";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  private _tracesTimer: Subscription = Subscription.EMPTY;

  public project$: Observable<Project>;
  public openApprovals$: Observable<Trace[]>;

  public currentRoot: Root;
  public error: string = null;

  private _rootEventsTimerInterval = 30;
  private _tracesTimerInterval = 10;

  public projectName: string;
  public serviceName: string;
  public contextId: string;
  public eventId: string;

  public view: string = 'services';
  public selectedStage: Stage = null;

  public eventTypes: string[] = [];
  public filterEventTypes: string[] = [];

  public filterEventType: string = null;

  public integrationsExternalDetails = null;

  public useCaseExamples = {
    'cli': [],
    'api': []
  };

  public keptnInfo: any;
  public currentTime: String;

  @ViewChild('problemFilterEventButton') public problemFilterEventButton: DtToggleButtonItem<string>;
  @ViewChild('evaluationFilterEventButton') public evaluationFilterEventButton: DtToggleButtonItem<string>;
  @ViewChild('approvalFilterEventButton') public approvalFilterEventButton: DtToggleButtonItem<string>;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private location: Location, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService, private clipboard: ClipboardService) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        if(params["shkeptncontext"]) {
          this.contextId = params["shkeptncontext"];
          this.apiService.getTraces(this.contextId)
            .pipe(
              map(response => response.body),
              map(result => result.events||[]),
              map(traces => traces.map(trace => Trace.fromJSON(trace)))
            )
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe((traces: Trace[]) => {
              if(traces.length > 0) {
                if(params["eventselector"]) {
                  let trace = traces.find((t: Trace) => t.data.stage == params["eventselector"] && !!t.getProject() && !!t.getService());
                  if(!trace)
                    trace = traces.reverse().find((t: Trace) => t.type == params["eventselector"] && !!t.getProject() && !!t.getService());

                  if(trace)
                    this.router.navigate(['/project', trace.getProject(), trace.getService(), trace.shkeptncontext, trace.id]);
                  else
                    this.error = "trace";
                } else {
                  let trace = traces.find((t: Trace) => !!t.getProject() && !!t.getService());
                  this.router.navigate(['/project', trace.getProject(), trace.getService(), trace.shkeptncontext]);
                }
              } else {
                this.error = "trace";
              }
            });
        } else {
          this.projectName = params["projectName"];
          this.serviceName = params["serviceName"];
          this.contextId = params["contextId"];
          this.eventId = params["eventId"];
          this.currentRoot = null;

          this.project$ = this.dataService.projects.pipe(
            map(projects => projects ? projects.find(project => {
              return project.projectName === params['projectName'];
            }) : null)
          );
          this.openApprovals$ = this.dataService.openApprovals;

          this.project$
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(project => {
              if(project === undefined)
                this.error = 'project';
              this._changeDetectorRef.markForCheck();
            }, error => {
              this.error = 'projects';
            });

          this.dataService.roots
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(roots => {
              if(roots) {
                if(!this.currentRoot)
                  this.currentRoot = roots.find(r => r.shkeptncontext == params["contextId"]);
                this.eventTypes = this.eventTypes.concat(roots.map(r => r.type)).filter((r, i, a) => a.indexOf(r) === i);
              }
              if(this.currentRoot && !this.eventId)
                this.eventId = this.currentRoot.traces[this.currentRoot.traces.length-1].id;
            });

          timer(0, this._rootEventsTimerInterval*1000)
            .pipe(
              startWith(0),
              switchMap(() => this.project$),
              filter(project => !!project && !!project.getServices())
            )
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(project => {
              this.updateIntegrations();
              project.getServices().forEach(service => {
                this.dataService.loadRoots(project, service);
              });
            });
        }
      });

    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.keptnInfo = keptnInfo;
        if(this.keptnInfo.bridgeInfo.keptnInstallationType.indexOf("CONTINUOUS_DELIVERY") != -1) {
          this.addDeploymentUseCaseToIntegrations();
        }
        if(this.keptnInfo.bridgeInfo.keptnInstallationType.indexOf("QUALITY_GATES_ONLY") != -1) {
          this.addEvaluationUseCaseToIntegrations();
        }
        if(this.keptnInfo.bridgeInfo.keptnInstallationType.indexOf("CONTINUOUS_OPERATIONS") != -1) {
          this.addRemediationUseCaseToIntegrations();
        }
      });
  }

  updateIntegrations() {
    if(this.keptnInfo.bridgeInfo.keptnInstallationType.indexOf("QUALITY_GATES_ONLY") != -1) {
      this.currentTime = moment().startOf('minute').format("YYYY-MM-DDTHH:mm:ss");
      this.useCaseExamples['cli'].find(e => e.label == 'Trigger a quality gate evaluation').code = `keptn send event start-evaluation --project=\${PROJECT} --stage=\${STAGE} --service=\${SERVICE} --start=${this.currentTime} --timeframe=5m`;
      this.useCaseExamples['api'].find(e => e.label == 'Trigger a quality gate evaluation').code = `curl -X POST "\${KEPTN_API_ENDPOINT}/v1/project/\${PROJECT}/stage/\${STAGE}/service/\${SERVICE}/evaluation" \\
    -H "accept: application/json; charset=utf-8" \\
    -H "x-token: \${KEPTN_API_TOKEN}" \\
    -H "Content-Type: application/json; charset=utf-8" \\
    -d "{"start": "${this.currentTime}", "timeframe": "5m", "labels":{"buildId":"build-17","owner":"JohnDoe","testNo":"47-11"}"`;
    }
  }

  addEvaluationUseCaseToIntegrations() {
    this.useCaseExamples['cli'].push({
      label: 'Trigger a quality gate evaluation',
      code: ''
    });
    this.useCaseExamples['api'].push({
      label: 'Trigger a quality gate evaluation',
      code: ''
    });
  }

  addDeploymentUseCaseToIntegrations() {
    this.useCaseExamples['cli'].push({
      label: 'Trigger deployment with a new artifact',
      code: `keptn send event new-artifact --project=\${PROJECT} --service=\${SERVICE}--image=\${IMAGE} --tag=\${TAG}`
    });
    this.useCaseExamples['api'].push({
      label: 'Trigger deployment with a new artifact',
      code: `curl -X POST "\${KEPTN_API_ENDPOINT}/v1/event" \\
      -H "accept: application/json; charset=utf-8" -H "x-token: \${KEPTN_API_TOKEN}" -H "Content-Type: application/json; charset=utf-8" \\
      -d "{"type":"sh.keptn.event.configuration.change","specversion":"0.2","source":"api","contenttype":"application\\/json","data":{"project":"\${PROJECT}","stage":"\${STAGE}","service":"\${SERVICE}","valuesCanary":{"image":"\${IMAGE}"}}}"`
    });
  }

  addRemediationUseCaseToIntegrations() {
    this.useCaseExamples['cli'].push({
      label: 'Trigger remediation with a dummy problem event (Note: Linux/mac OS only)',
      code: `echo '{"type":"sh.keptn.event.problem.open","specversion":"0.2","source":"api","contenttype":"application\\/json","data":{"State":"OPEN","ProblemID":"\${PROBLEM_ID}","ProblemTitle":"\${PROBLEM}","project":"\${PROJECT}","stage":"\${STAGE}","service":"\${SERVICE}"}}' > dummy_problem.json \\
      keptn send event -f=dummy_problem.json`
    });
    this.useCaseExamples['api'].push({
      label: 'Trigger remediation with a dummy problem event',
      code: `curl -X POST "\${KEPTN_API_ENDPOINT}/v1/event" \\
      -H "accept: application/json; charset=utf-8" -H "x-token: \${KEPTN_API_TOKEN}" -H "Content-Type: application/json; charset=utf-8" \\
      -d "{"type":"sh.keptn.event.problem.open","specversion":"0.2","source":"api","contenttype":"application\\/json","data":{"State":"OPEN","ProblemID":"\${PROBLEM_ID}","ProblemTitle":"\${PROBLEM}","project":"\${PROJECT}","stage":"\${STAGE}","service":"\${SERVICE}"}}"`
    });
  }

  selectRoot(event: any): void {
    this.projectName = event.root.getProject();
    this.serviceName = event.root.getService();
    this.contextId = event.root.data.shkeptncontext;
    this.eventId = null;
    if(event.stage) {
      let focusEvent = event.root.traces.find(trace => trace.data.stage == event.stage);
      let routeUrl = this.router.createUrlTree(['/project', focusEvent.getProject(), focusEvent.getService(), focusEvent.shkeptncontext, focusEvent.id]);
      this.eventId = focusEvent.id;
      this.location.go(routeUrl.toString());
    } else {
      let routeUrl = this.router.createUrlTree(['/project', event.root.getProject(), event.root.getService(), event.root.shkeptncontext]);
      this.eventId = event.root.traces[event.root.traces.length-1].id;
      this.location.go(routeUrl.toString());
    }

    this.currentRoot = event.root;
    this.loadTraces(this.currentRoot);
  }

  selectDeployment(deployment: Trace, project: Project) {
    this.selectRoot({
      root: project.getServices().find(service => service.serviceName === deployment.data.service).roots.find(root => root.shkeptncontext === deployment.shkeptncontext),
      stage: deployment.data.stage
    });
  }

  loadTraces(root: Root): void {
    this._tracesTimer.unsubscribe();
    if(moment().subtract(1, 'day').isBefore(root.time)) {
      this._tracesTimer = timer(0, this._tracesTimerInterval*1000)
        .subscribe(() => {
          this.dataService.loadTraces(root);
        });
    } else {
      this.dataService.loadTraces(root);
      this._tracesTimer = Subscription.EMPTY;
    }
  }

  getCalendarFormats() {
    return DateUtil.getCalendarFormats(true);
  }

  getRootsLastUpdated(project: Project, service: Service): Date {
    return this.dataService.getRootsLastUpdated(project, service);
  }

  getTracesLastUpdated(root: Root): Date {
    return this.dataService.getTracesLastUpdated(root);
  }

  showReloadButton(root: Root) {
    return moment().subtract(1, 'day').isAfter(root.time);
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  trackStage(index: number, stage: Stage) {
    return stage.stageName;
  }

  selectView(view) {
    this.view = view;
  }

  filterEvents(event: DtCheckboxChange<string>, eventType: string): void {
    let index = this.filterEventTypes.indexOf(eventType);
    if(index == -1) {
      this.filterEventTypes.push(eventType);
    } else {
      this.filterEventTypes.splice(index, 1);
    }
  }

  isFilteredEvent(eventType: string) {
    return this.filterEventTypes.indexOf(eventType) == -1;
  }

  getEventLabel(key): string {
    return EVENT_LABELS[key] || key;
  }

  getFilteredRoots(roots: Root[]) {
    if(roots)
      return roots.filter(r => this.filterEventTypes.indexOf(r.type) == -1);
  }

  selectStage($event, stage: Stage, filterType?: string) {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();

    this.selectedStage = stage;
    this.filterEventType = filterType;
    $event.stopPropagation();
  }

  selectFilterEvent($event) {
    if($event.isUserInput)
      this.filterEventType = $event.source.selected ? $event.value : null;
  }

  countOpenApprovals(openApprovals: Trace[], project: Project, stage: Stage, service?: Service) {
    return this.getOpenApprovals(openApprovals, project, stage, service).length;
  }

  getOpenApprovals(openApprovals: Trace[], project: Project, stage: Stage, service?: Service) {
    return openApprovals.filter(approval => approval.data.project == project.projectName && approval.data.stage == stage.stageName && (!service || approval.data.service == service.serviceName));
  }

  findFailedRootEvent(failedRootEvents: Root[], service: Service) {
    return failedRootEvents.find(root => root.data.service == service.serviceName);
  }

  findProblemEvent(problemEvents: Root[], service: Service) {
    return problemEvents.find(root => root?.data.service == service.serviceName);
  }

  loadIntegrations() {
    this.integrationsExternalDetails = '<p>Loading ...</p>';
    this.apiService.getIntegrationsPage()
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((result: string) => {
        this.integrationsExternalDetails = result;
      }, (err: Error) => {
        this.integrationsExternalDetails = '<p>Couldn\'t load page. For more details see <a href="https://keptn.sh/docs/integrations/" target="_blank" rel="noopener noreferrer">https://keptn.sh/docs/integrations/</a>';
      });
  }

  copyApiToken() {
    this.clipboard.copy(this.keptnInfo.bridgeInfo.apiToken, 'API token');
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this._tracesTimer.unsubscribe();
  }

}
