import { isAuthenticated, removeSession, sessionRouter } from './session';
import { Express, NextFunction, Request, Response } from 'express';
import { getRootLocation, oauthRouter, reduceRefreshDateBy } from './oauth-routes';
import { BaseClient, errors, Issuer, TokenSet } from 'openid-client';

const refreshPromises: { [sessionId: string]: Promise<TokenSet> } = {};
const reduceRefreshDateSeconds = 10;

async function setupOAuth(app: Express, discoveryEndpoint: string, clientId: string, baseUrl: string): Promise<void> {
  let prefix = getRootLocation();
  baseUrl = baseUrl.endsWith('/') ? baseUrl.substring(0, baseUrl.length - 1) : baseUrl;
  prefix = prefix.endsWith('/') ? prefix : `${prefix}/`;
  const redirectUri = `${baseUrl}${prefix}oauth/redirect`;
  // Initialise session middleware
  app.use(sessionRouter(app));
  const client = await setupClient(discoveryEndpoint, clientId, redirectUri);
  setRoutes(app, client, redirectUri);
}

async function setupClient(discoveryEndpoint: string, clientId: string, redirectUri: string): Promise<BaseClient> {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const ssoIssuer = await (global.issuer ?? Issuer).discover(discoveryEndpoint);
  const clientSecret = process.env.OAUTH_CLIENT_SECRET;

  if (!ssoIssuer.metadata.authorization_endpoint) {
    throw Error('OAuth service discovery must contain the authorization endpoint.');
  }

  if (!ssoIssuer.metadata.token_endpoint) {
    throw Error('OAuth service discovery must contain the token_decision endpoint.');
  }

  console.log(`Using authorization endpoint : ${ssoIssuer.metadata.authorization_endpoint}.`);
  console.log(`Using token decision endpoint : ${ssoIssuer.metadata.token_endpoint}.`);

  if (ssoIssuer.metadata.end_session_endpoint) {
    console.log(
      `RP logout is supported by OAuth service. Using logout endpoint : ${ssoIssuer.metadata.end_session_endpoint}.`
    );
  }

  return new ssoIssuer.Client({
    client_id: clientId,
    ...(clientSecret && { client_secret: clientSecret }),
    redirect_uris: [redirectUri],
    response_types: ['code'],
    token_endpoint_auth_method: clientSecret ? 'client_secret_basic' : 'none',
    id_token_signed_response_alg: process.env.OAUTH_ID_TOKEN_ALG || 'RS256',
  });
}

function setRoutes(app: Express, client: BaseClient, redirectUri: string): void {
  // Initializing OAuth middleware.
  app.use(oauthRouter(client, redirectUri, reduceRefreshDateSeconds));
  // Authentication filter for API requests
  app.use('/api', async (req: Request, resp: Response, next: NextFunction) => {
    if (!isAuthenticated(req.session)) {
      return next({ response: { status: 401 } });
    } else {
      const tokenSet = new TokenSet(req.session.tokenSet);
      if (tokenSet.expired()) {
        refreshPromises[req.session.id] ??= client.refresh(tokenSet).then((token) => {
          reduceRefreshDateBy(token, reduceRefreshDateSeconds);
          return token;
        });
        try {
          req.session.tokenSet = await refreshPromises[req.session.id];
          req.session.save((error) => {
            if (error) {
              console.log(`Error while saving tokenSet. Cause: ${error?.message}`);
            }
            delete refreshPromises[req.session.id];
            next();
          });
        } catch (error) {
          const err = error as errors.OPError | errors.RPError;
          console.error(`Renewal of access_token did not work. Cause ${err.message}`);

          delete refreshPromises[req.session.id];
          removeSession(req);
          resp.redirect(getRootLocation());
        }
      } else {
        return next();
      }
    }
  });
}

export { setupOAuth };
