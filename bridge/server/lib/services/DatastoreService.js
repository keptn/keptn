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
      project: event.data.project,
      service: event.data.service,
      stage: event.data.stage,
      tag: event.data.tag,
    };
  }

  async getRoots() {
    const url = `${this.api}/events/type/newartifact`;
    const result = await axios.get(url);
    const { data } = result;
    return data.map(event => DatastoreService.mapEvent(event));
  }

  async getTraces(contextId) {
    const url = `${this.api}/events/id/${contextId}`;
    const result = await axios.get(url);
    const { data } = result;
    return data.map(event => DatastoreService.mapEvent(event));
  }

  async findRoots(keptnContext) {
    const result = await this.elastic.search({
      index: 'logstash-*',
      body: {
        from: 0,
        size: 20,
        query: {
          bool: {
            must: [
              { match_all: {} },
              { match_phrase: { keptnEntry: { query: true } } },
              {
                match_phrase: {
                  keptnContext: { query: keptnContext },
                },
              },
            ],
            must_not: [],
          },
        },
        sort: { '@timestamp': 'desc' },
        _source: ['message', 'keptnContext', '@timestamp'],
      },
    });

    return result.body.hits.hits.map(hit => ({
      timestamp: hit._source['@timestamp'],
      keptnContext: hit._source.keptnContext,
      message: hit._source.message,
    }));
  }
}

module.exports = DatastoreService;
