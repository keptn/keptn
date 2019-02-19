import express = require("express");

let router = express.Router();

router.post('/', (request: express.Request, response: express.Response) => {

    let result = {
        "foo": "bar"
    };

    response.send(result);
});
// add more route handlers here
// e.g. router.post('/', (req,res,next)=> {/*...*/})
export = router;