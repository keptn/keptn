import { join } from 'path';
import { existsSync, readFileSync } from 'fs';
import { BasicSecrets, MongoDBSecrets, OAuthSecrets } from '../interfaces/configuration';

export const mongodbUserFileName = 'mongodb-user';
export const mongodbPasswordFileName = 'mongodb-passwords';

function getOAuthSecrets(configFolder: string): OAuthSecrets {
  const oauthFolder = join(configFolder, 'oauth');
  const sessionSecretPath = join(oauthFolder, 'session_secret');
  const clientSecretPath = join(oauthFolder, 'client_secret');
  const databaseEncryptSecretPath = join(oauthFolder, 'database_encrypt_secret');

  return {
    sessionSecret: readSecret(sessionSecretPath),
    databaseEncryptSecret: readSecret(databaseEncryptSecretPath),
    clientSecret: readSecret(clientSecretPath),
  };
}

function getOAuthMongoExternalConnectionString(configFolder: string): string {
  const mongodbFolder = join(configFolder, 'oauth_mongodb_connection_string');
  const mongoSecretPath = join(mongodbFolder, 'external_connection_string');
  return readSecret(mongoSecretPath);
}

function getMongoDbSecrets(configFolder: string): MongoDBSecrets {
  const mongodbFolder = getMongodbFolder(configFolder);
  const userPath = join(mongodbFolder, mongodbUserFileName);
  const passwordPath = join(mongodbFolder, mongodbPasswordFileName);

  return {
    user: readSecret(userPath),
    password: readSecret(passwordPath),
  };
}

function getMongodbFolder(configFolder: string): string {
  return join(configFolder, 'oauth_mongodb');
}

function getBasicSecrets(configFolder: string): BasicSecrets {
  const basicCredentialFolder = join(configFolder, 'basic');
  const apiFolder = join(configFolder, 'api-token');
  const apiTokenPath = join(apiFolder, 'keptn-api-token');
  const basicUser = join(basicCredentialFolder, 'BASIC_AUTH_USERNAME');
  const basicPassword = join(basicCredentialFolder, 'BASIC_AUTH_PASSWORD');
  return {
    apiToken: readSecret(apiTokenPath),
    user: readSecret(basicUser),
    password: readSecret(basicPassword),
  };
}

function readSecret(path: string): string {
  if (existsSync(path)) {
    return readFileSync(path, { encoding: 'utf8', flag: 'r' });
  }
  return '';
}

export { getOAuthSecrets, getOAuthMongoExternalConnectionString, getMongoDbSecrets, getBasicSecrets, getMongodbFolder };
