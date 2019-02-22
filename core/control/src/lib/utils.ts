const fs = require('fs');

class Utils {
  readFileContent(filePath: string) {
    return String(fs.readFileSync(filePath));
  }
}

export { Utils };