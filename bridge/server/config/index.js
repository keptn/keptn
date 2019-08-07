module.exports = {
  development: {
    elasticEndpoint: 'http://localhost:8001/api/v1/namespaces/knative-monitoring/services/elasticsearch-logging/proxy',
  },
  production: {
    elasticEndpoint: process.env.ELASTIC_ENDPOINT,
  },
};
