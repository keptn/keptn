import { Express, NextFunction, Request, Response } from 'express';
import { BaseClient, errors, Issuer, TokenSet } from 'openid-client';
import { SessionService } from './session';
import { getBuildableRootLocation, getRootLocation, oauthRouter, reduceRefreshDateBy } from './oauth-routes';
import { defaultContentSecurityPolicyOptions } from '../app';
import { contentSecurityPolicy } from 'helmet';
import { ComponentLogger } from '../utils/logger';
import { BridgeConfiguration } from '../interfaces/configuration';

const refreshPromises: { [sessionId: string]: Promise<TokenSet> } = {};
const reduceRefreshDateSeconds = 10;

const log = new ComponentLogger('OAuth');

async function setupOAuth(app: Express, configuration: BridgeConfiguration): Promise<SessionService> {
  const session = await new SessionService(configuration).init();
  const prefix = getBuildableRootLocation(configuration);
  let baseUrl = configuration.oauth.baseURL;
  baseUrl = baseUrl.endsWith('/') ? baseUrl.substring(0, baseUrl.length - 1) : baseUrl;
  const site = `${baseUrl}${prefix}`;
  const redirectUri = `${site}oauth/redirect`;
  const logoutUri = `${site}logoutsession`;
  // Initialise session middleware
  app.use(session.sessionRouter(app));
  const client = await setupClient(configuration, redirectUri);
  setEndSessionPolicy(app, client, configuration);
  setRoutes(app, client, redirectUri, logoutUri, session, configuration);
  return session;
}

async function setupClient(configuration: BridgeConfiguration, redirectUri: string): Promise<BaseClient> {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const ssoIssuer = await (global.issuer ?? Issuer).discover(configuration.oauth.discoveryURL);
  const clientSecret = configuration.oauth.clientSecret;

  if (!ssoIssuer.metadata.authorization_endpoint) {
    throw Error('OAuth service discovery must contain the authorization endpoint.');
  }

  if (!ssoIssuer.metadata.token_endpoint) {
    throw Error('OAuth service discovery must contain the token_decision endpoint.');
  }

  log.info(`Using authorization endpoint : ${ssoIssuer.metadata.authorization_endpoint}.`);
  log.info(`Using token decision endpoint : ${ssoIssuer.metadata.token_endpoint}.`);

  if (ssoIssuer.metadata.end_session_endpoint) {
    log.info(
      `RP logout is supported by OAuth service. Using logout endpoint : ${ssoIssuer.metadata.end_session_endpoint}.`
    );
  }

  return new ssoIssuer.Client({
    client_id: configuration.oauth.clientID,
    ...(clientSecret && { client_secret: clientSecret }),
    redirect_uris: [redirectUri],
    response_types: ['code'],
    token_endpoint_auth_method: clientSecret ? 'client_secret_basic' : 'none',
    id_token_signed_response_alg: configuration.oauth.tokenAlgorithm,
  });
}

function setEndSessionPolicy(app: Express, client: BaseClient, configuration: BridgeConfiguration): void {
  if (client.issuer.metadata.end_session_endpoint && defaultContentSecurityPolicyOptions.directives) {
    defaultContentSecurityPolicyOptions.directives['form-action'] = [
      `'self'`,
      client.issuer.metadata.end_session_endpoint,
      configuration.oauth.allowedLogoutURL,
    ];
    app.use(contentSecurityPolicy(defaultContentSecurityPolicyOptions));
  }
}

function setRoutes(
  app: Express,
  client: BaseClient,
  redirectUri: string,
  logoutUri: string,
  session: SessionService,
  configuration: BridgeConfiguration
): void {
  // Initializing OAuth middleware.
  app.use(oauthRouter(client, redirectUri, logoutUri, reduceRefreshDateSeconds, session, configuration));
  // Authentication filter for API requests
  app.use('/api', async (req: Request, resp: Response, next: NextFunction) => {
    if (!session.isAuthenticated(req.session)) {
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
              log.error(`Error while saving tokenSet. Cause: ${error?.message}`);
            }
            delete refreshPromises[req.session.id];
            next();
          });
        } catch (error) {
          const err = error as errors.OPError | errors.RPError;
          log.error(`Renewal of access_token did not work. Cause ${err.message}`);

          delete refreshPromises[req.session.id];
          session.removeSession(req);
          resp.redirect(getRootLocation(configuration.features.prefixPath));
        }
      } else {
        return next();
      }
    }
  });
}

export { setupOAuth };
