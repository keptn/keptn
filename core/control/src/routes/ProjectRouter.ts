import express = require('express');

const router = express.Router();

router.post('/', async (request: express.Request, response: express.Response) => {

  // TODO: convert payload into a CloudEvent containing the following data block:

  /*
    data : {
      project: 'sockshop',
      stages: [
        {
          name: 'dev',
          deployment_strategy: 'direct',
        },
        {
          name: 'staging',
          deployment_strategy: 'blue_green_service',
        },
        {
          name: 'production',
          deployment_strategy: 'blue_green_service',
        },
      ],
    }
  */

  // Post this CloudEvent into the Channel.

  const result = {
    result: 'success',
  };

  response.send(result);
});

router.get('/', (request: express.Request, response: express.Response) => {

  const result = {
    result: 'success',
  };

  response.send(result);
});

router.delete('/', async (request: express.Request, response: express.Response) => {

  const result = {
    result: 'success',
  };

  response.send(result);
});

export = router;
