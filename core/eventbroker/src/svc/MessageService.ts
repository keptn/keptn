import { injectable } from 'inversify';
import axios from 'axios';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';

@injectable()
export class MessageService {

  static KEPTN_EVENTS = [
    {
      eventType: 'sh.keptn.events.new-artefact',
      channelUri: process.env.NEW_ARTEFACT_CHANNEL,
    },
    {
      eventType: 'sh.keptn.events.start-deployment',
      channelUri: process.env.START_DEPLOYMENT_CHANNEL,
    },
    {
      eventType: 'sh.keptn.events.deployment-finished',
      channelUri: process.env.DEPLOYMENT_FINISHED_CHANNEL,
    },
    {
      eventType: 'sh.keptn.events.start-tests',
      channelUri: process.env.START_TESTS_CHANNEL,
    },
    {
      eventType: 'sh.keptn.events.tests-finished',
      channelUri: process.env.TESTS_FINISHED_CHANNEL,
    },
    {
      eventType: 'sh.keptn.events.start-evaluation',
      channelUri: process.env.START_EVALUATION_CHANNEL,
    },
    {
      eventType: 'sh.keptn.events.evaluation-finished',
      channelUri: process.env.EVALUATION_FINISHED_CHANNEL,
    },
  ];

  constructor() {}

  public async sendMessage(message: KeptnRequestModel): Promise<boolean> {
    let channelUri;
    const eventType = message.type;
    if (eventType !== undefined) {
      const eventSpec = MessageService.KEPTN_EVENTS.find(item => item.eventType === eventType);
      if (eventSpec !== undefined) {
        channelUri = eventSpec.channelUri;
      } else {
        console.log(`No event found for eventType ${eventType}`);
        return;
      }
    }
    console.log(`Sending message to ${process.env.CHANNEL_URI}`);

    axios.post(`http://${channelUri}`, message).then().catch(() => {});

    return true;
  }
}
