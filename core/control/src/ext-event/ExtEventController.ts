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
import { ExtEventService } from './ExtEventService';
import { ExtEventRequestModel } from './ExtEventRequestModel'

@ApiPath({
  name: 'Event',
  path: '/event',
  security: { apiKeyHeader: [] },
})
@controller('/event')
export class ExtEventController implements interfaces.Controller {

  constructor(@inject('ExtEventService') private readonly extEventService: ExtEventService) {}

  @ApiOperationPost({
    description: 'Handle incoming external event',
    parameters: {
      body: {
        description: 'Keptn CloudEvent',
        model: 'ExtEventRequestModel',
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
  public async handleExtEvent(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    let keptnContext = '';
    if (request.body !== undefined &&
      request.body.shkeptncontext !== undefined) {
      keptnContext = request.body.shkeptncontext;
    }
    console.log(JSON.stringify({
      keptnContext,
      keptnService: 'eventbroker-ext',
      logLevel: 'INFO',
      message: `received event: ${JSON.stringify(request.body)}`,
    }));

    this.extEventService.handleExtEvent(request.body);
  }
}
