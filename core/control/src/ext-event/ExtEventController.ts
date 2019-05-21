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
import { ExtEventRequestModel } from './ExtEventRequestModel'
import { MessageService } from '../svc/MessageService';
import { WebSocketService } from '../svc/WebSocketService';

const uuidv4 = require('uuid/v4');

@ApiPath({
  name: 'Event',
  path: '/event',
  security: { apiKeyHeader: [] },
})
@controller('/event')
export class ExtEventController implements interfaces.Controller {

  constructor(@inject('MessageService') private readonly messageService: MessageService) {}

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
    const keptnContext = uuidv4();
    
    console.log(JSON.stringify({
      keptnContext,
      keptnService: 'control',
      logLevel: 'INFO',
      message: `received external event: ${JSON.stringify(request.body)}`,
    }));
    const result = request.body;

    const channelInfo = await WebSocketService.getInstance().createChannel(keptnContext);
    if (result && result.data !== undefined) {
      result.data.channelInfo = channelInfo;
      result.shkeptncontext = keptnContext;
    }
    this.messageService.sendMessage(result);
    response.send(result);
  }
}
