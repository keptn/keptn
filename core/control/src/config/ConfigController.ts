import 'reflect-metadata';
import * as express from 'express';
import { inject, injectable } from 'inversify';
import {
  controller,
  httpGet,
  httpPost,
  interfaces,
} from 'inversify-express-utils';
import {
  ApiOperationGet,
  ApiOperationPost,
  ApiPath,
  SwaggerDefinitionConstant,
} from 'swagger-express-ts';
import { ConfigRequestModel } from './ConfigRequestModel';
import { MessageService } from '../svc/MessageService';
import { WebSocketService } from '../svc/WebSocketService';

@ApiPath({
  name: 'Config',
  path: '/config',
  security: { apiKeyHeader: [] },
})
@controller('/config')
export class ConfigController implements interfaces.Controller {

  constructor(@inject('MessageService') private readonly messageService: MessageService) {}

  @ApiOperationPost({
    description: 'Set Github credentials for keptn',
    parameters: {
      body: {
        description: 'Github Credentials',
        model: 'ConfigRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Set Github credentials for keptn',
  })
  @httpPost('/')
  public async setGithubConfig(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log(`received config command...`);
    const channelInfo = await WebSocketService.getInstance().createChannel();
    if (request.body && request.body.data !== undefined) {
      request.body.data.channelInfo = channelInfo;
    }
    const result = await this.messageService.sendMessage(request.body);
    response.send({
      success: result,
      websocketChannel: channelInfo,
    });
  }
}
