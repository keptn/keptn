import { Express } from 'express';
import * as http from 'http';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace NodeJS {
    interface Global {
      app: Express;
      baseUrl: string;
      server?: http.Server;
    }
  }
}
