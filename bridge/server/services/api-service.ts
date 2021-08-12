import { axios } from './axios-instance';
import { AxiosResponse } from 'axios';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Project } from '../models/project';
import { ResultTypes } from '../../shared/models/result-types';
import { SequenceResult } from '../interfaces/sequence-result';
import { EventResult } from '../interfaces/event-result';
import { UniformRegistration } from '../interfaces/uniform-registration';
import { UniformRegistrationLogResponse } from '../interfaces/uniform-registration-log';

export class ApiService {
  private readonly defaultHeaders: object;
  constructor(private readonly baseUrl: string, private readonly apiToken: string) {
    this.defaultHeaders = {
      'x-token': apiToken,
      'Content-Type': 'application/json'
    };
  }


  public getProject(projectName: string): Promise<AxiosResponse<Project>> {
    return axios.get<Project>(`${this.baseUrl}/controlPlane/v1/project/${projectName}`, {
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

  public getOpenTriggeredEvents(projectName: string, stageName: string, serviceName: string, eventType: EventTypes): Promise<AxiosResponse<EventResult>> {
    const params = {
      project: projectName,
      stage: stageName,
      service: serviceName,
    };
    return axios.get<EventResult>(`${this.baseUrl}/controlPlane/v1/event/triggered/${eventType}`, {params, headers: this.defaultHeaders});
  }

  public getUniformRegistrations(): Promise<AxiosResponse<UniformRegistration[]>> {
    return axios.get<UniformRegistration[]>(`${this.baseUrl}/controlPlane/v1/uniform/registration`, {headers: this.defaultHeaders});
  }

  public getUniformRegistrationLogs(integrationId: string, fromTime?: string, pageSize = 100): Promise<AxiosResponse<UniformRegistrationLogResponse>> {
    const params = {
      integrationId,
      ...fromTime && {fromTime: new Date(new Date(fromTime).getTime() + 1).toISOString()}, // > fromTime instead of >= fromTime
      pageSize: pageSize.toString()
    };
    return axios.get<UniformRegistrationLogResponse>(`${this.baseUrl}/controlPlane/v1/log`, {params, headers: this.defaultHeaders});
  }
}
