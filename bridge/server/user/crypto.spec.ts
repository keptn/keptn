import { Crypto, IEncrypted } from './crypto';

describe('Test crypto.ts', () => {
  it('should encrypt and decrypt data', () => {
    const crypto = new Crypto('secret__'.repeat(4)); // length of 32

    for (const data of ['myText', '', 't']) {
      const encrypted = crypto.encrypt(data);
      expect(encrypted).not.toBe(data);
      const decrypted = crypto.decrypt(encrypted);
      expect(decrypted).toBe(data);
    }
  });

  it('should throw exception if secret is invalid', () => {
    const text = 'myText';
    const crypto = new Crypto('secret__'.repeat(4)); // length of 32
    const encrypted = crypto.encrypt(text);
    const crypto2 = new Crypto('secret_1'.repeat(4)); // length of 32
    expect(() => crypto2.decrypt(encrypted)).toThrowError();
  });

  it('should throw exception if initialization vector is invalid', () => {
    const text = 'myText';
    const crypto = new Crypto('secret__'.repeat(4)); // length of 32
    const encrypted = crypto.encrypt(text);
    const data: IEncrypted = JSON.parse(encrypted);
    data.iv = Buffer.from('iv'.repeat(8), 'utf-8').toString('hex');
    expect(() => crypto.decrypt(JSON.stringify(data))).toThrowError();
  });

  it('should throw exception if auth is invalid', () => {
    const text = 'myText';
    const crypto = new Crypto('secret__'.repeat(4)); // length of 32
    const encrypted = crypto.encrypt(text);
    const data: IEncrypted = JSON.parse(encrypted);
    data.auth = Buffer.from('iv'.repeat(8), 'utf-8').toString('hex');
    expect(() => crypto.decrypt(JSON.stringify(data))).toThrowError();
  });

  it('should throw exception if content is invalid', () => {
    const text = 'myText';
    const crypto = new Crypto('secret__'.repeat(4)); // length of 32
    const encrypted = crypto.encrypt(text);
    const data: IEncrypted = JSON.parse(encrypted);
    data.content = Buffer.from('something', 'utf-8').toString('hex');
    expect(() => crypto.decrypt(JSON.stringify(data))).toThrowError();
  });
});
