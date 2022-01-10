import Axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Project } from '../models/project';
import { ResultTypes } from '../../shared/models/result-types';
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
import { SequenceState } from '../../shared/models/sequence';
import { IStage } from '../../shared/interfaces/stage';

export class ApiService {
  private readonly axios: AxiosInstance;
  private readonly escapeSlash = '%252F';

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
    pageSize: number,
    sequenceName?: string,
    state?: SequenceState,
    fromTime?: string,
    beforeTime?: string,
    keptnContext?: string
  ): Promise<AxiosResponse<SequenceResult>> {
    const params: { [key: string]: string } = {
      pageSize: pageSize.toString(),
      ...(sequenceName && { name: sequenceName }),
      ...(state && { state }),
      ...(fromTime && { fromTime }),
      ...(beforeTime && { beforeTime }),
      ...(keptnContext && { keptnContext }),
    };

    return this.axios.get<SequenceResult>(`${this.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getTraces(
    accessToken: string | undefined,
    eventType?: string,
    pageSize?: number,
    projectName?: string,
    stageName?: string,
    serviceName?: string,
    keptnContext?: string,
    eventSource?: KeptnService
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      ...(projectName && { project: projectName }),
      ...(serviceName && { service: serviceName }),
      ...(stageName && { stage: stageName }),
      ...(eventType && { type: eventType }),
      ...(pageSize && { pageSize: pageSize.toString() }),
      ...(keptnContext && { keptnContext }),
      ...(eventSource && { source: eventSource }),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, {
      params,
      ...this.getAuthHeaders(accessToken),
    });
  }

  public getTracesByContext(
    accessToken: string | undefined,
    keptnContext: string | undefined,
    projectName?: string | undefined,
    fromTime?: string | undefined,
    nextPageKey?: string | undefined
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      keptnContext,
      pageSize: '100',
      ...(nextPageKey && { nextPageKey }),
      ...(projectName && { project: projectName }),
      ...(fromTime && { fromTime }),
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
    eventType: EventTypes,
    pageSize: number,
    projectName: string,
    stageName: string,
    serviceName: string,
    resultType?: ResultTypes,
    source?: KeptnService
  ): Promise<AxiosResponse<EventResult>> {
    const resultString = resultType ? ` AND data.result:${resultType}` : '';
    const sourceString = source ? ` AND data.source:${source}` : '';
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName}${sourceString}${resultString}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString(),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${eventType}`, {
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
