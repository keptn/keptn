/* eslint-disable no-underscore-dangle */
const axios = require('axios');

class DatastoreService {
  constructor(endpoint) {
    this.api = endpoint;
  }

  static mapEventsResult(result, sortCompareCallback) {
    const { data } = result;
    if (data.events) {
      const events = data.events.map(event => DatastoreService.mapEvent(event));
      if(sortCompareCallback)
        events.sort(sortCompareCallback);
      return events;
    }
    return [];
  }

  static mapEvent(event) {
    return event;

    // TODO: check if this mappedEvent is necessary
    // eventTypeHeadline should be translated on client side
    const mappedEvent = {
      timestamp: event.time,
      type: event.type,
      keptnContext: event.shkeptncontext,
      data: event.data,
      id: event.id,
      source: event.source,
      plainEvent: JSON.stringify(event, null, 2),
    };

    switch (mappedEvent.type) {
      case 'sh.keptn.event.configuration.change': mappedEvent.eventTypeHeadline = 'Configuration change'; break;
      case 'sh.keptn.event.problem.open': mappedEvent.eventTypeHeadline = 'Problem'; break;
      case 'sh.keptn.events.deployment-finished': mappedEvent.eventTypeHeadline = 'Deployment finished'; break;
      case 'sh.keptn.events.evaluation-done': mappedEvent.eventTypeHeadline = 'Evaluation done'; break;
      case 'sh.keptn.events.tests-finished': mappedEvent.eventTypeHeadline = 'Tests finished'; break;
      case 'sh.keptn.event.start-evaluation': mappedEvent.eventTypeHeadline = 'Start Evaluation'; break;
      case 'sh.keptn.internal.event.get-sli': mappedEvent.eventTypeHeadline = 'Start SLI retrieval'; break;
      case 'sh.keptn.internal.event.get-sli.done': mappedEvent.eventTypeHeadline = 'SLI retrieval done'; break;

      default: mappedEvent.eventTypeHeadline = event.type; break;
    }

    if (event.source === 'https://github.com/keptn/keptn/remediation-service') {
      mappedEvent.eventTypeHeadline = 'Remediation';
    }

    return mappedEvent;
  }

  async getRoots(projectName, serviceName) {
    const url = `${this.api}/event?root=true&project=${projectName}&service=${serviceName}&pageSize=100`;
    const result = await axios.get(url);
    return DatastoreService.mapEventsResult(result, (a, b) => (a.timestamp < b.timestamp ? 1 : -1));
  }

  async getTraces(contextId) {
    const url = `${this.api}/event?keptnContext=${contextId}&pageSize=100`;
    const result = await axios.get(url);
    return DatastoreService.mapEventsResult(result, (a, b) => (a.timestamp > b.timestamp ? 1 : -1));
  }

}

module.exports = DatastoreService;
