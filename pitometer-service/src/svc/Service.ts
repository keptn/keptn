import 'reflect-metadata';
import { injectable, inject } from 'inversify';
import { RequestModel } from './RequestModel';
import axios, { AxiosRequestConfig, AxiosPromise, AxiosError } from 'axios';
import YAML from 'yaml';
const { base64encode, base64decode } = require('nodejs-base64');

const Pitometer = require('@keptn/pitometer').Pitometer;
// tslint:disable-next-line: variable-name
const PrometheusSource = require('@keptn/pitometer-source-prometheus').Source;
// tslint:disable-next-line: variable-name
const DynatraceSource = require('@keptn/pitometer-source-dynatrace').Source;
// tslint:disable-next-line: variable-name
const ThresholdGrader = require('@keptn/pitometer-grader-threshold').Grader;

import moment from 'moment';

import { Logger } from '../lib/Logger';
import { Keptn } from '../lib/Keptn';
import { Credentials } from '../lib/Credentials';
import { DynatraceCredentialsModel } from '../lib/DynatraceCredentialsModel';
import { ServiceIndicators, Indicator } from './ServiceIndicators';
import { ServiceObjectives, Objective } from './ServiceObjectives';
import { create } from 'domain';
import { json } from 'body-parser';

let configServiceUrl;

@injectable()
export class Service {

  constructor() { }

  public async handleRequest(event: RequestModel): Promise<boolean> {
    try {
      if (event.data.teststrategy !== 'performance' && event.data.teststrategy !== 'real-user') {
        Logger.log(
          event.shkeptncontext, event.id,
          `No performance gate specified for stage ${event.data.stage}`,
        );
        this.handleEvaluationResult({ result: 'pass' }, event);
        return true;
      }
      const pitometer = new Pitometer();

      let prometheusUrl;
      if (process.env.NODE_ENV === 'production') {
        configServiceUrl = 'http://configuration-service.keptn.svc.cluster.local:8080';
        prometheusUrl =
          `http://prometheus-service.monitoring.svc.cluster.local:8080/api/v1/query`;
      } else {
        configServiceUrl = 'http://localhost:6060';
        prometheusUrl = 'http://localhost:8080/api/v1/query';
      }

      let testRunDuration = 0;
      if (event.time === undefined) {
        event.time = moment().format();
      }
      if (event.time !== undefined && event.data.startedat !== undefined) {
        testRunDuration = Math.ceil(
          moment.duration(moment(event.time).diff(moment(event.data.startedat))).asSeconds(),
        );
      }

      pitometer.addGrader('Threshold', new ThresholdGrader());

      let perfspecString;

      try {
        perfspecString = await this.getPerfspecString(event);
      } catch (e) {
        this.handleEvaluationResult(
          {
            result: 'failed',
            error: 'Error while retrieving perfspec content.',
          },
          event,
        );
        return false;
      }
      if (perfspecString === '') {
        Logger.log(
          event.shkeptncontext, event.id,
          `No perfspec file defined for `
          + `${event.data.project}:${event.data.service}:${event.data.stage}`);
        this.handleEvaluationResult({ result: 'pass' }, event);
        return true;
      }

      let perfspec;

      Logger.log(
        event.shkeptncontext, event.id,
        perfspecString,
      );

      try {
        perfspec = JSON.parse(perfspecString);
      } catch (e) {
        this.handleEvaluationResult(
          {
            result: 'failed',
            error: 'Bad perfspec format.',
          },
          event,
        );
        return false;
      }

      const envPlaceHolderRegex = new RegExp('\\$ENVIRONMENT', 'g');
      perfspecString =
          perfspecString.replace(envPlaceHolderRegex, `${event.data.stage}`);
      /*
      /* TODO: going forward, setting the duration via the $DURATION_MINUTES
      /* placeholder will become obsolete, since this is handled by pitometer now.
      /* For backwards compatibility reasons we have to keep this for now.
      */
      let durationRegex = new RegExp('\\$DURATIONm', 'g');
      perfspecString = perfspecString.replace(durationRegex, '$DURATION');
      durationRegex = new RegExp('\\$DURATIONs', 'g');
      perfspecString = perfspecString.replace(durationRegex, '$DURATION');
      durationRegex = new RegExp('\\$DURATION', 'g');
      if (testRunDuration > 0) {
        perfspecString = perfspecString.replace(durationRegex, `${testRunDuration}s`);
      } else {
        perfspecString = perfspecString.replace(durationRegex, `10s`);
      }

      perfspec = JSON.parse(perfspecString);
      Logger.log(
        event.shkeptncontext, event.id,
        `Perfspec file content: ${JSON.stringify(perfspec)}`,
      );

      const indicators = [];
      if (perfspec === undefined || perfspec.indicators === undefined) {
        this.handleEvaluationResult(
          {
            result: 'failed',
            error: 'Bad perfspec format.',
          },
          event,
        );
        return false;
      }

      if (perfspec.indicators
        .find(indicator => indicator.source.toLowerCase() === 'prometheus') !== undefined) {
        this.addPrometheusSource(event, pitometer, prometheusUrl);
      }
      // get dynatrace service entity ID if Dynatrace source is defined in perfspec
      let serviceEntityId = '';
      if (perfspec.indicators
        .find(indicator => indicator.source.toLowerCase() === 'dynatrace') !== undefined) {
        try {
          const dynatraceCredentials = await this.addDynatraceSource(event, pitometer);
          serviceEntityId = await this.getDTServiceEntityId(event, dynatraceCredentials);

          if (serviceEntityId === undefined || serviceEntityId === '') {
            this.handleEvaluationResult(
              {
                result: 'failed',
                error: 'No Dynatrace Service Entity found.',
              },
              event,
            );
            return false;
          }
        } catch (e) {
          this.handleEvaluationResult(
            {
              result: 'failed',
              error: `Error while fetching Dynatrace data: ${e}`,
            },
            event,
          );
          return false;
        }
      }

      for (let i = 0; i < perfspec.indicators.length; i += 1) {
        const indicator = perfspec.indicators[i];
        if (indicator.source.toLowerCase() === 'dynatrace' && indicator.query !== undefined) {
          if (serviceEntityId !== undefined && serviceEntityId !== '') {
            indicator.query.entityIds = [serviceEntityId];
          }
        }
        indicators.push(indicator);
      }
      perfspec.indicators = indicators;
      try {
        const evaluationResult = await pitometer.run(
          perfspec,
          {
            timeStart: moment(event.data.startedat).unix(),
            timeEnd: moment(event.time).unix(),
          },
        );
        Logger.log(
          event.shkeptncontext, event.id,
          evaluationResult,
        );

        this.handleEvaluationResult(evaluationResult, event);
      } catch (e) {
        Logger.log(
          event.shkeptncontext,
          JSON.stringify(e.config.data),
          'ERROR',
        );
        this.handleEvaluationResult(
          {
            result: 'failed',
            error: `${e}`,
          },
          event,
        );
      }
      return true;
    } catch (e) {
      this.handleEvaluationResult(
        {
          result: 'failed',
          error: `${e}`,
        },
        event,
      );
    }
  }

