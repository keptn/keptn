import express = require('express');
import bodyParser = require('body-parser');
import configRouter = require('./routes/configRouter');
import projectRouter = require('./routes/ProjectRouter');
import serviceRouter = require('./routes/serviceRouter');
import RequestLogger = require('./middleware/requestLogger');
import Authenticator = require('./middleware/authenticator');

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
    app.use('/config', configRouter);
    app.use('/project', projectRouter);
    app.use('/service', serviceRouter);

    // mount more routers here
    // e.g. app.use("/organisation", organisationRouter);
  }

  public run() {
    this.app.listen(this.port);
  }
}
