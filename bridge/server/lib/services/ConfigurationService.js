/* eslint-disable no-underscore-dangle */
const axios = require('axios');

class ConfigurationService {

  constructor(endpoint) {
    this.api = endpoint;
  }

  async getProjects() {
    const url = `${this.api}/project?pageSize=50`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.projects) {
      return data.projects;
    }
    return [];
  }

  async getProjectResources(projectName) {
    const url = `${this.api}/project/${projectName}/resource?pageSize=50`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.resources) {
      return data.resources;
    }
    return [];
  }

  async getStages(projectName) {
    const url = `${this.api}/project/${projectName}/stage?pageSize=50`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.stages) {
      return data.stages;
    }
    return [];
  }

  async getServices(projectName, stageName) {
    const url = `${this.api}/project/${projectName}/stage/${stageName}/service?pageSize=50`;
    const result = await axios.get(url);
    const { data } = result;
    if (data.services) {
      return data.services;
    }
    return [];
  }

}

module.exports = ConfigurationService;
