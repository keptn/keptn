// server.js
var express = require('express');
var bodyParser = require('body-parser');
var app = express();

// parse application/json
app.use(bodyParser.json())

var gitHubListener = require('./app/event-broker.js').githubWebhookListener;
var jenkinsListener = require('./app/event-broker.js').jenkinsNotificationListener;
var cloudEventsListener = require('./app/event-broker.js').cloudEventsListener;

app.post('/github', function(req, res) {
  gitHubListener(req, res);
});
app.post('/jenkins', jenkinsListener);

app.post('/', cloudEventsListener);

app.get('/health', function (req, res, next) {
    // check my health
    res.sendStatus(200)
  });

var server = app.listen(process.env.PORT || 8079, function () {
    var port = server.address().port;
    console.log("Keptn Event Broker now running in %s mode on port %d", app.get("env"), port);
    console.log(`Channel URI: ${process.env.CHANNEL_URI}`);
  });