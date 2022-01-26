import { Request, Response, Router } from 'express';
import { authenticateSession, getLogoutHint, isAuthenticated, removeSession } from './session';
import oClient, { BaseClient, errors, TokenSet } from 'openid-client';
import { EndSessionData } from '../../shared/interfaces/end-session-data';

const generators = oClient.generators; // else jest isn't working
const prefixPath = process.env.PREFIX_PATH;
const codeVerifiers: { [state: string]: { codeVerifier: string; nonce: string; expiresAt: number } } = {};
const stateExpireMilliSeconds = 60 * 60_000; // expires in 60 minutes

/**
 * Build the root path. The exact path depends on the deployment & PREFIX_PATH value
 *
 * If PREFIX_PATH is defined, root location is set to <PREFIX_PATH>/bridge. Otherwise, root is set to / .
 *
 * Redirection to / will be either handled by Nginx (ex:- generic keptn deployment) OR the Express layer (ex:- local bridge development).
 */
function getRootLocation(): string {
  if (prefixPath !== undefined) {
    return `${prefixPath}/bridge`;
  }
  return '/';
}

function oauthRouter(client: BaseClient, redirectUri: string, reduceRefreshDateSeconds: number): Router {
  const router = Router();
  const additionalScopes = process.env.OAUTH_SCOPE ? ` ${process.env.OAUTH_SCOPE.trim()}` : '';
  const scope = `openid${additionalScopes}`;

  /**
   * Router level middleware for login
   */
  router.get('/login', async (req: Request, res: Response) => {
    const codeVerifier = generators.codeVerifier();
    const codeChallenge = generators.codeChallenge(codeVerifier);
    const nonce = generators.nonce();
    const state = generators.state();
    codeVerifiers[state] = { codeVerifier, nonce, expiresAt: new Date().getTime() + stateExpireMilliSeconds };

    const authorizationUrl = client.authorizationUrl({
      scope,
      state,
      nonce,
      code_challenge: codeChallenge,
      code_challenge_method: 'S256',
    });

    res.redirect(authorizationUrl);
    return res;
  });

  /**
   * Router level middleware for redirect handling
   */
  router.get('/oauth/redirect', async (req: Request, res: Response) => {
    const params = client.callbackParams(req);

    if (!params.code || !params.state) {
      return res.redirect(getRootLocation());
    }

    const verifiers = codeVerifiers[params.state];
    if (verifiers) {
      delete codeVerifiers[params.state];
      params.state = undefined;
    } else {
      return res.render('error', {
        title: 'Permission denied',
        message: 'Forbidden',
      });
    }

    try {
      const tokenSet = await client.callback(redirectUri, params, {
        code_verifier: verifiers.codeVerifier,
        nonce: verifiers.nonce,
        scope,
      });
      reduceRefreshDateBy(tokenSet, reduceRefreshDateSeconds);
      await authenticateSession(req, tokenSet);
      res.redirect(getRootLocation());
    } catch (error) {
      const err = error as errors.OPError | errors.RPError;
      console.log(`Error while handling the redirect. Cause : ${err.message}`);

      if (err.response?.statusCode === 403) {
        const response = {
          title: 'Permission denied',
          message: '',
        };
        response.message =
          (err.response.body as Record<string, string>).message ?? 'User is not allowed to access the instance.';
        return res.render('error', response);
      } else {
        return res.render('error', {
          title: 'Internal error',
          message: 'Error while handling the redirect. Please retry and check whether the problem exists.',
          location: getRootLocation(),
        });
      }
    }
  });

  /**
   * Router level middleware for logout
   */
  router.post('/logout', async (req: Request, res: Response) => {
    if (!isAuthenticated(req.session)) {
      // Session is not authenticated, redirect to root
      return res.json();
    }

    const hint = getLogoutHint(req) ?? '';
    if (req.session.tokenSet.access_token && client.issuer.metadata.revocation_endpoint) {
      client.revoke(req.session.tokenSet.access_token);
    }
    removeSession(req);

    if (client.issuer.metadata.end_session_endpoint) {
      const params: EndSessionData = {
        id_token_hint: hint,
        state: generators.state(),
        post_logout_redirect_uri: redirectUri,
        end_session_endpoint: client.issuer.metadata.end_session_endpoint,
      };
      return res.json(params);
    } else {
      return res.json();
    }
  });

  router.get('/loggedOut', (req: Request, res: Response) => {
    return res.render('logout', { location: getRootLocation() });
  });

  return router;
}

/**
 * Sets the expiry date x seconds before the real one
 * @param tokenSet
 * @param seconds
 */
function reduceRefreshDateBy(tokenSet: TokenSet, seconds: number): void {
  tokenSet.expires_at = tokenSet.expires_at ? tokenSet.expires_at - seconds : undefined; // token should be refreshed x seconds earlier
}

/**
 * Delete states created by the authentication-flow if expiry is reached
 */
function validateStateTimes(): void {
  for (const key of Object.keys(codeVerifiers)) {
    if (new Date().getTime() > codeVerifiers[key]?.expiresAt) {
      delete codeVerifiers[key];
    }
  }
}

if (process.env.NODE_ENV !== 'test') {
  setInterval(validateStateTimes, 5_000); // could lead to missing exit of jest without using fakeTimers
}

export { oauthRouter, getRootLocation, reduceRefreshDateBy };
