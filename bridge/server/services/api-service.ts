import Axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Project } from '../models/project';
import { SequenceResult } from '../interfaces/sequence-result';
import { EventResult } from '../interfaces/event-result';
import { UniformRegistration } from '../models/uniform-registration';
import { UniformRegistrationLogResponse } from '../../shared/interfaces/uniform-registration-log';
import { Resource, ResourceResponse } from '../../shared/interfaces/resource';
import https from 'https';
import { ProjectResult } from '../interfaces/project-result';
import { UniformSubscription } from '../../shared/interfaces/uniform-subscription';
import { Secret } from '../../shared/interfaces/secret';
import { KeptnService } from '../../shared/models/keptn-service';
import { IStage } from '../../shared/interfaces/stage';
import { SequenceOptions, TraceOptions } from './data-service';
import { ComponentLogger } from '../utils/logger';

export class ApiService {
  private readonly axios: AxiosInstance;
  private readonly escapeSlash = '%252F';
  private readonly log = new ComponentLogger('API');

  constructor(private readonly baseUrl: string, readonly apiToken: string | undefined) {
    if (process.env.NODE_ENV === 'test' && global.axiosInstance) {
      this.axios = global.axiosInstance;
    } else {
      this.axios = Axios.create({
        // accepts self-signed ssl certificate
        httpsAgent: new https.Agent({
          rejectUnauthorized: false,
        }),
        headers: {
          ...(apiToken && { 'x-token': apiToken }),
          'Content-Type': 'application/json',
        },
      });
      // log request using the same format of morgan but without response time
      this.axios.interceptors.response.use((response) => {
        //Example: GET /api/project/podtato 200
        this.log.info(`${response.config.method} ${response.config.url} ${response.status}`);
        return response;
      });
    }
  }

  public getProjects(accessToken: string | undefined): Promise<AxiosResponse<ProjectResult>> {
    return this.axios.get<ProjectResult>(`${this.baseUrl}/controlPlane/v1/project`, this.getAuthHeaders(accessToken));
  }

  public getProject(accessToken: string | undefined, projectName: string): Promise<AxiosResponse<Project>> {
    return this.axios.get<Project>(
      `${this.baseUrl}/controlPlane/v1/project/${projectName}`,
      this.getAuthHeaders(accessToken)
    );
  }

  public getSequences(
    accessToken: string | undefined,
    projectName: string,
    sequenceOptions: SequenceOptions
  ): Promise<AxiosResponse<SequenceResult>> {
    this.log.debug(`Sequence options: ${this.log.prettyPrint(sequenceOptions)}`);
    return this.axios.get<SequenceResult>(`${this.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
      params: sequenceOptions,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getTraces(
    accessToken: string | undefined,
    traceOptions: Partial<TraceOptions>
  ): Promise<AxiosResponse<EventResult>> {
    this.log.debug(`Trace options: ${this.log.prettyPrint(traceOptions)}`);
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, {
      params: traceOptions,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getTracesByContext(
    accessToken: string | undefined,
    keptnContext: string | undefined,
    projectName?: string | undefined,
    fromTime?: string | undefined,
    nextPageKey?: string | undefined,
    type?: EventTypes,
    source?: KeptnService
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      keptnContext,
      pageSize: '100',
      ...(nextPageKey && { nextPageKey }),
      ...(projectName && { project: projectName }),
      ...(fromTime && { fromTime }),
      ...(type && { type }),
      ...(source && { source }),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getRoots(
    accessToken: string | undefined,
    projectName: string | undefined,
    pageSize: string | undefined,
    serviceName: string | undefined,
    fromTime?: string | undefined,
    beforeTime?: string | undefined,
    keptnContext?: string | undefined
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      root: 'true',
      project: projectName,
      limit: pageSize,
      ...(serviceName && { serviceName }),
      ...(fromTime && { fromTime }),
      ...(beforeTime && { beforeTime }),
      ...(keptnContext && { keptnContext }),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getTracesWithResultAndSource(
    accessToken: string | undefined,
    traceOptions: TraceOptions
  ): Promise<AxiosResponse<EventResult>> {
    const resultString = traceOptions.result ? ` AND data.result:${traceOptions.result}` : '';
    const sourceString = traceOptions.source ? ` AND source:${traceOptions.source}` : '';
    const params = {
      filter: `data.project:${traceOptions.project} AND data.service:${traceOptions.service} AND data.stage:${traceOptions.stage}${sourceString}${resultString}`,
      excludeInvalidated: 'true',
      limit: traceOptions.pageSize,
    };
    this.log.debug(`getTracesWithResultAndSource options: ${this.log.prettyPrint(params)}`);
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${traceOptions.type}`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getEvaluationResults(
    accessToken: string | undefined,
    projectName: string,
    serviceName: string,
    stageName: string,
    pageSize: number,
    keptnContext?: string
  ): Promise<AxiosResponse<EventResult>> {
    const contextString = keptnContext ? ` AND shkeptncontext:${keptnContext}` : '';
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName}${contextString} AND source:${KeptnService.LIGHTHOUSE_SERVICE}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString(),
    };
    this.log.debug(`getEvaluationResults options: ${this.log.prettyPrint(params)}`);
    return this.axios.get<EventResult>(
      `${this.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`,
      { params, ...this.getAuthHeaders(accessToken) }
    );
  }

