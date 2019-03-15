import * as bodyParser from 'body-parser';
import * as express from 'express';
import 'reflect-metadata';
import { Container } from 'inversify';
import {
  interfaces,
  InversifyExpressServer,
  TYPE,
} from 'inversify-express-utils';
import * as swagger from 'swagger-express-ts';

// import controllers
import './github/GitHubController';
import './docker/DockerController';
import './dynatrace/DynatraceController';

// import models

// tslint:disable-next-line: import-name
import RequestLogger = require('./middleware/requestLogger');
import authenticator = require('./middleware/authenticator');
import * as path from 'path';
import { GitHubService } from './github/GitHubService';
import { DockerService } from './docker/DockerService';
import { DynatraceService } from './dynatrace/DynatraceService';

const port: number = Number(process.env.PORT) || 5001; // or from a configuration file
const swaggerUiAssetPath = require('swagger-ui-dist').getAbsoluteFSPath();
// import models

// set up container
const container = new Container();

// set up bindings
container.bind<GitHubService>('GitHubService').to(GitHubService);
container.bind<DockerService>('DockerService').to(DockerService);
container.bind<DynatraceService>('DynatraceService').to(DynatraceService);

// create server
const server = new InversifyExpressServer(container);

server.setConfig((app: any) => {
  app.use('/api-docs/swagger', express.static(path.join(__dirname, '/src/swagger')));
  app.use('/api-docs/swagger/assets',
          express.static(
            swaggerUiAssetPath,
          ),
    );
  app.use(bodyParser.json());
  app.use(RequestLogger);
  app.use(
    swagger.express({
      definition: {
        info: {
          title: 'Keptn Eventbroker (ext) API',
          version: '0.2',
        },
        externalDocs: {
          url: '',
        },
        // Models can be defined here
      },
    }),
  );
  if (process.env.NODE_ENV === 'production') {
    app.use(authenticator);
  }
});

server.setErrorConfig((app: any) => {
  app.use(
    (
      err: Error,
      request: express.Request,
      response: express.Response,
      next: express.NextFunction,
    ) => {
      console.error(err.stack);
      response.status(500).send('Something broke!');
    },
  );
});

const app = server.build();
app.listen(port);
console.info(`Server is listening on port : ${port}`);
