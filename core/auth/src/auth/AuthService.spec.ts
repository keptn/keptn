import 'reflect-metadata';
import { expect } from 'chai';
import 'mocha';
import { AuthService } from './AuthService';
import { AuthRequestModel } from './AuthRequestModel';

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
});
