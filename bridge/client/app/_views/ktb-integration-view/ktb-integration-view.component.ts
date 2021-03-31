import {Component, OnDestroy, OnInit} from '@angular/core';
import {filter, takeUntil} from 'rxjs/operators';
import {DataService} from '../../_services/data.service';
import {Subject, timer} from 'rxjs';
import * as moment from 'moment';
import {ClipboardService} from '../../_services/clipboard.service';
import {ApiService} from "../../_services/api.service";

@Component({
  selector: 'ktb-integration-view',
  templateUrl: './ktb-integration-view.component.html',
  styleUrls: ['./ktb-integration-view.component.scss']
})
export class KtbIntegrationViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private _rootEventsTimerInterval = 30;

  public currentTime: string;
  public keptnInfo: any;
  public integrationsExternalDetails = null;
  public useCaseExamples = {
    cli: [],
    api: []
  };

  constructor(private dataService: DataService, private clipboard: ClipboardService, private apiService: ApiService) { }

  ngOnInit(): void {
    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.keptnInfo = keptnInfo;
        if (this.keptnInfo.bridgeInfo.keptnInstallationType) {
          if (this.keptnInfo.bridgeInfo.keptnInstallationType.includes('CONTINUOUS_DELIVERY')) {
            this.addDeploymentUseCaseToIntegrations();
          }
          if (this.keptnInfo.bridgeInfo.keptnInstallationType.includes('QUALITY_GATES')) {
            this.addEvaluationUseCaseToIntegrations();
          }
          if (this.keptnInfo.bridgeInfo.keptnInstallationType.includes('CONTINUOUS_OPERATIONS')) {
            this.addRemediationUseCaseToIntegrations();
          }
          this.updateIntegrations();
        }
      });
    timer(this._rootEventsTimerInterval*1000, this._rootEventsTimerInterval*1000)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.updateIntegrations();
      });
  }

  updateIntegrations() {
    if (this.keptnInfo && this.keptnInfo.bridgeInfo.keptnInstallationType && this.keptnInfo.bridgeInfo.keptnInstallationType.includes('QUALITY_GATES')) {
      this.currentTime = moment.utc().startOf('minute').format('YYYY-MM-DDTHH:mm:ss');
      this.useCaseExamples.cli.find(e => e.label === 'Trigger a quality gate evaluation').code = `keptn trigger evaluation --project=\${PROJECT} --stage=\${STAGE} --service=\${SERVICE} --start=\${this.currentTime} --timeframe=5m`;
      this.useCaseExamples.api.find(e => e.label === 'Trigger a quality gate evaluation').code = `curl -X POST "\${KEPTN_API_ENDPOINT}/v1/project/\${PROJECT}/stage/\${STAGE}/service/\${SERVICE}/evaluation" \\
    -H "accept: application/json; charset=utf-8" \\
    -H "x-token: \${KEPTN_API_TOKEN}" \\
    -H "Content-Type: application/json; charset=utf-8" \\
    -d "{"start": "${this.currentTime}", "timeframe": "5m", "labels":{"buildId":"build-17","owner":"JohnDoe","testNo":"47-11"}"`;
    }
  }

  addEvaluationUseCaseToIntegrations() {
    this.useCaseExamples.cli.push({
      label: 'Trigger a quality gate evaluation',
      code: ''
    });
    this.useCaseExamples.api.push({
      label: 'Trigger a quality gate evaluation',
      code: ''
    });
  }

  addDeploymentUseCaseToIntegrations() {
    this.useCaseExamples.cli.push({
      label: 'Trigger deployment with a new artifact',
      code: `keptn trigger delivery --project=\${PROJECT} --service=\${SERVICE} --image=\${IMAGE} --tag=\${TAG} --sequence=\${SEQUENCE}`
    });
    this.useCaseExamples.api.push({
      label: 'Trigger deployment with a new artifact',
      code: `curl -X POST "\${KEPTN_API_ENDPOINT}/v1/event" \\
      -H "accept: application/json; charset=utf-8" -H "x-token: \${KEPTN_API_TOKEN}" -H "Content-Type: application/json; charset=utf-8" \\
      -d "{"type":"sh.keptn.event.configuration.change","specversion":"0.2","source":"api","contenttype":"application\\/json","data":{"project":"\${PROJECT}","stage":"\${STAGE}","service":"\${SERVICE}","configurationChange":{"values":{"image":"\${IMAGE}"}}}}"`
    });
  }

  addRemediationUseCaseToIntegrations() {
    this.useCaseExamples.cli.push({
      label: 'Trigger remediation with a dummy problem event (Note: Linux/mac OS only)',
      code: `echo '{"type":"sh.keptn.event.problem.open","specversion":"0.2","source":"api","contenttype":"application\\/json","data":{"State":"OPEN","ProblemID":"\${PROBLEM_ID}","ProblemTitle":"\${PROBLEM}","project":"\${PROJECT}","stage":"\${STAGE}","service":"\${SERVICE}"}}' > dummy_problem.json \\
      keptn send event -f=dummy_problem.json`
    });
    this.useCaseExamples.api.push({
      label: 'Trigger remediation with a dummy problem event',
      code: `curl -X POST "\${KEPTN_API_ENDPOINT}/v1/event" \\
      -H "accept: application/json; charset=utf-8" -H "x-token: \${KEPTN_API_TOKEN}" -H "Content-Type: application/json; charset=utf-8" \\
      -d "{"type":"sh.keptn.event.problem.open","specversion":"0.2","source":"api","contenttype":"application\\/json","data":{"State":"OPEN","ProblemID":"\${PROBLEM_ID}","ProblemTitle":"\${PROBLEM}","project":"\${PROJECT}","stage":"\${STAGE}","service":"\${SERVICE}"}}"`
    });
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
  }

}
