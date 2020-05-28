module.exports = {
  development: {
    apiUrl: 'http://localhost:8088/'
  },
  production: {
    apiUrl: process.env.API_URL
  },
};
