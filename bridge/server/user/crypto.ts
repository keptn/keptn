import { createCipheriv, createDecipheriv, randomBytes } from 'crypto';

export interface IEncrypted {
  iv: string;
  content: string;
  auth: string;
}

export class Crypto {
  private readonly iv: Buffer;
  private readonly algorithm = 'aes-256-gcm';

  constructor(private readonly secretKey: string) {
    this.iv = randomBytes(16);
  }

  public encrypt(text: string): string {
    const cipher = createCipheriv(this.algorithm, this.secretKey, this.iv);
    const enc1 = cipher.update(text, 'utf8');
    const enc2 = cipher.final();
    const data: IEncrypted = {
      iv: this.iv.toString('hex'),
      content: Buffer.concat([enc1, enc2]).toString('hex'),
      auth: cipher.getAuthTag().toString('hex'),
    };
    return JSON.stringify(data);
  }

  public decrypt(text: string): string {
    const parsed: IEncrypted = JSON.parse(text);
    const decipher = createDecipheriv(this.algorithm, this.secretKey, Buffer.from(parsed.iv, 'hex'));
    decipher.setAuthTag(Buffer.from(parsed.auth, 'hex'));
    return decipher.update(parsed.content, 'hex', 'utf8') + decipher.final('utf8');
  }
}
