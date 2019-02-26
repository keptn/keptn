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

@ApiPath({
  name: 'Auth',
  path: '/auth',
  security: { apiKeyHeader: [] },
})
@controller('/auth')
export class AuthController implements interfaces.Controller {
  constructor() { }

  @ApiOperationPost({
    description: 'Verify API Token',
    parameters: {
      body: {
        description: 'Signed payload',
        model: 'AuthRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Verifies if payload has been signed correctly with API Token',
  })
  @httpPost('/')
  public async verifyToken(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    // actual authentication is done via middleware
    response.send({ status: 'OK' });
  }
}
