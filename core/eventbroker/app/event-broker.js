var axios = require("axios");

var verifyGithubWebhook = require('verify-github-webhook');

const KEPTN_EVENTS = [
  {
    eventType: "sh.keptn.events.new-artefact",
    channelUri: process.env.NEW_ARTEFACT_CHANNEL
  },
  {
    eventType: "sh.keptn.events.start-deployment",
    channelUri: process.env.START_DEPLOYMENT_CHANNEL
  },
  {
    eventType: "sh.keptn.events.deployment-finished",
    channelUri: process.env.DEPLOYMENT_FINISHED_CHANNEL
  },
  {
    eventType: "sh.keptn.events.start-tests",
    channelUri: process.env.START_TESTS_CHANNEL
  },
  {
    eventType: "sh.keptn.events.tests-finished",
    channelUri: process.env.TESTS_FINISHED_CHANNEL
  },
  {
    eventType: "sh.keptn.events.start-evaluation",
    channelUri: process.env.START_EVALUATION_CHANNEL
  },
  {
    eventType: "sh.keptn.events.evaluation-finished",
    channelUri: process.env.EVALUATION_FINISHED_CHANNEL
  }
]

async function sendMessage(msg, eventType) {
  let channelUri;
  if (eventType !== undefined) {
    const eventSpec = KEPTN_EVENTS.find(item => item.eventType === eventType);
    if (eventSpec !== undefined) {
      channelUri = eventSpec.channelUri;
    } else {
      console.log(`No event found for eventType ${eventType}`);
      return;
    }
  } else {
    channelUri = process.env.CHANNEL_URI;
  }
  console.log(`Sending message to ${process.env.CHANNEL_URI}`);
  var config = {
    method: 'POST',
    url: `http://${channelUri}`,
    data: msg
  };

  axios.request(config).then(response => {
    console.log(`Sent message, received response: ${response}`);
  }).catch(e => {
    console.log(`Error while sending message: ${e}`);
  }); 
}

module.exports.jenkinsNotificationListener = (event, response) => {
  console.log(event.body);
  const msg = {
    channel: 'jenkins',
    body: event.body
  };
  sendMessage(msg);
  response.writeHeader(200); response.end();
}

module.exports.githubWebhookListener = async (event, response) => {
  console.log(event.body);
  var signature = event.headers['X-Hub-Signature'];
  console.log(`Verifying signature ${signature}`);
  /*
  if (signature === undefined) {
    response.writeHeader(403); response.end;
    return;
  }
  if (!verifyGithubWebhook(signature, event.body, process.env.GITHUB_WEBHOOK_SECRET)) {
    response.writeHeader(403); response.end;
    return;
  }
  */
  const githubEvent = event.get('X-GitHub-Event');
  const msg = {
    channel: 'github',
    eventType: githubEvent,
    body: event.body
  };
  sendMessage(msg);
  response.writeHeader(200); response.end();
};

module.exports.cloudEventsListener = async (event, response) => {
  console.log(event.body);
  sendMessage(event.body, event.body.type);
  response.writeHeader(200); response.end();
}
