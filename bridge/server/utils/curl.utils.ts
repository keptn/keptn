import { WebhookConfig } from '../../shared/models/webhook-config';

function generateWebhookConfigCurl(webhookConfig: WebhookConfig): string {
  let params = '';
  for (const header of webhookConfig?.header || []) {
    params += `--header '${header.name}: ${header.value}' `;
  }
  params += `--request ${webhookConfig.method} `;
  if (webhookConfig.proxy) {
    params += `--proxy ${webhookConfig.proxy} `;
  }
  if (webhookConfig.payload) {
    let stringify = webhookConfig.payload;
    try {
      stringify = JSON.stringify(JSON.parse(webhookConfig.payload));
    } catch {
      stringify = stringify.replace(/\r\n|\n|\r/gm, '');
    }
    params += `--data '${stringify}' `;
  }
  return `curl ${params}${webhookConfig.url}`;
}

function parseCurl(curl: string): { [key: string]: string[] } {
  const startCommand = 'curl ';
  const result: { [key: string]: string[] } = {};
  if (curl.startsWith(startCommand)) {
    let i = startCommand.length;
    while (i < curl.length) {
      i = skipSpace(curl, i);
      let command = '_';
      if (curl[i] === '-') {
        const commandInfo = getNextCommand(curl, i);
        i = commandInfo.index + 1;
        command = commandInfo.data;
      }
      i = skipSpace(curl, i);
      if (i < curl.length) {
        const commandData = getNextCommandData(curl, i);
        i = commandData.index;
        const data = result[command];
        if (data) {
          data.push(commandData.data);
        } else {
          result[command] = [commandData.data];
        }
        ++i;
      }
    }
  }
  return result;
}

function skipSpace(curl: string, index: number): number {
  while (curl[index] === ' ') {
    ++index;
  }
  return index;
}

function getNextCommandData(curl: string, i: number): { data: string; index: number } {
  const startsWith = curl[i];
  let data = '';
  const startIndex = i;
  if (startsWith === "'" || startsWith === '"') {
    ++i;
    while (i < curl.length && (curl[i] !== startsWith || (curl[i] === startsWith && curl[i - 1] === '\\'))) {
      ++i;
    }
    data = curl.substring(startIndex + 1, i);
  } else {
    i = curl.indexOf(' ', startIndex);
    if (i === -1) {
      i = curl.length;
    }
    data = curl.substring(startIndex, i);
  }
  return {
    data,
    index: i,
  };
}

function getNextCommand(curl: string, i: number): { data: string; index: number } {
  let startCommandIndex = i + 1;
  if (curl[i + 1] === '-') {
    ++startCommandIndex;
  }
  i = curl.indexOf(' ', startCommandIndex);
  return {
    data: curl.substring(startCommandIndex, i),
    index: i === -1 ? curl.length : i,
  };
}

export { parseCurl, generateWebhookConfigCurl };
