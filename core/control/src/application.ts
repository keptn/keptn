import express = require('express');
import bodyParser = require('body-parser');
import configRouter = require('./routes/ConfigRouter');
import projectRouter = require('./routes/ProjectRouter');
import RequestLogger = require('./middleware/RequestLogger');
import Authenticator = require('./middleware/Authenticator');

export class WebApi {

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

    // mount more routers here
    // e.g. app.use("/organisation", organisationRouter);
  }

  public run() {
    this.app.listen(this.port);
  }
}
