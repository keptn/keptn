import * as express from 'express';
import { inject, injectable } from 'inversify';
import {
  controller,
  httpGet,
  httpPost,
  interfaces,
  httpDelete,
} from 'inversify-express-utils';
import {
  ApiOperationGet,
  ApiOperationPost,
  ApiPath,
  SwaggerDefinitionConstant,
  ApiOperationDelete,
} from 'swagger-express-ts';

import { MessageService } from '../svc/MessageService';
import { WebSocketService } from '../svc/WebSocketService';

const uuidv4 = require('uuid/v4');

@ApiPath({
  name: 'Project',
  path: '/project',
  security: { apiKeyHeader: [] },
})
@controller('/project')
export class ProjectController implements interfaces.Controller {

  @inject('MessageService') private readonly messageService: MessageService;

  constructor() { }

  @ApiOperationPost({
    description: 'Create a new project',
    parameters: {
      body: {
        description: 'Project information',
        model: 'ProjectRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Create a new keptn project',
  })
  @httpPost('/')
  public async createProject(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    const keptnContext = uuidv4();
    const result = {
      keptnContext,
      success: true,
    };
    const channelInfo = await WebSocketService.getInstance().createChannel(keptnContext);
    if (request.body && request.body.data !== undefined) {
      request.body.data.channelInfo = channelInfo;
      request.body.shkeptncontext = keptnContext;
    }
    result.success = await this.messageService.sendMessage(request.body);
    response.send({
      success: result,
      websocketChannel: channelInfo,
    });
  }

  @ApiOperationGet({
    description: 'Get projects',
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Get keptn projects',
  })
  @httpGet('/')
  public getProjects(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ) {
    const result = {
      result: 'success',
    };

    response.send(result);
  }

  @ApiOperationDelete({
    description: 'Delete a project',
    parameters: {

    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Delete a keptn project',
  })
  @httpDelete('/')
  public deleteProject(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ) {
    const result = {
      result: 'success',
    };

    response.send(result);
  }
}
