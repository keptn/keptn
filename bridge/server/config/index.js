module.exports = {
  development: {
    datastore: 'http://localhost:8080',
  },
  production: {
    datastore: process.env.DATASTORE,
  },
};
