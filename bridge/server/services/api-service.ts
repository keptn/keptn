import { axios } from './axios-instance.js';
import { SequenceResult } from '../interfaces/sequence-result.js';
import { AxiosResponse } from 'axios';
import { ProjectResult } from '../interfaces/project-result.js';
import { EventResult } from '../interfaces/event-result.js';
import { EventTypes } from '../models/event-types.js';
import { ResultTypes } from '../models/result-types.js';

export class ApiService {
  private readonly defaultHeaders: object;
  constructor(private readonly baseUrl: string, private readonly apiToken: string) {
    this.defaultHeaders = {
      'x-token': apiToken,
      'Content-Type': 'application/json'
    };
  }


  public getProject(projectName: string): Promise<AxiosResponse<ProjectResult>> {
    return axios.get<ProjectResult>(`${this.baseUrl}/controlPlane/v1/project/${projectName}`, {
      headers: this.defaultHeaders
    });
  }
  public getSequences(projectName: string, pageSize: number, sequenceName?: string, state?: string,
                      fromTime?: string, beforeTime?: string, keptnContext?: string): Promise<AxiosResponse<SequenceResult>> {
    const params = {
      pageSize: pageSize.toString(),
      ...(sequenceName && {name: sequenceName}),
      ...(state && {state}),
      ...(fromTime && {fromTime}),
      ...(beforeTime && {beforeTime}),
      ...(keptnContext && {keptnContext})
    };

    return axios.get<SequenceResult>(`${this.baseUrl}/controlPlane/v1/sequence/${projectName}`, {params, headers: this.defaultHeaders});
  }

  public getTraces(eventType: string, pageSize: number, projectName: string, stageName: string, serviceName: string, keptnContext?: string, fromTime?: string): Promise<AxiosResponse<EventResult>> {
    const params = {
      project: projectName,
      service: serviceName,
      stage: stageName,
      type: eventType,
      limit: pageSize.toString(),
      ...keptnContext && {keptnContext}
    };
    return axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, {params, headers: this.defaultHeaders});
  }

  public getTracesWithResult(eventType: EventTypes, pageSize: number, projectName: string, stageName: string, serviceName: string, resultType: ResultTypes): Promise<AxiosResponse<EventResult>> {
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName} AND data.result:${resultType}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString()
    };
    return axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${eventType}`, {params, headers: this.defaultHeaders});
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, pageSize: number, keptnContext?: string): Promise<AxiosResponse<EventResult>> {
    const contextString = keptnContext ? ` AND shkeptncontext:${keptnContext}` : '';
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName}${contextString}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString()
    };
    return axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {params, headers: this.defaultHeaders});
  }
}
