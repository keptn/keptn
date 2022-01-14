import { Express } from 'express';
import { AxiosInstance } from 'axios';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace NodeJS {
    interface Global {
      app: Express;
      baseUrl: string;
      axiosInstance: AxiosInstance;
    }
  }
}
