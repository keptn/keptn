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
import { GitHubService } from './GitHubService';

@ApiPath({
  name: 'GitHub',
  path: '/github',
  security: { apiKeyHeader: [] },
})
@controller('/github')
export class GitHubController implements interfaces.Controller {

  constructor(@inject('GitHubService') private readonly gitHubService: GitHubService) {}

  @ApiOperationPost({
    description: 'Dispatch a new keptn event',
    parameters: {
      body: {
        description: 'GitHub Webhook payload',
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
  public async handleGitHubEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log(`received event: ${JSON.stringify(request.body)}`);

  }
}
