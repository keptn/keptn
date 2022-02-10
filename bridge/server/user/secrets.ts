import { dirname, join } from 'path';
import { existsSync, readFileSync } from 'fs';
import { fileURLToPath } from 'url';

const configFolder = process.env.CONFIG_DIR || join(dirname(fileURLToPath(import.meta.url)), '../../../../config');
const oauthFolder = join(configFolder, 'oauth');
const mongodbFolder = join(configFolder, 'oauth_mongodb');
const sessionSecretPath = join(oauthFolder, 'session_secret');
const databaseEncryptSecretPath = join(oauthFolder, 'database_encrypt_secret');
const mongoSecretPath = join(mongodbFolder, 'external_connection_string');

function getOAuthSecrets(): { sessionSecret: string; databaseEncryptSecret: string } {
  const secrets = {
    sessionSecret: '',
    databaseEncryptSecret: '',
  };

  if (existsSync(sessionSecretPath)) {
    secrets.sessionSecret = readFileSync(sessionSecretPath, { encoding: 'utf8', flag: 'r' });
  }
  if (existsSync(databaseEncryptSecretPath)) {
    secrets.databaseEncryptSecret = readFileSync(databaseEncryptSecretPath, { encoding: 'utf8', flag: 'r' });
  }
  return secrets;
}

function getOAuthMongoExternalConnectionString(): string {
  let externalConnection = '';
  if (existsSync(mongoSecretPath)) {
    externalConnection = readFileSync(mongoSecretPath, { encoding: 'utf8', flag: 'r' });
  }
  return externalConnection;
}

export { getOAuthSecrets, getOAuthMongoExternalConnectionString };
