/* eslint-disable no-underscore-dangle */
const { Client } = require('@elastic/elasticsearch');

class ElasticService {
  constructor(endpoint) {
    this.elastic = new Client({ node: endpoint });
  }

  async getRoots() {
    const result = await this.elastic.search({
      index: 'logstash-*',
      body: {
        query: {
          match: {
            keptnEntry: {
              query: true,
            },
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

  async getTraces(contextId) {
    const result = await this.elastic.search({
      index: 'logstash-*',
      body: {
        from: 0,
        size: 999,
        query: {
          match: {
            keptnContext: {
              query: contextId,
            },
          },
        },
        sort: { '@timestamp': 'desc' },
      },
    });

    return result.body.hits.hits.map(hit => ({
      ...hit,
    })).filter(elm => elm._source.keptnService !== 'eventbroker');
  }
}

module.exports = ElasticService;
