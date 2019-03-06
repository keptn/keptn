import moment from 'moment';

export class KeptnRequestModel {

  static EVENT_TYPES = {
    NEW_ARTEFACT: 'sh.keptn.events.new-artefact',
    CONFIGURATION_CHANGED: 'sh.keptn.events.configuration-changed',
    DEPLOYMENT_FINISHED: 'sh.keptn.events.deployment-finished',
    TESTS_FINISHED: 'sh.keptn.events.tests-finished',
    EVALUATION_DONE: 'sh.keptn.events.evaluation-done',
    PROBLEM: 'sh.keptn.events.problem',
  };

  public specversion: string;
  public type: string;
  public source: string;
  public id: string;
  public time: string;
  public datacontenttype: string;
  public data: any;

  constructor() {
    this.specversion = '0.2';
    this.time = moment().format();
    this.datacontenttype = 'application/json';
  }
}
