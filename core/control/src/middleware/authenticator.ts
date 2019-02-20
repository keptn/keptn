import express = require("express");
import { AuthRequest } from '../types/authRequest';
import axios from 'axios';

let authenticator: express.RequestHandler = async (
    request: express.Request,
    response: express.Response,
    next: express.NextFunction
) => {
    
    // TODO: insert call to authenticator.keptn.svc.cluster.local here
    // get signature from header 
    let signature: string = request.headers['X-Keptn-Signature'] as string;
    let payload = JSON.stringify(request.body);

    let authRequest: AuthRequest = {
        signature,
        payload
    }

    let authResult = await axios.post('http://authenticator/auth', authRequest);

    console.log(authResult);
    
    if (authResult.data.authenticated) {
        next();
    } else {
        response.status(401);
    }
}

export = authenticator;