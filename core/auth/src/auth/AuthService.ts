import { AuthRequestModel } from './AuthRequestModel';
const crypto = require('crypto');
const bufferEq = require('buffer-equal-constant-time');
import { injectable } from 'inversify';

@injectable()
export class AuthService {
  private sign(data: string): string {
    const signature =
      `sha1=${crypto.createHmac('sha1', process.env.SECRET_TOKEN || '')
        .update(data).digest('hex')}`;

    console.log(`Calculated signature: ${signature}`);
    return signature;
  }

  verify(authRequest: AuthRequestModel): boolean {
    return bufferEq(
        Buffer.from(authRequest.signature),
        Buffer.from(this.sign(authRequest.payload)));
  }
}
