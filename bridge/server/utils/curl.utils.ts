function parseCurl(curl: string, ignoreStart = false): { [key: string]: string[] } {
  const startCommand = ignoreStart ? '' : 'curl ';
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

export { parseCurl };
