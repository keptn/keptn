import express = require('express');
import { CredentialsService } from '../service/CredentialsService';

const router = express.Router();

router.post('/', async (request: express.Request, response: express.Response) => {
  const credentialsService = CredentialsService.getInstance();
  try {
    await credentialsService.updateGithubConfig(request.body.data);
    const secret = await credentialsService.getGithubCredentials();
  } catch (e) {
    console.log(e);
  }

  response.send({ status: 'OK' });
});

// add more route handlers here
// e.g. router.post('/', (req,res,next)=> {/*...*/})
export = router;
