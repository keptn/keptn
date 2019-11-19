module.exports = {
  development: {
    datastore: 'http://localhost:8081',
  },
  production: {
    datastore: process.env.DATASTORE,
  },
};
