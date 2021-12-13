import { Request, Response, Router } from 'express';
import { authenticateSession, getLogoutHint, isAuthenticated, removeSession } from './session';
import { BaseClient, errors, generators, TokenSet } from 'openid-client';

const prefixPath = process.env.PREFIX_PATH;
const codeVerifiers: { [state: string]: { codeVerifier: string; nonce: string; expiresAt: number } } = {};

/**
 * Build the root path. The exact path depends on the deployment & PREFIX_PATH value
 *
 * If PREFIX_PATH is defined, root location is set to <PREFIX_PATH> or <PREFIX_PATH>/bridge for production environments. Otherwise, root is set to / .
 *
 * Redirection to / will be either handled by Nginx (ex:- generic keptn deployment) OR the Express layer (ex:- local bridge development).
 */
function getRootLocation(): string {
  if (prefixPath !== undefined) {
    return process.env.NODE_ENV === 'production' ? `${prefixPath}/bridge` : prefixPath;
  }

  return '/';
}

function oauthRouter(client: BaseClient, redirectUri: string, reduceRefreshDateSeconds: number): Router {
  const router = Router();

  /**
   * Router level middleware for login
   */
  router.get('/login', async (req: Request, res: Response) => {
    const codeVerifier = generators.codeVerifier();
    const codeChallenge = generators.codeChallenge(codeVerifier);
    const nonce = generators.nonce();
    const state = generators.state();
    codeVerifiers[state] = { codeVerifier, nonce, expiresAt: new Date().getTime() + 5 * 60_000 }; // expires in 5 minutes

    const authorizationUrl = client.authorizationUrl({
      scope: 'openid',
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
        title: 'Error',
        message: 'Error while handling request. State does not exist.',
        location: getRootLocation(),
      });
    }

    try {
      const tokenSet = await client.callback(redirectUri, params, {
        code_verifier: verifiers.codeVerifier,
        nonce: verifiers.nonce,
        scope: 'openid',
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
          (err.response.body as Record<string, string>).message ?? 'User is not allowed access the instance.';
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
  router.get('/logout', async (req: Request, res: Response) => {
    if (!isAuthenticated(req.session)) {
      // Session is not authenticated, redirect to root
      return res.redirect(getRootLocation());
    }

    let logoutUrl;
    if (client.issuer.metadata.end_session_endpoint) {
      logoutUrl = client.endSessionUrl({
        id_token_hint: getLogoutHint(req),
        state: generators.state(),
        post_logout_redirect_uri: redirectUri,
      });
    } else {
      logoutUrl = getRootLocation();
    }
    removeSession(req);
    return res.redirect(logoutUrl);
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

setInterval(validateStateTimes, 5_000);

export { oauthRouter, getRootLocation, reduceRefreshDateBy };
