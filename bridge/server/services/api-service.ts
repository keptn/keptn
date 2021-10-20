import Axios, { AxiosInstance, AxiosResponse } from 'axios';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Project } from '../models/project';
import { ResultTypes } from '../../shared/models/result-types';
import { SequenceResult } from '../interfaces/sequence-result';
import { EventResult } from '../interfaces/event-result';
import { UniformRegistration } from '../models/uniform-registration';
import { UniformRegistrationLogResponse } from '../interfaces/uniform-registration-log';
import { Resource, ResourceResponse } from '../../shared/interfaces/resource';
import https from 'https';
import { ProjectResult } from '../interfaces/project-result';
import { UniformSubscription } from '../../shared/interfaces/uniform-subscription';
import { Secret } from '../../shared/interfaces/secret';
import { KeptnService } from '../../shared/models/keptn-service';
import { SequenceState } from '../../shared/models/sequence';

export class ApiService {
  private readonly axios: AxiosInstance;
  private readonly escapeSlash = '%252F';

  constructor(private readonly baseUrl: string, private readonly apiToken: string) {
    if (process.env.NODE_ENV === 'test' && global.axiosInstance) {
      this.axios = global.axiosInstance;
    } else {
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
  }

  public getProjects(): Promise<AxiosResponse<ProjectResult>> {
    return this.axios.get<ProjectResult>(`${this.baseUrl}/controlPlane/v1/project`);
  }

  public getProject(projectName: string): Promise<AxiosResponse<Project>> {
    return this.axios.get<Project>(`${this.baseUrl}/controlPlane/v1/project/${projectName}`);
  }

  public getSequences(
    projectName: string,
    pageSize: number,
    sequenceName?: string,
    state?: SequenceState,
    fromTime?: string,
    beforeTime?: string,
    keptnContext?: string
  ): Promise<AxiosResponse<SequenceResult>> {
    const params = {
      pageSize: pageSize.toString(),
      ...(sequenceName && { name: sequenceName }),
      ...(state && { state }),
      ...(fromTime && { fromTime }),
      ...(beforeTime && { beforeTime }),
      ...(keptnContext && { keptnContext }),
    };

    return this.axios.get<SequenceResult>(`${this.baseUrl}/controlPlane/v1/sequence/${projectName}`, { params });
  }

  public getTraces(
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
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, { params });
  }

  public getTracesByContext(
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
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, { params });
  }

  public getRoots(
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
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event`, { params });
  }

  public getTracesWithResult(
    eventType: EventTypes,
    pageSize: number,
    projectName: string,
    stageName: string,
    serviceName: string,
    resultType: ResultTypes
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName} AND data.result:${resultType}`,
      excludeInvalidated: 'true',
      limit: pageSize.toString(),
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/mongodb-datastore/event/type/${eventType}`, { params });
  }

  public getEvaluationResults(
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
      { params }
    );
  }

  public getOpenTriggeredEvents(
    projectName: string,
    stageName: string,
    serviceName: string,
    eventType: EventTypes
  ): Promise<AxiosResponse<EventResult>> {
    const params = {
      project: projectName,
      stage: stageName,
      service: serviceName,
    };
    return this.axios.get<EventResult>(`${this.baseUrl}/controlPlane/v1/event/triggered/${eventType}`, { params });
  }

  public getUniformRegistrations(integrationId?: string): Promise<AxiosResponse<UniformRegistration[]>> {
    return this.axios.get<UniformRegistration[]>(`${this.baseUrl}/controlPlane/v1/uniform/registration`, {
      params: {
        ...(integrationId && { id: integrationId }),
      },
    });
  }

  public getUniformRegistrationLogs(
    integrationId: string,
    fromTime?: string,
    pageSize = 100
  ): Promise<AxiosResponse<UniformRegistrationLogResponse>> {
    const params = {
      integrationId,
      ...(fromTime && { fromTime: new Date(new Date(fromTime).getTime() + 1).toISOString() }), // > fromTime instead of >= fromTime
      pageSize: pageSize.toString(),
    };
    return this.axios.get<UniformRegistrationLogResponse>(`${this.baseUrl}/controlPlane/v1/log`, { params });
  }

  public getShipyard(projectName: string): Promise<AxiosResponse<Resource>> {
    return this.axios.get<Resource>(
      `${this.baseUrl}/configuration-service/v1/project/${projectName}/resource/shipyard.yaml`
    );
  }

  public createSubscription(
    integrationId: string,
    subscription: UniformSubscription
  ): Promise<AxiosResponse<{ id: string }>> {
    return this.axios.post(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription`,
      subscription
    );
  }

  public updateSubscription(
    integrationId: string,
    subscriptionId: string,
    subscription: UniformSubscription
  ): Promise<AxiosResponse<Record<string, unknown>>> {
    return this.axios.put(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`,
      subscription
    );
  }

  public getWebhookConfig(
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

    return this.axios.get<Resource>(url);
  }

  public deleteWebhookConfig(
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

    return this.axios.delete<Resource>(url);
  }

  public saveWebhookConfig(
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

    return this.axios.put<Resource>(url, {
      resources: [
        {
          resourceURI: 'webhook/webhook.yaml',
          resourceContent: Buffer.from(content).toString('base64'),
        },
      ],
    });
  }

  public deleteUniformSubscription(
    integrationId: string,
    subscriptionId: string
  ): Promise<AxiosResponse<Record<string, unknown>>> {
    return this.axios.delete(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`
    );
  }

  public getUniformSubscription(
    integrationId: string,
    subscriptionId: string
  ): Promise<AxiosResponse<UniformSubscription>> {
    return this.axios.get<UniformSubscription>(
      `${this.baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`
    );
  }

  public getServiceResource(
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

    return this.axios.get<ResourceResponse>(url, params);
  }

  public getSecrets(): Promise<AxiosResponse<{ Secrets: Secret[] }>> {
    const url = `${this.baseUrl}/secrets/v1/secret`;
    return this.axios.get<{ Secrets: Secret[] }>(url);
  }
}
