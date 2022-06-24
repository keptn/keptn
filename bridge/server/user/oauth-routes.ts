import { Request, Response, Router } from 'express';
import { BaseClient, errors, generators, TokenSet } from 'openid-client';
import { EndSessionData } from '../../shared/interfaces/end-session-data';
import { SessionService } from './session';
import { ComponentLogger } from '../utils/logger';

const prefixPath = process.env.PREFIX_PATH;

const log = new ComponentLogger('OAuth');

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

/**
 * returns the root location and ends with "/"
 */
export function getBuildableRootLocation(): string {
  const prefix = getRootLocation();
  return prefix.endsWith('/') ? prefix : `${prefix}/`;
  // currently "/bridge/" leads to an empty page. That's why getRootLocation is not changed to end with "/"
}

function oauthRouter(
  client: BaseClient,
  redirectUri: string,
  logoutUri: string,
  reduceRefreshDateSeconds: number,
  session: SessionService
): Router {
  const router = Router();
  const additionalScopes = process.env.OAUTH_SCOPE ? ` ${process.env.OAUTH_SCOPE.trim()}` : '';
  const scope = `openid${additionalScopes}`;
  log.info(`Using scope: ${scope}`);

  /**
   * Router level middleware for login
   */
  router.get('/oauth/login', async (_req: Request, res: Response) => {
    const codeVerifier = generators.codeVerifier();
    const codeChallenge = generators.codeChallenge(codeVerifier);
    const nonce = generators.nonce();
    const state = generators.state();
    try {
      await session.saveValidationData(state, codeVerifier, nonce);

      const authorizationUrl = client.authorizationUrl({
        scope,
        state,
        nonce,
        code_challenge: codeChallenge,
        code_challenge_method: 'S256',
      });
      res.redirect(authorizationUrl);
    } catch (e) {
      const msg = e instanceof Error ? `${e.name}: ${e.message}` : `${e}`;
      log.error(msg);
    }
    return res;
  });

  /**
   * Router level middleware for redirect handling
   */
  router.get('/oauth/redirect', async (req: Request, res: Response) => {
    const params = client.callbackParams(req);
    const errorPageUrl = `${getBuildableRootLocation()}error`;

    if (!params.code || !params.state) {
      return res.redirect(errorPageUrl);
    }

    try {
      const validationData = await session.getAndRemoveValidationData(params.state);
      if (!validationData) {
        return res.redirect(errorPageUrl);
      }

      const tokenSet = await client.callback(redirectUri, params, {
        code_verifier: validationData.codeVerifier,
        nonce: validationData.nonce,
        state: validationData._id,
        scope,
      });
      reduceRefreshDateBy(tokenSet, reduceRefreshDateSeconds);
      await session.authenticateSession(req, tokenSet);
      res.redirect(getRootLocation());
    } catch (error) {
      const err = error as errors.OPError | errors.RPError;
      log.error(`Error while handling the redirect. Cause : ${err.message}`);

      if (err.response?.statusCode === 403) {
        return res.redirect(`${errorPageUrl}?status=403`);
      } else {
        return res.redirect(errorPageUrl);
      }
    }
  });

  /**
   * Router level middleware for logout
   */
  router.post('/oauth/logout', async (req: Request, res: Response) => {
    if (!session.isAuthenticated(req.session)) {
      // Session is not authenticated, redirect to root
      return res.json();
    }

    const hint = session.getLogoutHint(req) ?? '';
    if (req.session.tokenSet.access_token && client.issuer.metadata.revocation_endpoint) {
      client.revoke(req.session.tokenSet.access_token);
    }
    session.removeSession(req);

    if (client.issuer.metadata.end_session_endpoint) {
      const params: EndSessionData = {
        id_token_hint: hint,
        state: generators.state(),
        post_logout_redirect_uri: logoutUri,
        end_session_endpoint: client.issuer.metadata.end_session_endpoint,
      };
      return res.json(params);
    } else {
      return res.json();
    }
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

export { oauthRouter, getRootLocation, reduceRefreshDateBy };
