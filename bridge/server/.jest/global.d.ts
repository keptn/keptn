import { AxiosInstance } from 'axios';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace NodeJS {
    interface Global {
      baseUrl: string;
      axiosInstance: AxiosInstance;
      issuer?: unknown;
    }
  }
}
