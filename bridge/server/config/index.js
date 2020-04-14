module.exports = {
  development: {
    datastore: 'http://localhost:8085',
    configurationService: 'http://localhost:8086/v1',
  },
  production: {
    datastore: process.env.DATASTORE,
    configurationService: process.env.CONFIGURATION_SERVICE,
  },
};
