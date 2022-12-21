// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';
import { join } from 'path';
import { BasicSecrets, MongoDBSecrets } from '../interfaces/configuration';

const readFileSyncSpy = jest.fn();
const existsSyncSpy = jest.fn();

jest.unstable_mockModule('fs', () => {
  return {
    readFileSync: readFileSyncSpy,
    existsSync: existsSyncSpy,
  };
});

const { getOAuthMongoExternalConnectionString, getOAuthSecrets, getBasicSecrets, getMongoDbSecrets } = await import(
  './secrets'
);

describe('Test fetching secrets from the file system', () => {
  const options = { encoding: 'utf8', flag: 'r' };

  afterEach(() => {
    existsSyncSpy.mockReset();
    readFileSyncSpy.mockReset();
  });

  it('should read OAuth secrets', () => {
    // given
    existsSyncSpy.mockReturnValue(true);
    readFileSyncSpy.mockReturnValueOnce('secret1');
    readFileSyncSpy.mockReturnValueOnce('secret2');
    readFileSyncSpy.mockReturnValueOnce('secret3');

    // when
    const secrets = getOAuthSecrets('config');

    // then
    expect(secrets).toEqual({
      sessionSecret: 'secret1',
      databaseEncryptSecret: 'secret2',
      clientSecret: 'secret3',
    });
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'oauth', 'session_secret'), options);
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'oauth', 'client_secret'), options);
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'oauth', 'database_encrypt_secret'), options);
  });

  it('should read mongo secret', () => {
    existsSyncSpy.mockReturnValue(true);
    readFileSyncSpy.mockReturnValueOnce('secretMongo');

    const secret = getOAuthMongoExternalConnectionString('config');

    expect(secret).toBe('secretMongo');
    expect(readFileSyncSpy).toHaveBeenCalledWith(
      join('config', 'oauth_mongodb_connection_string', 'external_connection_string'),
      options
    );
  });

  it('should return empty string if the directory can not be found', () => {
    // given
    existsSyncSpy.mockReturnValue(false);

    // when
    const mongoSecret = getOAuthMongoExternalConnectionString('config');
    const oAuthSecrets = getOAuthSecrets('config');

    // then
    expect(mongoSecret).toBe('');
    expect(oAuthSecrets).toEqual({
      sessionSecret: '',
      databaseEncryptSecret: '',
      clientSecret: '',
    });
    expect(readFileSyncSpy).not.toHaveBeenCalled();
  });

  it('should return mongodb secrets', () => {
    existsSyncSpy.mockReturnValue(true);
    readFileSyncSpy.mockReturnValueOnce('secret1');
    readFileSyncSpy.mockReturnValueOnce('secret2');

    const mongoSecrets = getMongoDbSecrets('config');

    expect(mongoSecrets).toEqual<MongoDBSecrets>({
      user: 'secret1',
      password: 'secret2',
    });
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'oauth_mongodb', 'mongodb-user'), options);
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'oauth_mongodb', 'mongodb-passwords'), options);
  });

  it('should return basic secrets', () => {
    existsSyncSpy.mockReturnValue(true);
    readFileSyncSpy.mockReturnValueOnce('secret1');
    readFileSyncSpy.mockReturnValueOnce('secret2');
    readFileSyncSpy.mockReturnValueOnce('secret3');

    const mongoSecrets = getBasicSecrets('config');

    expect(mongoSecrets).toEqual<BasicSecrets>({
      apiToken: 'secret1',
      user: 'secret2',
      password: 'secret3',
    });
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'api-token', 'keptn-api-token'), options);
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'basic', 'BASIC_AUTH_USERNAME'), options);
    expect(readFileSyncSpy).toHaveBeenCalledWith(join('config', 'basic', 'BASIC_AUTH_PASSWORD'), options);
  });
});
