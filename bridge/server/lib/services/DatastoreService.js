/* eslint-disable no-underscore-dangle */
const axios = require('axios');

class DatastoreService {
  constructor(endpoint) {
    this.api = endpoint;
  }


  static mapEvent(event) {
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
      default: mappedEvent.eventTypeHeadline = event.type; break;
    }

    if (event.source === 'https://github.com/keptn/keptn/remediation-service') {
      mappedEvent.eventTypeHeadline = 'Remediation';
    }

    return mappedEvent;
  }

  async getRoots() {
    const deploymentRoots = await this.getDeploymentRoots();
    const problemRoots = await this.getProblemRoots();
    const combinedRoots = deploymentRoots.concat(problemRoots);
    combinedRoots.sort((a, b) => (a.timestamp < b.timestamp ? 1 : -1));
    return combinedRoots;
  }

  async getDeploymentRoots() {
    const url = `${this.api}/event?type=sh.keptn.event.configuration.change&pageSize=100`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event)).filter(e => e.data.stage === '');
    }
    return [];
  }

  async getProblemRoots() {
    const url = `${this.api}/event?type=sh.keptn.event.configuration.change&pageSize=100`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event)).filter(e => e.source.includes('remediation-service'));
    }
    return [];
  }

  async getTraces(contextId) {
    const url = `${this.api}/event?keptnContext=${contextId}&pageSize=100`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      const traces = data.events.map(event => DatastoreService.mapEvent(event));
      traces.sort((a, b) => (a.timestamp > b.timestamp ? 1 : -1));
      return traces;
    }
    return [];
  }

  async findRoots(contextId) {
    const deploymentRoots = await this.findDeploymentRoots(contextId);
    const problemRoots = await this.findProblemRoots(contextId);
    const combinedRoots = deploymentRoots.concat(problemRoots);
    combinedRoots.sort((a, b) => (a.timestamp < b.timestamp ? 1 : -1));
    return combinedRoots;
  }

  async findDeploymentRoots(contextId) {
    const url = `${this.api}/event?keptnContext=${contextId}&type=sh.keptn.event.configuration.change&pageSize=10`;
    console.log(url);
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event)).filter(e => e.data.stage === '');
    }
    return [];
  }

  async findProblemRoots(contextId) {
    const url = `${this.api}/event?keptnContext=${contextId}&type=sh.keptn.event.problem.open&pageSize=100`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event)).filter(e => e.data.state === 'OPEN');
    }
    return [];
  }
}

module.exports = DatastoreService;