  private addPrometheusSource(event: RequestModel, pitometer: any, prometheusUrl: any) {
    Logger.log(event.shkeptncontext, `Adding Prometheus source`, 'DEBUG');
    pitometer.addSource('Prometheus', new PrometheusSource({
      queryUrl: prometheusUrl,
    }));
  }

  private async addDynatraceSource(
    event: RequestModel, pitometer: any): Promise<DynatraceCredentialsModel> {
    const dynatraceCredentials = await Credentials.getInstance().getDynatraceCredentials();
    if (dynatraceCredentials !== undefined &&
      dynatraceCredentials.tenant !== undefined &&
      dynatraceCredentials.apiToken !== undefined) {
      Logger.log(
        event.shkeptncontext, event.id,
        `Adding Dynatrace Source for tenant ${dynatraceCredentials.tenant}`,
        'DEBUG',
      );
      pitometer.addSource('Dynatrace', new DynatraceSource({
        baseUrl: `https://${dynatraceCredentials.tenant}`,
        apiToken: dynatraceCredentials.apiToken,
        log: console.log,
      }));

      return dynatraceCredentials;
    }
    throw new Error('no Dynatrace credentials available in the cluster.');
  }

  private async getDTServiceEntityId(
    event: RequestModel, dynatraceCredentials: DynatraceCredentialsModel) {
    try {
      let entityId = '';
      for (let i = 0; i < 10; i += 1) {
        entityId = await this.getEntityId(event, dynatraceCredentials);
        if (entityId !== '') {
          return entityId;
        }
        await this.delay(20000);
      }
      return '';
    } catch (e) {
      throw e;
    }
  }

