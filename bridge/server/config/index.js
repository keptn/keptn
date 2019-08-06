module.exports = {
  development: {
    elasticEndpoint: 'http://localhost:8001/api/v1/namespaces/knative-monitoring/services/elasticsearch-logging/proxy',
    datastore: 'http://localhost:8080',
  },
  production: {
    elasticEndpoint: process.env.ELASTIC_ENDPOINT,
    datastore: process.env.DATASTORE,
  },
};
