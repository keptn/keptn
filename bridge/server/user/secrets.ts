import { join } from 'path';
import { existsSync, readFileSync } from 'fs';
import { OAuthSecrets } from '../interfaces/configuration';

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
  const mongodbFolder = join(configFolder, 'oauth_mongodb');
  const mongoSecretPath = join(mongodbFolder, 'external_connection_string');
  return readSecret(mongoSecretPath);
}

function readSecret(path: string): string {
  if (existsSync(path)) {
    return readFileSync(path, { encoding: 'utf8', flag: 'r' });
  }
  return '';
}

export { getOAuthSecrets, getOAuthMongoExternalConnectionString };
