import express = require("express");

const router = express.Router();

router.post('/', (request: express.Request, response: express.Response) => {

  // TODO: Onboarding - convert payload into a CloudEvent containing the following data block:

  /*
    data : {
      project: 'sockshop',
      file : // deployment and service definition
    }
  */

  // Post this CloudEvent into the Channel.

  const result = {
    result: 'success',
  };

  response.send(result);
});

// add more route handlers here
// e.g. router.post('/', (req,res,next)=> {/*...*/})
export = router;
