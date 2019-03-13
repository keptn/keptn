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
import { AuthRequestModel } from './AuthRequestModel';
import { BearerAuthRequestModel } from './BearerAuthRequestModel';
import { AuthService } from './AuthService';

@ApiPath({
  name: 'Auth',
  path: '/auth',
  security: { apiKeyHeader: [] },
})
@controller('/auth')
export class AuthController implements interfaces.Controller {

  constructor(@inject('AuthService') private readonly authService: AuthService) {}

  @ApiOperationPost({
    description: 'Verifiy authentication request (Sha1 signature)',
    parameters: {
      body: {
        description: 'AuthRequest',
        model: 'AuthRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Validate auth request',
  })
  @httpPost('/')
  public async authenticate(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log('Starting authentication');
    const authRequest: AuthRequestModel = request.body;

    if (!AuthRequestModel.isAuthRequestModel(authRequest)) {
      console.log('Not a valid request');
      response.status(422);
      response.send();
      return;
    }
    console.log(`Received auth request: ${JSON.stringify(authRequest)}`);

    const authResult = {
      authenticated: this.authService.verify(authRequest),
    };

    console.log(`Response: ${JSON.stringify(authResult)}`);

    response.send(authResult);
  }

  @ApiOperationPost({
    description: 'Verifiy authentication request (Bearer token)',
    parameters: {
      body: {
        description: 'BearerAuthRequest',
        model: 'BearerAuthRequestModel',
        required: true,
      },
    },
    responses: {
      200: {
      },
      400: { description: 'Parameters fail' },
    },
    summary: 'Validate auth request',
  })
  @httpPost('/token')
  public async authenticateToken(
    request: express.Request,
    response: express.Response,
    next: express.NextFunction,
  ): Promise<void> {
    console.log('Starting authentication');
    const authRequest: BearerAuthRequestModel = request.body;

    if (!BearerAuthRequestModel.isBearerAuthRequestModel(authRequest)) {
      console.log('Not a valid request');
      response.status(422);
      response.send();
      return;
    }
    console.log(`Received auth request: ${JSON.stringify(authRequest)}`);

    const authResult = {
      authenticated: this.authService.verifyBearerToken(authRequest),
    };

    console.log(`Response: ${JSON.stringify(authResult)}`);

    response.send(authResult);
  }
}
