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

    console.log(`Sending auth request: ${JSON.stringify(authRequest)}`);

    let authResult;
    try {
        authResult = await axios.post('http://authenticator.keptn.svc.cluster.local/auth', authRequest);
    } catch (e) {
        console.log(e);
    }

    console.log(authResult);
    
    if (authResult.data.authenticated) {
        next();
    } else {
        response.status(401);
    }
}

export = authenticator;