const express = require('express');
const router = express.Router();
const axios = require('axios');
const sessionAuthentication = require('./session').setAuthenticatedSession;
const removeSession = require('./session').removeSession;
const isAuthenticated = require('./session').isAuthenticated;
const getLogoutHint = require('./session').getLogoutHint;

const AUTHORIZATION = 'authorization';
const AUTH_URL = 'authorization_url';
const TOKEN_DECISION = 'token_decision';
const USER = 'user';
const LOGOUT_HINT = 'logout_hint';
const RP_LOGOUT = 'rp_logout';
const LOGOUT_URL = 'logout_path';

const prefixPath = process.env.PREFIX_PATH;

/**
 * Bridge root redirection. The exact path depends on deployment & PREFIX_PATH value
 *
 * If PREFIX_PATH is defined, redirect to <PREFIX_PATH>/bridge. Otherwise, redirect to root.
 *
 * Redirection to root will either handled by Nginx (ex:- generic keptn deployment) OR the Express layer (ex:- local bridge development).
 * */
function redirectToRoot(resp) {
  if (prefixPath !== undefined) {
    return resp.redirect(`${prefixPath}/bridge`);
  }

  return resp.redirect('/');

}

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

  let logoutEndpoint;

  if (discoveryResp.data.hasOwnProperty(RP_LOGOUT)) {
    logoutEndpoint = discoveryResp.data[RP_LOGOUT];
    console.log(`RP logout is supported by OAuth service. Using logout endpoint : ${logoutEndpoint}.`);
  }

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
      return redirectToRoot(res);
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

    sessionAuthentication(req, tokenDecision.data[USER], tokenDecision.data[LOGOUT_HINT]);
    return redirectToRoot(res);
  });

  /**
   * Router level middleware for logout
   */
  router.get('/logout', async (req, res) => {
    if (!isAuthenticated(req)) {
      // Session is not authenticated, redirect to root
      return redirectToRoot(res);
    }

    if (!logoutEndpoint) {
      removeSession(req);
      return redirectToRoot(res);
    }

    const hint = getLogoutHint(req);
    removeSession(req);

    let logoutResponse;

    try {
      logoutResponse = await axios({
        method: "post",
        url: logoutEndpoint,
        headers: {
          'Content-Type': 'application/json'
        },
        data: {
          logout_hint: hint
        }
      });
    } catch (err) {
      console.log(`Error while handling the RP logout. Cause : ${err.message}`);
      return res.status(500).send({
        error: 'Logout was successfully handled.' +
          ' However, there was an error while redirecting you to the correct endpoint.'
      });
    }

    if (!logoutResponse.data.hasOwnProperty(LOGOUT_URL)) {
      console.log('Invalid response from rp_logout.');
      return res.status(500).send({
        error: 'Logout was successfully handled.' +
          ' However, there was an error while redirecting you to the correct endpoint.'
      });
    }

    return res.redirect(logoutResponse.data[LOGOUT_URL]);

  });

  return router;
})();
