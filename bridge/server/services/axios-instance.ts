import axios from 'axios';
import https from 'https';

const ax = axios.create({
  // accepts self-signed ssl certificate
  httpsAgent: new https.Agent({
    rejectUnauthorized: false
  })
});

export {ax as axios};
