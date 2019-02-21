import express = require('express');

const router = express.Router();

router.post('/', (request: express.Request, response: express.Response) => {

  const result = {
    foo: 'bar',
  };

  response.send(result);
});

// add more route handlers here
// e.g. router.post('/', (req,res,next)=> {/*...*/})
export = router;
