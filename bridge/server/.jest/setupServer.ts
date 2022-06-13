import { init } from '../app';
import Axios from 'axios';
import https from 'https';
import { Express } from 'express';
import { BridgeConfiguration, BridgeOption, getConfiguration } from '../utils/configuration';

const baseOptions: BridgeOption = {
  api: {
    token: 'apiToken',
    url: 'http://localhost/api/',
  },
  oauth: {
    baseURL: 'http://baseoauthurl',
    clientID: 'myclientid',
    discoveryURL: 'http://discoveryurl',
    enabled: false,
  },
  mongo: {
    host: 'mongo://localhost',
    password: 'pwd',
    user: 'usr',
  },
  version: 'develop',
};

const baseConfig = getConfiguration(baseOptions);

const setupServer = async (config: BridgeConfiguration = Object.create({})): Promise<Express> => {
  global.baseUrl = 'http://localhost/api/';

  // create a new fresh configuration
  if (config && Object.keys(config).length === 0) {
    config = getConfiguration(baseOptions);
  }

  global.axiosInstance = Axios.create({
    // accepts self-signed ssl certificate
    httpsAgent: new https.Agent({
      rejectUnauthorized: false,
    }),
    headers: {
      ...(config.api.token && { 'x-token': config.api.token }),
      'Content-Type': 'application/json',
    },
  });

  return init(config);
};

export { setupServer, baseOptions, baseConfig };