  private delay(t) {
    return new Promise((resolve, reject) => {
      setTimeout(resolve, t);
    });
  }

  private async getEntityId(
    event: RequestModel, dynatraceCredentials: DynatraceCredentialsModel): Promise<string> {
    let entityId = '';
    Logger.log(
      event.shkeptncontext, event.id,
      `Trying to get serviceEntityId with most requests ` +
      `during test run execution for ${event.data.service} ` +
      `in namespace ${event.data.project}-${event.data.stage}`,
    );
    try {
      const dtApiUrl =
        `https://${dynatraceCredentials.tenant}` +
        `/api/v1/timeseries/com.dynatrace.builtin%3Aservice.server_side_requests?` +
        `Api-Token=${dynatraceCredentials.apiToken}`;

      // TODO: check for real-user test strategy
      const data = {
        aggregationType: 'count',
        startTimestamp: moment(event.data.startedat).unix() * 1000,
        endTimestamp: moment(event.time).unix() * 1000,
        tags: [
          `service:${event.data.service}`,
          `environment:${event.data.project}-${event.data.stage}`,
        ],
        queryMode: 'TOTAL',
        timeseriesId: 'com.dynatrace.builtin:service.server_side_requests',
      };
      if (event.data.teststrategy !== 'real-user') {
        data.tags.push('test-subject:true');
      }
      const timeseries = await axios.post(dtApiUrl, data);
      if (timeseries.data &&
        timeseries.data.result &&
        timeseries.data.result.dataPoints) {
        let max = 0;
        for (const entity in timeseries.data.result.dataPoints) {
          const dataPoint = timeseries.data.result.dataPoints[entity][0];
          if (dataPoint.length > 1 && dataPoint[1] >= max) {
            entityId = entity;
            max = dataPoint[1];
          }
        }
      }
    } catch (e) {
      if (
        e.response !== undefined &&
        e.response.status !== undefined &&
        e.response.status === 400) {
        Logger.log(
          event.shkeptncontext,
          `No data in Dynatrace available yet.`,
          'DEBUG',
        );
        return '';
      }
      Logger.log(
        event.shkeptncontext,
        `Error while requesting serviceEntityId with most
        requests during test run execution: ${e}`,
        'ERROR',
      );
      throw e;
    }
    if (entityId !== undefined && entityId !== '') {
      Logger.log(
        event.shkeptncontext, event.id,
        `Found serviceEntityId: ${entityId}`,
      );
    } else {
      Logger.log(
        event.shkeptncontext, event.id,
        'No Dynatrace serviceEntityId found.',
      );
    }
    return entityId;
  }

  async getServiceIndicators(event: RequestModel): Promise<ServiceIndicators> {
    const indicatorString =
      await this.getServiceResourceContent(event, 'service-indicators.yaml');

    if (indicatorString === '') {
      Logger.log(
        event.shkeptncontext, event.id,
        `No service-indicators file defined for `
        + `${event.data.project}:${event.data.service}:${event.data.stage}`);
      return null;
    }

    Logger.log(
      event.shkeptncontext, event.id,
      `Service Indicator file content: ${indicatorString}`,
    );

    try {
      const indicators = <ServiceIndicators>YAML.parse(indicatorString);
      return indicators;
    } catch (e) {
      Logger.log(
        event.shkeptncontext, event.id,
        `Invalid Service indicators format: `
        + `${event.data.project}:${event.data.service}:${event.data.stage}`);
      throw e;
    }
  }

  async getServiceObjectives(event: RequestModel): Promise<ServiceObjectives> {
    const objectivesString =
      await this.getServiceResourceContent(event, 'service-objectives.yaml');

    if (objectivesString === '') {
      Logger.log(
        event.shkeptncontext, event.id,
        `No service-objectives file defined for `
        + `${event.data.project}:${event.data.service}:${event.data.stage}`);
      return null;
    }

    Logger.log(
      event.shkeptncontext, event.id,
      `Service Objectives file content: ${objectivesString}`,
    );

    try {
      const objectives = <ServiceObjectives>YAML.parse(objectivesString);
      return objectives;
    } catch (e) {
      Logger.log(
        event.shkeptncontext, event.id,
        `Invalid Service objectives format: `
        + `${event.data.project}:${event.data.service}:${event.data.stage}`);
      throw e;
    }
  }

