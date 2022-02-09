import { dirname, join } from 'path';
import { existsSync, readFileSync } from 'fs';
import { fileURLToPath } from 'url';

const configFolder = join(
  dirname(fileURLToPath(import.meta.url)),
  process.env.NODE_ENV === 'development' ? '../../../../' : '../../../../../../..',
  'config/oauth'
);
const sessionSecretPath = join(configFolder, 'session_secret');
const databaseEncryptSecretPath = join(configFolder, 'database_encrypt_secret');

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

export { getOAuthSecrets };
