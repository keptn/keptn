import { init } from '../app';
import Axios from 'axios';
import https from 'https';
import { Express } from 'express';

const setupServer = async (token = process.env.API_TOKEN): Promise<Express> => {
  global.baseUrl = 'http://localhost/api/';

  global.axiosInstance = Axios.create({
    // accepts self-signed ssl certificate
    httpsAgent: new https.Agent({
      rejectUnauthorized: false,
    }),
    headers: {
      ...token && { 'x-token': token },
      'Content-Type': 'application/json',
    },
  });

  return init();
};

export { setupServer };