  async getPerfspecString(event: RequestModel): Promise<string> {
    try {
      const indicators = await this.getServiceIndicators(event);
      if (indicators === null
        || indicators.indicators === undefined
        || indicators.indicators.length === 0
      ) {
        return await this.getServiceResourceContent(event, 'perfspec.json');
      }
      const objectives = await this.getServiceObjectives(event);
      if (objectives === null
        || objectives.objectives === undefined
        || objectives.objectives.length === 0
      ) {
        return await this.getServiceResourceContent(event, 'perfspec.json');
      }
      const perfspecObject = this.createPerfspecObject(indicators, objectives);
      return JSON.stringify(perfspecObject);
    } catch (e) {
      throw e;
    }
  }

  createPerfspecObject(indicators: ServiceIndicators, objectives: ServiceObjectives): any {
    const perfspecObject = {
      spec_version: '1.0',
      indicators: [],
      objectives: {
        pass: objectives.pass,
        warning: objectives.warning,
      },
    };

    const getIndicator = (name: string): Indicator => {
      for (let i = 0; i < indicators.indicators.length; i += 1) {
        if (indicators.indicators[i].metric === name) {
          return indicators.indicators[i];
        }
      }
      return null;
    };

    for (let i = 0; i < objectives.objectives.length; i += 1) {
      const objective = objectives.objectives[i];
      const indicator = getIndicator(objective.metric);
      if (indicator !== null) {
        const newPerfspecIndicator = {
          id: indicator.metric,
          source: indicator.source,
          query: {},
          grading: {
            type: 'Threshold',
            thresholds: {
              upperSevere: objective.threshold,
            },
            metricScore: objective.score,
          },
        };

        if (indicator.queryObject !== undefined && indicator.queryObject.length > 0) {
          const newQueryObject = {};
          for (let j = 0; j < indicator.queryObject.length; j += 1) {
            newQueryObject[indicator.queryObject[j].key] = indicator.queryObject[j].value;
          }
          newPerfspecIndicator.query = newQueryObject;
        } else {
          newPerfspecIndicator.query = indicator.query;
        }
        perfspecObject.indicators.push(newPerfspecIndicator);
      }
    }

    return perfspecObject;
  }

  async getServiceResourceContent(event: RequestModel, resourceUri: string): Promise<string> {
    // tslint:disable-next-line: max-line-length
    const url = `${configServiceUrl}/v1/project/${event.data.project}/stage/${event.data.stage}/service/${event.data.service}/resource/${resourceUri}`;
    let response;

    try {
      response = await axios.get(url, {});
    } catch (e) {
      Logger.log(
        event.shkeptncontext, event.id,
        `Resource ${resourceUri} not found in`
        + `${event.data.project}:${event.data.service}:${event.data.stage}`);
      return '';
    }
    if (response.data !== undefined && response.data.resourceContent !== undefined) {
      // decode the base64 encoded string
      return base64decode(response.data.resourceContent);
    }
    return '';
  }

  async handleEvaluationResult(evaluationResult: any, sourceEvent: RequestModel): Promise<void> {
    const evaluationPassed: boolean =
      evaluationResult.result !== undefined &&
      (evaluationResult.result === 'pass' || evaluationResult.result === 'warning');

    Logger.log(
      sourceEvent.shkeptncontext,
      sourceEvent.id,
      `Evaluation passed: ${evaluationPassed}`,
    );
    try {
      Logger.log(
        sourceEvent.shkeptncontext, sourceEvent.id,
        `Pitometer Result: ${JSON.stringify(evaluationResult)}`,
      );
    } catch (e) {
      Logger.log(
        sourceEvent.shkeptncontext,
        e,
        'ERROR',
      );
    }

    const event: RequestModel = new RequestModel();
    event.type = RequestModel.EVENT_TYPES.EVALUATION_DONE;
    event.source = 'pitometer-service';
    event.shkeptncontext = sourceEvent.shkeptncontext;
    event.data.githuborg = sourceEvent.data.githuborg;
    event.data.project = sourceEvent.data.project;
    event.data.teststrategy = sourceEvent.data.teststrategy;
    event.data.deploymentstrategy = sourceEvent.data.deploymentstrategy;
    event.data.stage = sourceEvent.data.stage;
    event.data.service = sourceEvent.data.service;
    event.data.image = sourceEvent.data.image;
    event.data.tag = sourceEvent.data.tag;
    event.data.evaluationpassed = evaluationPassed;
    event.data.evaluationdetails = evaluationResult;

    Keptn.sendEvent(event);
  }
}
