import 'reflect-metadata';
import { expect } from 'chai';
import 'mocha';
import { AuthService } from './AuthService';
import { AuthRequestModel } from './AuthRequestModel';
import { BearerAuthRequestModel } from './BearerAuthRequestModel';

describe('AuthService', () => {
  let authService: AuthService;

  beforeEach(() => {
    authService = new AuthService();
  });

  it('Should return true for correct signatures', async () => {
    const authRequest: AuthRequestModel = {
      payload: '2344',
      signature: 'sha1=8a12fb3402b3d203aee37701f39e22a104b9a2a0',
    };

    const result = authService.verify(authRequest);

    expect(result).is.true;
  });

  it('Should return false for incorrect signatures', async () => {
    const authRequest: AuthRequestModel = {
      payload: '2344',
      signature: 'sha1=invalid',
    };

    const result = authService.verify(authRequest);

    expect(result).is.false;
  });

  it('Should return true for valid bearer tokens', async () => {
    const authRequest: BearerAuthRequestModel = {
      token: '',
    };

    const result = authService.verifyBearerToken(authRequest);

    expect(result).is.true;
  });

  it('Should return false for invalid bearer tokens', async () => {
    const authRequest: BearerAuthRequestModel = {
      token: 'invalid',
    };

    const result = authService.verifyBearerToken(authRequest);

    expect(result).is.false;
  });
});
