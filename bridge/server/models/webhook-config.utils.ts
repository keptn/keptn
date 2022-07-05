import {
  IWebhookConfigYamlResult,
  IWebhookConfigYamlResultV1Alpha1,
  IWebhookConfigYamlResultV1Beta1,
  IWebhookRequestV1Beta1,
  IWebhookSecret,
  IWebhookV1Alpha1,
  IWebhookV1Beta1,
  WebhookApiVersions,
} from '../interfaces/webhook-config-yaml-result';
import { parseCurl } from '../utils/curl.utils';
import { IWebhookConfigClient, WebhookConfigMethod } from '../../shared/interfaces/webhook-config';
import { WebhookConfigYaml } from './webhook-config-yaml';
import { IClientSecret } from '../../shared/interfaces/secret';

interface FlatSecret {
  path: string;
  name: string;
  key: string;
  parsedPath: string;
}

const mapHeaders = (header: string): { key: string; value: string } => {
  const headerInfo = header.split(':');
  return {
    key: headerInfo[0]?.trim(),
    value: headerInfo[1]?.trim(),
  };
};

const mapRequests = (curl: string): IWebhookRequestV1Beta1 => {
  const result = parseCurl(curl);
  const proxy = result.proxy?.[0];

  return {
    url: result._?.join(' ') ?? '',
    headers: result.header?.map(mapHeaders),
    method: (result.request?.[0] ?? 'GET') as WebhookConfigMethod,
    payload: formatJSON(result.data?.[0] ?? ''),
    ...(proxy && { options: `--proxy ${proxy}` }),
  };
};

const mapWebhooks = (webhook: IWebhookV1Alpha1): IWebhookV1Beta1 => {
  return {
    requests: webhook.requests.map(mapRequests),
    envFrom: webhook.envFrom,
    sendFinished: webhook.sendFinished,
    sendStarted: webhook.sendStarted,
    type: webhook.type,
    subscriptionID: webhook.subscriptionID,
  };
};

export function migrateWebhook(config: IWebhookConfigYamlResult): IWebhookConfigYamlResultV1Beta1 {
  if (config.apiVersion === WebhookApiVersions.V1ALPHA1) {
    return migrateV1Alpha1ToV1Beta1(config);
  }
  return config;
}

function migrateV1Alpha1ToV1Beta1(config: IWebhookConfigYamlResultV1Alpha1): IWebhookConfigYamlResultV1Beta1 {
  return {
    apiVersion: WebhookApiVersions.V1BETA1,
    kind: 'WebhookConfig',
    metadata: {
      name: 'webhook-configuration',
    },
    spec: {
      webhooks: config.spec.webhooks.map(mapWebhooks),
    },
  };
}

export function formatJSON(data: string): string {
  try {
    data = JSON.stringify(JSON.parse(data), null, 2);
  } catch {}
  return data;
}

export function stringifyPayload(payload: string): string {
  try {
    return JSON.stringify(JSON.parse(payload));
  } catch {
    return payload.replace(/\r\n|\n|\r/gm, '');
  }
}

export function parseToClientWebhookRequest(
  config: WebhookConfigYaml,
  subscriptionId: string
): { webhookConfig: IWebhookConfigClient; secrets?: IWebhookSecret[] } | undefined {
  const webhook = config.getWebhook(subscriptionId);
  const request = webhook?.requests[0];
  if (!webhook || !request) {
    return undefined;
  }

  return {
    webhookConfig: {
      header: request.headers ?? [],
      url: request.url,
      payload: formatJSON(request.payload ?? ''),
      method: request.method,
      sendFinished: webhook.sendFinished ?? false,
      sendStarted: webhook.sendStarted ?? true,
      type: webhook.type,
      proxy: request.options ? parseCurl(request.options, true)?.proxy?.[0] ?? '' : '',
    },
    secrets: webhook.envFrom,
  };
}

export function mapYamlSecretsToBridgeSecrets(webhookConfig: IWebhookConfigClient, secrets: IWebhookSecret[]): void {
  for (const webhookSecret of secrets) {
    const bridgeSecret = `{{.secret.${webhookSecret.secretRef.name}.${webhookSecret.secretRef.key}}}`;
    const regex = new RegExp(`{{.env.${webhookSecret.name}\}\}`, 'g');

    webhookConfig.url = webhookConfig.url.replace(regex, bridgeSecret);
    webhookConfig.payload = webhookConfig.payload.replace(regex, bridgeSecret);
    for (const header of webhookConfig.header) {
      header.value = header.value.replace(regex, bridgeSecret);
    }
  }
}

export function mapBridgeSecretsToYamlSecrets(
  config: IWebhookConfigClient,
  secrets: IClientSecret[]
): IWebhookSecret[] {
  const flatSecret = getSecretPathFlat(secrets);
  const webhookSecrets: IWebhookSecret[] = [];

  config.url = addWebhookSecretsFromString(config.url, flatSecret, webhookSecrets);
  config.payload = addWebhookSecretsFromString(config.payload, flatSecret, webhookSecrets);

  for (const head of config.header) {
    head.value = addWebhookSecretsFromString(head.value, flatSecret, webhookSecrets);
  }
  return webhookSecrets;
}

function getSecretPathFlat(secrets: IClientSecret[]): FlatSecret[] {
  return secrets
    .filter((secret): secret is IClientSecret & { keys: string[] } => !!secret.keys)
    .reduce(
      (flatSecrets: FlatSecret[], secret) => [
        ...flatSecrets,
        ...secret.keys.map((key) => mapSecretToFlatSecret(secret.name, key)),
      ],
      [] as FlatSecret[]
    );
}

function mapSecretToFlatSecret(name: string, key: string): FlatSecret {
  const sanitizedName = name.replace(/[^a-zA-Z0-9]/g, '');
  const sanitizedKey = key.replace(/[^a-zA-Z0-9]/g, '');
  return {
    path: `${name}.${key}`,
    name,
    key,
    parsedPath: `secret_${sanitizedName}_${sanitizedKey}`,
  };
}

function addWebhookSecretsFromString(
  parseString: string,
  allSecretPaths: FlatSecret[],
  existingSecrets: IWebhookSecret[]
): string {
  const foundSecrets = allSecretPaths.filter((scrt) => parseString.includes(scrt.path));
  let replacedString = parseString;
  for (const found of foundSecrets) {
    const idx = existingSecrets.findIndex((secret) => secret.name === found.parsedPath);
    if (idx === -1) {
      const secret: IWebhookSecret = {
        name: found.parsedPath,
        secretRef: {
          name: found.name,
          key: found.key,
        },
      };
      existingSecrets.push(secret);
    }

    replacedString = replacedString.replace(
      new RegExp(`{{.secret.${found.path}}}`, 'g'),
      `{{.env.${found.parsedPath}}}`
    );
  }

  return replacedString;
}
