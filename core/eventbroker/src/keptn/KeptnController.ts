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
import { KeptnRequestModel } from './KeptnRequestModel';
import { MessageService } from '../svc/MessageService';

@ApiPath({
  name: 'Keptn',
  path: '/keptn',
  security: { apiKeyHeader: [] },
})
@controller('/keptn')
export class KeptnController implements interfaces.Controller {

  constructor(@inject('MessageService') private readonly messageService: MessageService) {}

  @ApiOperationPost({
    description: 'Dispatch a new keptn event',
    parameters: {
      body: {
        description: 'Keptn CloudEvent',
        model: 'KeptnRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Dispatch a new keptn event',
  })
  @httpPost('/')
  public async dispatchEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log(`received event: ${JSON.stringify(request.body)}`);

    const result = await this.messageService.sendMessage(request.body);
    response.send({ success: result });
  }
}
