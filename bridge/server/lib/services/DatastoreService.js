/* eslint-disable no-underscore-dangle */
const axios = require('axios');

class DatastoreService {
  constructor(endpoint) {
    this.api = endpoint;
  }


  static mapEvent(event) {
    return {
      timestamp: event.time,
      type: event.type,
      keptnContext: event.shkeptncontext,
      data: event.data,
      id: event.id,
    };
  }

  async getRoots() {
    const url = `${this.api}/event?type=sh.keptn.events.new-artifact&pageSize=100`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event));
    }
    return [];
  }

  async getTraces(contextId) {
    const url = `${this.api}/event?keptnContext=${contextId}&pageSize=100`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event));
    }
    return [];
  }

  async findRoots(contextId) {
    const url = `${this.api}/event?keptnContext=${contextId}&type=sh.keptn.events.new-artifact&pageSize=10`;
    console.log(url);
    const result = await axios.get(url);
    const { data } = result;
    if (data.events) {
      return data.events.map(event => DatastoreService.mapEvent(event));
    }
    return [];
  }
}

module.exports = DatastoreService;
