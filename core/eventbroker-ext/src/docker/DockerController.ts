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

@ApiPath({
  name: 'Docker',
  path: '/docker',
  security: { apiKeyHeader: [] },
})
@controller('/docker')
export class DockerController implements interfaces.Controller {

  constructor(@inject('DockerService') private readonly dockerService: DockerService) {}

  @ApiOperationPost({
    description: 'Dispatch a new keptn event',
    parameters: {
      body: {
        description: 'Docker Webhook payload',
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
  public async handleDockerEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log(`received event: ${JSON.stringify(request.body)}`);

  }
}
