import { init } from '../app';
import Axios from 'axios';
import https from 'https';
import { Express } from 'express';

const setupServer = async (): Promise<Express> => {
  global.baseUrl = 'http://localhost/api/';

  global.axiosInstance = Axios.create({
    // accepts self-signed ssl certificate
    httpsAgent: new https.Agent({
      rejectUnauthorized: false,
    }),
    headers: {
      'x-token': process.env.API_TOKEN ?? '',
      'Content-Type': 'application/json',
    },
  });

  return init();
};

export { setupServer };
