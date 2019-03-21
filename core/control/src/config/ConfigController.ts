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
import { CredentialsService } from '../svc/CredentialsService';
import { MessageService } from '../svc/MessageService';

const uuidv4 = require('uuid/v4');

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
    const keptnContext = uuidv4();
    const result = {
      keptnContext,
      success: true,
    };
    if (request.body !== undefined && request.body.data !== undefined) {
      request.body.data.keptnContext = keptnContext;
    }
    result.success = await this.messageService.sendMessage(request.body);
    response.send(result);
  }
}
