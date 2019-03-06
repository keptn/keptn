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
import { DockerService } from './DockerService';
import { DockerRequestModel } from './DockerRequestModel';

@ApiPath({
  name: 'Docker',
  path: '/docker',
  security: { apiKeyHeader: [] },
})
@controller('/docker')
export class DockerController implements interfaces.Controller {

  constructor(@inject('DockerService') private readonly dockerService: DockerService) { }

  @ApiOperationPost({
    description: 'Handle an incoming docker event',
    parameters: {
      body: {
        description: 'Docker Webhook payload',
        model: 'DockerRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Handle an incoming docker event',
  })
  @httpPost('/')
  public async handleDockerEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log(`received event: ${JSON.stringify(request.body)}`);
    const messageSent =
      await this.dockerService.handleDockerRequest(request.body as DockerRequestModel);
    response.status(200);
    response.send({
      messageSent,
    });
  }
}
