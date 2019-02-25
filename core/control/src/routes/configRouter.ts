import express = require('express');
import { CredentialsService } from '../service/CredentialsService';
import {
  ApiOperationGet,
  ApiOperationPost,
  ApiPath,
  SwaggerDefinitionConstant,
} from 'swagger-express-ts';

const router = express.Router();

router.post('/', async (request: express.Request, response: express.Response) => {
  const credentialsService = CredentialsService.getInstance();
  
  try {
    await credentialsService.updateGithubConfig(request.body.data);
  } catch (e) {
    console.log(e);
  }

  response.send({ status: 'OK' });
});

// add more route handlers here
// e.g. router.post('/', (req,res,next)=> {/*...*/})
export = router;
