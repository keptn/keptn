import express = require('express');
import { ConfigHandler } from '../handler/config/configHandler';

const router = express.Router();

router.post('/', async (request: express.Request, response: express.Response) => {
  const configHandler = new ConfigHandler();
  await configHandler.init();

  try {
    await configHandler.updateKeptnConfig(request.body.data);
  } catch (e) {
    console.log(e);
  }


  response.send({ status: 'OK' });
});

// add more route handlers here
// e.g. router.post('/', (req,res,next)=> {/*...*/})
export = router;