  public getEvaluationResult(
    accessToken: string | undefined,
    keptnContext: string
  ): Promise<AxiosResponse<EventResult>> {
    const url = `${this.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`;
    const params = {
      filter: `shkeptncontext:${keptnContext} AND source:${KeptnService.LIGHTHOUSE_SERVICE}`,
      limit: '1',
    };
    this.log.debug(`getEvaluationResult options: ${this.log.prettyPrint(params)}`);
    return this.axios.get<EventResult>(url, { params, ...this.getAuthHeaders(accessToken) });
  }

  public getTracesOfMultipleServices(
    accessToken: string | undefined,
    projectName: string,
    eventType: EventTypes,
    eventIds: string,
    source?: KeptnService
  ): Promise<AxiosResponse<EventResult>> {
    const sourceString = source ? `AND source:${source} ` : '';
    const params = {
      filter: `data.project:${projectName} ${sourceString}AND id:${eventIds}`,
      excludeInvalidated: 'true',
    };
    this.log.debug(`getTracesOfMultipleServices options: ${this.log.prettyPrint(params)}`);
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${eventType}`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getOpenTriggeredEvents(
    accessToken: string | undefined,
    projectName: string,
    eventType: EventTypes,
    stageName?: string,
    serviceName?: string
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      project: projectName,
      ...(stageName && { stage: stageName }),
      ...(serviceName && { service: serviceName }),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/controlPlane/v1/event/triggered/${eventType}`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getUniformRegistrations(
    accessToken: string | undefined,
    integrationId?: string
  ): Promise<AxiosResponse<UniformRegistration[]>> {
    return this.axios.get<UniformRegistration[]>(`${this.baseUrl}/controlPlane/v1/uniform/registration`, {
      params: {
        ...(integrationId && { id: integrationId }),
      },
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getUniformRegistrationLogs(
    accessToken: string | undefined,
    integrationId: string,
    fromTime?: string,
    pageSize = 100
  ): Promise<AxiosResponse<UniformRegistrationLogResponse>> {
    const params = {
      integrationId,
      ...(fromTime && { fromTime: new Date(new Date(fromTime).getTime() + 1).toISOString() }), // > fromTime instead of >= fromTime
      pageSize: pageSize.toString(),
    };
    return this.axios.get<UniformRegistrationLogResponse>(`${this.baseUrl}/controlPlane/v1/log`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getShipyard(accessToken: string | undefined, projectName: string): Promise<AxiosResponse<Resource>> {
    return this.axios.get<Resource>(
      `${this.baseUrl}/configuration-service/v1/project/${projectName}/resource/shipyard.yaml`,
      this.getAuthHeaders(accessToken)
    );
  }

  public createSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscription: UniformSubscription
  ): Promise<AxiosResponse<{ id: string }>> {
    return this.axios.post(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription`,
      subscription,
      this.getAuthHeaders(accessToken)
    );
  }

  public updateSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscriptionId: string,
    subscription: UniformSubscription
  ): Promise<AxiosResponse<Record<string, unknown>>> {
    return this.axios.put(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`,
      subscription,
      this.getAuthHeaders(accessToken)
    );
  }

  public getWebhookConfig(
    accessToken: string | undefined,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<AxiosResponse<Resource>> {
    let url = `${this.baseUrl}/configuration-service/v1/project/${projectName}`;
    if (stageName) {
      url += `/stage/${stageName}`;
    }
    if (serviceName) {
      url += `/service/${serviceName}`;
    }
    url += `/resource/webhook${this.escapeSlash}webhook.yaml`;
    this.log.debug(`getWebhookConfig webhook path: ${url}`);
    return this.axios.get<Resource>(url, this.getAuthHeaders(accessToken));
  }

  public deleteWebhookConfig(
    accessToken: string | undefined,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<AxiosResponse<Resource>> {
    let url = `${this.baseUrl}/configuration-service/v1/project/${projectName}`;
    if (stageName) {
      url += `/stage/${stageName}`;
    }
    if (serviceName) {
      url += `/service/${serviceName}`;
    }
    url += `/resource/webhook${this.escapeSlash}webhook.yaml`;
    this.log.debug(`deleteWebhookConfig webhook path: ${url}`);
    return this.axios.delete<Resource>(url, this.getAuthHeaders(accessToken));
  }

  public saveWebhookConfig(
    accessToken: string | undefined,
    content: string,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<AxiosResponse<Resource>> {
    let url = `${this.baseUrl}/configuration-service/v1/project/${projectName}`;
    if (stageName) {
      url += `/stage/${stageName}`;
    }
    if (serviceName) {
      url += `/service/${serviceName}`;
    }
    url += `/resource`; // /resource/resourceURI does not overwrite, fallback to this endpoint
    this.log.debug(`saveWebhookConfig webhook path: ${url}`);
    return this.axios.put<Resource>(
      url,
      {
        resources: [
          {
            resourceURI: 'webhook/webhook.yaml',
            resourceContent: Buffer.from(content).toString('base64'),
          },
        ],
      },
      this.getAuthHeaders(accessToken)
    );
  }

  public deleteUniformSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscriptionId: string
  ): Promise<AxiosResponse<Record<string, unknown>>> {
    return this.axios.delete(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`,
      this.getAuthHeaders(accessToken)
    );
  }

  public getUniformSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscriptionId: string
  ): Promise<AxiosResponse<UniformSubscription>> {
    return this.axios.get<UniformSubscription>(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`,
      this.getAuthHeaders(accessToken)
    );
  }

  public getServiceResources(
    accessToken: string | undefined,
    projectName: string,
    stageName: string,
    serviceName: string,
    nextPageKey?: string
  ): Promise<AxiosResponse<ResourceResponse>> {
    const url = `${this.baseUrl}/configuration-service/v1/project/${projectName}/stage/${stageName}/service/${serviceName}/resource`;
    const params: { [key: string]: string } = {};
    if (nextPageKey) {
      params.nextPageKey = nextPageKey;
    }

    return this.axios.get<ResourceResponse>(url, { params, ...this.getAuthHeaders(accessToken) });
  }

  public getServiceResource(
    accessToken: string | undefined,
    projectName: string,
    stageName: string,
    serviceName: string,
    resourceURI: string
  ): Promise<AxiosResponse<Resource>> {
    const url = `${this.baseUrl}/configuration-service/v1/project/${projectName}/stage/${stageName}/service/${serviceName}/resource/${resourceURI}`;

    return this.axios.get<Resource>(url, this.getAuthHeaders(accessToken));
  }

  public getSecrets(accessToken: string | undefined): Promise<AxiosResponse<{ Secrets: Secret[] }>> {
    const url = `${this.baseUrl}/secrets/v1/secret`;
    return this.axios.get<{ Secrets: Secret[] }>(url, this.getAuthHeaders(accessToken));
  }

  public getStages(accessToken: string | undefined, projectName: string): Promise<AxiosResponse<{ stages: IStage[] }>> {
    const url = `${this.baseUrl}/controlPlane/v1/project/${projectName}/stage`;
    return this.axios.get(url, this.getAuthHeaders(accessToken));
  }

  private getAuthHeaders(accessToken: string | undefined): { headers: AxiosRequestConfig['headers'] } | undefined {
    return accessToken
      ? {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      : undefined;
  }
}
