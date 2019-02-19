import express = require("express");
import bodyParser = require("body-parser");
import customerRouter = require("./routes/authRouter");
import requestLogger = require("./middleware/requestLogger");

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
    }

    private configureRoutes(app: express.Express) {
        app.use("/auth", customerRouter );
        // mount more routers here
        // e.g. app.use("/organisation", organisationRouter);
    }

    public run() {
        this.app.listen(this.port);  
    }
}