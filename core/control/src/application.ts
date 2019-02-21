import express = require('express');
import bodyParser = require('body-parser');
import configRouter = require('./routes/configRouter');
import onboardRouter = require('./routes/onboardRouter');
import appRouter = require('./routes/appRouter');
import requestLogger = require('./middleware/requestLogger');
import authenticator = require('./middleware/authenticator');

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
    app.use(requestLogger);
    app.use(authenticator);
  }

  /**
   * @param app - express application
   */
  private configureRoutes(app: express.Express) {
    app.use('/config', configRouter);
    app.use('/onboard', onboardRouter);
    app.use('/app', appRouter);

    // mount more routers here
    // e.g. app.use("/organisation", organisationRouter);
  }

  public run() {
    this.app.listen(this.port);
  }
}
