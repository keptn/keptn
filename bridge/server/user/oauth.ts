import { Request, Response, Router } from 'express';
import axios from 'axios';
import { authenticateSession, getLogoutHint, isAuthenticated, removeSession } from './session';

const router = Router();
const AUTHORIZATION = 'authorization';
const AUTH_URL = 'authorization_url';
const TOKEN_DECISION = 'token_decision';
const USER = 'user';
const LOGOUT_HINT = 'logout_hint';
const RP_LOGOUT = 'rp_logout';
const LOGOUT_URL = 'logout_path';
const prefixPath = process.env.PREFIX_PATH;

/**
 * Build the root path. The exact path depends on the deployment & PREFIX_PATH value
 *
 * If PREFIX_PATH is defined, root location is set to <PREFIX_PATH>/bridge. Otherwise, root is set to / .
 *
 * Redirection to / will be either handled by Nginx (ex:- generic keptn deployment) OR the Express layer (ex:- local bridge development).
 */
function getRootLocation() {
  if (prefixPath !== undefined) {
    return `${prefixPath}/bridge`;
  }

  return '/';
}

async function oauthRouter() {
  console.log('Enabling OAuth for bridge.');

  const discoveryEndpoint = process.env.OAUTH_DISCOVERY;

  if (!discoveryEndpoint) {
    throw Error('OAUTH_DISCOVERY must be defined when oauth based login (OAUTH_ENABLED) is activated.' +
    ' Please check your environment variables.');
  }

  const discoveryResp = await axios({
    method: 'get',
    url: discoveryEndpoint,
  });

  if (discoveryResp.status !== 200) {
    throw Error(`Invalid oauth service discovery response. Received status : ${discoveryResp.statusText}.`);
  }

  if (!discoveryResp.data.hasOwnProperty(AUTHORIZATION)) {
    throw Error('OAuth service discovery must contain the authorization endpoint.');
  }

  if (!discoveryResp.data.hasOwnProperty(TOKEN_DECISION)) {
    throw Error('OAuth service discovery must contain the token_decision endpoint.');
  }

  const authorizationEndpoint = discoveryResp.data[AUTHORIZATION];
  const tokenDecisionEndpoint = discoveryResp.data[TOKEN_DECISION];

  console.log(`Using authorization endpoint : ${authorizationEndpoint}.`);
  console.log(`Using token decision endpoint : ${tokenDecisionEndpoint}.`);

  let logoutEndpoint = '';

  if (discoveryResp.data.hasOwnProperty(RP_LOGOUT)) {
    logoutEndpoint = discoveryResp.data[RP_LOGOUT];
    console.log(`RP logout is supported by OAuth service. Using logout endpoint : ${logoutEndpoint}.`);
  }

  /**
   * Router level middleware for login
   */
  router.get('/login', async (req: Request, res: Response) => {

    let authResponse;
    try {
      authResponse = await axios({
        method: 'get',
        url: authorizationEndpoint,
      });
    } catch (err) {
      console.log(`Error while handling the login request. Cause : ${err.message}`);
      return res.render('error',
        {
          title: 'Internal error',
          message: 'Error while handling the login request.',
          location: getRootLocation()
        });
    }


    if (!authResponse.data.hasOwnProperty(AUTH_URL)) {
      return res.render('error',
        {
          title: 'Invalid state',
          message: 'Failure to obtain login details.',
          location: getRootLocation()
        });
    }

    res.redirect(authResponse.data[AUTH_URL]);
    return res;
  });

  /**
   * Router level middleware for redirect handling
   */
  router.get('/oauth/redirect', async (req: Request, res: Response) => {
    const authCode = req.query.code;
    const state = req.query.state;

    if (authCode === undefined || state === undefined) {
      return res.redirect(getRootLocation());
    }

    const tokensPayload = {
      code: authCode,
      state,
    };

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
        const response = {
          title: 'Permission denied',
          message: ''
        };

        if (err.response.data.hasOwnProperty('message')) {
          response.message = err.response.data.message;
        } else {
          response.message = 'User is not allowed access the instance.';
        }

        return res.render('error', response);
      } else {
        return res.render('error',
          {
            title: 'Internal error',
            message: 'Error while handling the redirect. Please retry and check whether the problem exists.',
            location: getRootLocation()
          });
      }
    }

    authenticateSession(req, tokenDecision.data[USER], tokenDecision.data[LOGOUT_HINT],
      () => res.redirect(getRootLocation()));
  });

  /**
   * Router level middleware for logout
   */
  router.get('/logout', async (req: Request, res: Response) => {
    if (!isAuthenticated(req)) {
      // Session is not authenticated, redirect to root
      return res.redirect(getRootLocation());
    }

    if (!logoutEndpoint) {
      removeSession(req);
      return res.redirect(getRootLocation());
    }

    const hint = getLogoutHint(req);
    removeSession(req);

    let logoutResponse;

    try {
      logoutResponse = await axios({
        method: 'post',
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

      return res.render('error',
        {
          title: 'Internal error',
          message: 'Logout was successfully handled.' +
            ' However, there was an error while redirecting you to the correct endpoint.',
          location : getRootLocation(),
          locationMessage: 'Home'
        });
    }

    if (!logoutResponse.data.hasOwnProperty(LOGOUT_URL)) {
      console.log('Invalid response from rp_logout.');
      return res.render('error',
        {
          title: 'Internal error',
          message: 'Logout was successfully handled.' +
            ' However, there was an error while redirecting you to the correct endpoint.',
          location : getRootLocation(),
          locationMessage: 'Home'
        });
    }

    return res.redirect(logoutResponse.data[LOGOUT_URL]);

  });

  return router;
}

const authRouter = oauthRouter();
export { authRouter as oauthRouter };
