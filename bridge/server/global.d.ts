import { AxiosInstance } from 'axios';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace NodeJS {
    interface Global {
      axiosInstance?: AxiosInstance;
      issuer?: unknown;
    }
  }
}
