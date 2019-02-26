import * as express from 'express';
import { inject, injectable } from 'inversify';
import {
  controller,
  httpGet,
  httpPost,
  interfaces,
  httpDelete,
} from 'inversify-express-utils';
import 'reflect-metadata';
import {
  ApiOperationGet,
  ApiOperationPost,
  ApiPath,
  SwaggerDefinitionConstant,
  ApiOperationDelete,
} from 'swagger-express-ts';

@ApiPath({
  name: 'Project',
  path: '/project',
  security: { apiKeyHeader: [] },
})
@controller('/project')
export class ProjectController implements interfaces.Controller {
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
  public setGithubConfig(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): void {
    const result = {
      result: 'success',
    };

    response.send(result);
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
