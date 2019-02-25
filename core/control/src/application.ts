import express = require('express');
import bodyParser = require('body-parser');
// tslint:disable-next-line: import-name
import ConfigRouter from './routes/configRouter';
import ProjectRouter = require('./routes/ProjectRouter');
// tslint:disable-next-line: import-name
import RequestLogger = require('./middleware/requestLogger');
// tslint:disable-next-line: import-name
import Authenticator = require('./middleware/authenticator');
import * as path from 'path';

import * as swagger from 'swagger-express-ts';
const swaggerUiAssetPath = require('swagger-ui-dist').getAbsoluteFSPath();

export class WebApi {

  private swaggerSpec: any;

  /**
   * @param app - express application
   * @param port - port to listen on
   */
  constructor(private app: express.Express, private port: number) {
    this.configureMiddleware(app);
    this.configureRoutes(app);
  }

  /**
   * @param app - express application
   */
  private configureMiddleware(app: express.Express) {
    app.use('/api-docs/swagger', express.static(path.join(__dirname, '/src/swagger')));
    app.use('/api-docs/swagger/assets',
            express.static(
              swaggerUiAssetPath,
            ),
      );
    app.use(bodyParser.json());
    app.use(RequestLogger);
    if (process.env.NODE_ENV === 'production') {
      app.use(Authenticator);
    }
  }

  /**
   * @param app - express application
   */
  private configureRoutes(app: express.Express) {
    app.use(swagger.express(
      {
        definition: {
          info: {
            title: 'Keptn Control API',
            version: '0.2',
          },
          externalDocs: {
            url: '',
          },
          // Models can be defined here
        },
      },
    ));
    app.use('/config', ConfigRouter);
    app.use('/project', ProjectRouter);

    // mount more routers here
    // e.g. app.use("/organisation", organisationRouter);
  }

  public run() {
    this.app.listen(this.port);
  }
}
