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
import { Service } from './Service';
import { RequestModel } from './RequestModel';

@ApiPath({
  name: '',
  path: '/',
  security: { apiKeyHeader: [] },
})
@controller('/')
export class Controller implements interfaces.Controller {

  constructor(@inject('Service') private readonly service: Service) { }

  @ApiOperationPost({
    description: 'Handle an incoming event',
    parameters: {
      body: {
        description: 'payload',
        model: 'RequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Handle an incoming  event',
  })
  @httpPost('/')
  public async handleEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    await this.service.handleRequest(request.body as RequestModel);
    response.status(200);
    response.send({});
  }
}
