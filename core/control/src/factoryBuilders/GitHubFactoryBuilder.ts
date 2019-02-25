import {FactoryBuilder} from './FactoryBuilder';
import { CreateRequest } from '../types/createRequest';
import { Utils } from '../lib/utils';

const GitHub = require('github-api');
const Mustache = require('mustache');
const YAML = require('yamljs');

// Util class
const utils = new Utils();

// Basic authentication
const gh = new GitHub({
  username: '**',
  password: '**',
  auth: 'basic',
});

export class GitHubFactoryBuilder {

  private static instance: GitHubFactoryBuilder;

  private constructor() {
  }

  static getInstance() {
    if (GitHubFactoryBuilder.instance === undefined) {
      GitHubFactoryBuilder.instance = new GitHubFactoryBuilder();
    }
    return GitHubFactoryBuilder.instance;
  }

  async createRepository(gitHubOrgName : string, payload : CreateRequest) : Promise<any> {
    const repository = {
      name : payload.data.application,
    };
  
    try {
      const organization = await gh.getOrganization(gitHubOrgName);
      await organization.createRepo(repository);
    } catch (e) {
      if (e.response.statusText == 'Not Found') {
        console.log(`[keptn] Could not find organziation ${gitHubOrgName}.`);
        console.log(e.message);
      } else if (e.response.statusText == 'Unprocessable Entity'){
        console.log(`[keptn] Repository ${payload.data.application} already available.`);
        console.log(e.message);
      }
    }  
  }

  async setHook(gitHubOrgName : string, payload : CreateRequest) : Promise<any> {
    try {
      const repo = await gh.getRepo(gitHubOrgName, payload.data.application);
  
      //const istioIngressGatewayService = await utils.getK8sServiceUrl('istio-ingressgateway', 'istio-system');
      //const eventBrokerUri = `event-broker.keptn.${istioIngressGatewayService.ip}.xip.io`;
      const eventBrokerUri = 'need-to-be-set';
  
      await repo.createHook({
          name: 'web',
          events: ['push'],
          config: {
              url: `http://${eventBrokerUri}/github`,
              content_type: 'json'   
          }
      });
      console.log(`Webhook created: http://${eventBrokerUri}/github`);
    } catch (e) {
      console.log('Setting hook failed.');
      console.log(e.message);
    }
  }
  
  async initialCommit(gitHubOrgName : string, payload : CreateRequest) : Promise<any> {
    try {
      const repo = await gh.getRepo(gitHubOrgName, payload.data.application);
      
      await repo.writeFile('master',
                           'README.md',
                           `# keptn takes care of your ${payload.data.application}`,
                           '[keptn]: Initial commit', { encode: true });
    } catch (e) {
      console.log('[keptn] Initial commit failed.');
      console.log(e.message);
    }
  }
  
  async createBranchesForEachStages(gitHubOrgName : string, payload : CreateRequest) : Promise<any> {
  
    const chart = {
      apiVersion: 'v1',
      description: 'A Helm chart for Kubernetes',
      name: 'mean-k8s',
      version: '0.1.0'
    };
    
    try {
      const repo = await gh.getRepo(gitHubOrgName, payload.data.application);
  
      payload.data.stages.forEach(async stage => {
        await repo.createBranch('master', stage.name);
  
        await repo.writeFile(stage.name,
                             'helm-chart/Chart.yml',
                             YAML.stringify(chart),
                             '[keptn]: Added helm-chart Chart.yml file.',
                             { encode: true });
  
        await repo.writeFile(stage.name,
                             'helm-chart/values.yml',
                             '',
                             '[keptn]: Added helm-chart values.yml file.',
                             { encode: true });
  
        if(stage.deployment_strategy === 'blue_green_service' ) {
          // Add istio gateway to stage
          let gatewaySpec = await utils.readFileContent('keptn/keptn/core/control/src/routes/istio-manifests/gateway.tpl');
          gatewaySpec = Mustache.render(gatewaySpec, { application: payload.data.application, stage: stage.name});
          
          await repo.writeFile(stage.name,
                               'helm-chart/templates/istio-gateway.yaml',
                               gatewaySpec, 
                               '[keptn]: Added istio gateway.',
                               {encode: true});
        }
  
      });
    } catch (e) {
      console.log('[keptn] Creating branches failed.');
      console.log(e.message);
    }
  }
  
  async addShipyardToMaster(gitHubOrgName : string, payload : CreateRequest) : Promise<any> {
    try {
      const repo = await gh.getRepo(gitHubOrgName, payload.data.application);
      await repo.writeFile('master',
                           'shipyard.yml',
                           YAML.stringify(payload.data),
                           '[keptn]: Added shipyard containing the definition of each stage',
                           { encode: true });
    } catch (e) {
      console.log('[keptn] Adding shipyard to master failed.');
      console.log(e.message);
    }
  }

}