import axios from 'axios';

export default {
  async fetchRoots() {
    const response = await axios.get('/api/roots');
    return response.data;
  },
  async fetchTraces(contextId) {
    const response = await axios.get(`/api/traces/${contextId}`);
    return response.data;
  },
  async findRoots(contextId) {
    const response = await axios.get(`/api/roots/${contextId}`);
    return response.data;
  },
};
