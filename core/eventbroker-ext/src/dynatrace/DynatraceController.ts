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
import { DynatraceService } from './DynatraceService';

@ApiPath({
  name: 'Dynatrace',
  path: '/dynatrace',
  security: { apiKeyHeader: [] },
})
@controller('/dynatrace')
export class DynatraceController implements interfaces.Controller {

  constructor(@inject('DynatraceService') private readonly dynatraceService: DynatraceService) {}

  @ApiOperationPost({
    description: 'Handle incoming Dynatrace problem notifications',
    parameters: {
      body: {
        description: 'Dynatrace problem payload',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: '',
  })
  @httpPost('/')
  public async handleDynatraceEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log(`received event: ${JSON.stringify(request.body)}`);
    this.dynatraceService.handleDynatraceEvent(request.body);
  }
}
