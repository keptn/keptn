import Axios, { AxiosInstance, AxiosResponse } from 'axios';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Project } from '../models/project';
import { ResultTypes } from '../../shared/models/result-types';
import { SequenceResult } from '../interfaces/sequence-result';
import { EventResult } from '../interfaces/event-result';
import { UniformRegistration } from '../interfaces/uniform-registration';
import { UniformRegistrationLogResponse } from '../interfaces/uniform-registration-log';
import { Resource, ResourceResponse } from '../../shared/interfaces/resource';
import https from 'https';

export class ApiService {
  private readonly axios: AxiosInstance;

  constructor(private readonly baseUrl: string, private readonly apiToken: string) {
    this.axios = Axios.create({
      // accepts self-signed ssl certificate
      httpsAgent: new https.Agent({
        rejectUnauthorized: false,
      }),
      headers: {
        'x-token': apiToken,
        'Content-Type': 'application/json',
      },
    });
  }


  public getProject(projectName: string): Promise<AxiosResponse<Project>> {
    return this.axios.get<Project>(`${this.baseUrl}/controlPlane/v1/project/${projectName}`);
  }

  public getSequences(projectName: string, pageSize: number, sequenceName?: string, state?: string,
                      fromTime?: string, beforeTime?: string, keptnContext?: string): Promise<AxiosResponse<SequenceResult>> {
    const params = {
      pageSize: pageSize.toString(),
      ...(sequenceName && {name: sequenceName}),
      ...(state && {state}),
      ...(fromTime && {fromTime}),
      ...(beforeTime && {beforeTime}),
      ...(keptnContext && {keptnContext}),
    };

    return this.axios.get<SequenceResult>(`${this.baseUrl}/controlPlane/v1/sequence/${projectName}`, {params});
  }

  public getTraces(eventType: string, pageSize: number, projectName: string, stageName: string, serviceName: string, keptnContext?: string, fromTime?: string): Promise<AxiosResponse<EventResult>> {
    const params = {
      project: projectName,
      service: serviceName,
      stage: stageName,
      type: eventType,
      limit: pageSize.toString(),
      ...keptnContext && {keptnContext},
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, {params});
  }

  public getTracesWithResult(eventType: EventTypes, pageSize: number, projectName: string, stageName: string, serviceName: string, resultType: ResultTypes): Promise<AxiosResponse<EventResult>> {
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName} AND data.result:${resultType}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString(),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${eventType}`, {params});
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, pageSize: number, keptnContext?: string): Promise<AxiosResponse<EventResult>> {
    const contextString = keptnContext ? ` AND shkeptncontext:${keptnContext}` : '';
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName}${contextString}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString(),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {params});
  }

  public getOpenTriggeredEvents(projectName: string, stageName: string, serviceName: string, eventType: EventTypes): Promise<AxiosResponse<EventResult>> {
    const params = {
      project: projectName,
      stage: stageName,
      service: serviceName,
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/controlPlane/v1/event/triggered/${eventType}`, {params});
  }

  public getUniformRegistrations(integrationId?: string): Promise<AxiosResponse<UniformRegistration[]>> {
    return this.axios.get<UniformRegistration[]>(`${this.baseUrl}/controlPlane/v1/uniform/registration`, {
      params: {
        ...integrationId && {id: integrationId},
      },
    });
  }

  public getUniformRegistrationLogs(integrationId: string, fromTime?: string, pageSize = 100): Promise<AxiosResponse<UniformRegistrationLogResponse>> {
    const params = {
      integrationId,
      ...fromTime && {fromTime: new Date(new Date(fromTime).getTime() + 1).toISOString()}, // > fromTime instead of >= fromTime
      pageSize: pageSize.toString(),
    };
    return this.axios.get<UniformRegistrationLogResponse>(`${this.baseUrl}/controlPlane/v1/log`, {params});
  }

  public getShipyard(projectName: string): Promise<AxiosResponse<Resource>> {
    return this.axios.get<Resource>(`${this.baseUrl}/configuration-service/v1/project/${projectName}/resource/shipyard.yaml`);
  }

  public getServiceResource(projectName: string, stageName: string, serviceName: string, nextPageKey?: string): Promise<AxiosResponse<ResourceResponse>> {
    const url = `${this.baseUrl}/configuration-service/v1/project/${projectName}/stage/${stageName}/service/${serviceName}/resource`;
    const params: { [key: string]: string } = {};
    if (nextPageKey) {
      params.nextPageKey = nextPageKey;
    }

    return this.axios.get<ResourceResponse>(url, params);
  }
}
