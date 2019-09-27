module.exports = {
  development: {
    datastore: 'http://localhost:8085',
  },
  production: {
    datastore: process.env.DATASTORE,
  },
};
