import { init } from '../app';
import Axios from 'axios';
import https from 'https';

const setup = async (): Promise<void> => {
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

  global.app = await init();
};

export default setup();
