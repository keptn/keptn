const express = require('express');
const router = express.Router();
const axios = require('axios');
const sessionAuthentication = require('./session').setAuthenticatedPrincipal;
const removeSession = require('./session').removeSession;

const AUTHORIZATION = 'authorization';
const AUTH_URL = 'authorization_url';
const TOKEN_DECISION = 'token_decision';

module.exports = (async () => {
  console.log('Enabling OAuth for bridge.');

  const discoveryEndpoint = process.env.OAUTH_DISCOVERY;

  if (!discoveryEndpoint) {
    throw 'OAUTH_DISCOVERY must be defined when oauth based login (OAUTH_ENABLED) is activated.' +
    ' Please check your environment variables.';
  }

  const discoveryResp = await axios({
    method: 'get',
    url: discoveryEndpoint,
  })

  if (discoveryResp.status !== 200) {
    throw `Invalid oauth service discovery response. Received status : ${discoveryResp.statusText}.`;
  }

  if (!discoveryResp.data.hasOwnProperty(AUTHORIZATION)) {
    throw 'OAuth service discovery must contain the authorization endpoint.';
  }

  if (!discoveryResp.data.hasOwnProperty(TOKEN_DECISION)) {
    throw 'OAuth service discovery must contain the token_decision endpoint.';
  }

  const authorizationEndpoint = discoveryResp.data[AUTHORIZATION];
  const tokenDecisionEndpoint = discoveryResp.data[TOKEN_DECISION];

  console.log(`Using authorization endpoint : ${authorizationEndpoint}.`);
  console.log(`Using token decision endpoint : ${tokenDecisionEndpoint}.`);

  /**
   * Router level middleware for login
   */
  router.get('/login', async (req, res, next) => {

    let authResponse;
    try {
      authResponse = await axios({
        method: 'get',
        url: authorizationEndpoint,
      })
    } catch (err) {
      console.log(`Error while handling the login request. Cause : ${err.message}`);
      return res.status(500).send({error: 'Internal error while handling login request.'});
    }


    if (!authResponse.data.hasOwnProperty(AUTH_URL)) {
      throw 'OAuth service discovery must contain the authorization endpoint.';
    }

    res.redirect(authResponse.data[AUTH_URL]);
    return res;
  });

  /**
   * Router level middleware for redirect handling
   */
  router.get('/oauth/redirect', async (req, res, next) => {
    const authCode = req.query.code;
    const state = req.query.state;

    if (authCode === undefined || state === undefined) {
      return res.redirect('/');
    }

    let tokensPayload = {
      code: authCode,
      state: state,
    }

    let tokenDecision;

    try {
      tokenDecision = await axios({
        method: 'post',
        url: tokenDecisionEndpoint,
        headers: {
          'Content-Type': 'application/json'
        },
        data: tokensPayload
      });
    } catch (err) {
      console.log(`Error while handling the redirect. Cause : ${err.message}`);

      if (err.response !== undefined && err.response.status === 403) {
        let response = {
          error: 'Permission denied.'
        };

        if (err.response.data.hasOwnProperty('message')) {
          response['message'] = err.response.data['message'];
        }

        return res.status(403).send(response);
      } else {
        return res.status(500).send({error: 'Error while handling the redirect.'});
      }
    }

    sessionAuthentication(req, tokenDecision.data['user']);

    return res.redirect('/');
  });

  /**
   * Router level middleware for logout
   */
  router.get('/logout', (req, res) => {
    removeSession(req);
    return res.redirect('/');
  });

  return router;
})();
