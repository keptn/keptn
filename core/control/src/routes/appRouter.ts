import express = require('express');

import { CreateRequest } from '../types/createRequest';
import { Utils } from '../lib/utils';

const GitHub = require('github-api');
const YAML = require('yamljs');

const router = express.Router();
const utils = new Utils();

// Basic authentication
const gh = new GitHub({
  username: '**',
  password: '**',
  auth: 'basic',
});

router.post('/', async (request: express.Request, response: express.Response) => {

  const payload : CreateRequest = {
    data : {
      application: 'sockshop8',
      stages: [
        {
          name: 'dev',
          deployment_strategy: 'direct',
        },
        {
          name: 'staging',
          deployment_strategy: 'blue_green_service',
        },
        {
          name: 'production',
          deployment_strategy: 'blue_green_service',
        },
      ],
    },
  };

  const gitHubOrgName = 'keptn-test';

  await createRepository(gitHubOrgName, payload); 

  await initialCommit(gitHubOrgName, payload);

  await createBranchesForEachStages(gitHubOrgName, payload);

  await addShipyardToMaster(gitHubOrgName, payload);

  await setHook(gitHubOrgName, payload);

  const result = {
    result: 'success',
  };

  response.send(result);
});

router.get('/', (request: express.Request, response: express.Response) => {

  const result = {
    result: 'success',
  };

  response.send(result);
});

router.delete('/', async (request: express.Request, response: express.Response) => {

  const result = {
    result: 'success',
  };

  response.send(result);
});

async function setHook(gitHubOrgName : string, payload : CreateRequest) {
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

async function createRepository(gitHubOrgName : string, payload : CreateRequest) {
  const repository = {
    name : payload.data.application,
  };

  try {
    const organization = await gh.getOrganization(gitHubOrgName);
    await organization.createRepo(repository);
  } catch (e) {
    console.log('Creating repository failed.');
    console.log(e.message);
  }  
}

async function initialCommit(gitHubOrgName : string, payload : CreateRequest) {
  try {
    const repo = await gh.getRepo(gitHubOrgName, payload.data.application);
    
    await repo.writeFile('master',
                         'README.md',
                         `# keptn takes care of your ${payload.data.application}`,
                         '[keptn]: Initial commit', { encode: true });
  } catch (e) {
    console.log('Initial commit failed.');
    console.log(e.message);
  }
}

async function createBranchesForEachStages(gitHubOrgName : string, payload : CreateRequest) {

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
                           `helm-chart/values.yml`,
                           '',
                           '[keptn]: Added helm-chart values.yml file.',
                           { encode: true });

      if(stage.deployment_strategy === 'blue_green_service' ) {
        let gatewaySpec = await utils.readFileContent('keptn/keptn/core/control/src/routes/istio-manifests/gateway.tpl');
        //gatewaySpec = mustache.render(gatewaySpec, { gitHubOrg: gitHubOrgName });
        await repo.writeFile(stage.name, 'helm-chart/templates/istio-gateway.yaml', gatewaySpec, `[keptn]: Added istio gateway`, {encode: true});
      }

    });
  } catch (e) {
    console.log('Creating branches failed.');
    console.log(e.message);
  }
}

async function createHelmChart() {
  console.log("create Helm charts");
}

async function addShipyardToMaster(gitHubOrgName : string, payload : CreateRequest) {
  try {
    const repo = await gh.getRepo(gitHubOrgName, payload.data.application);
    await repo.writeFile('master',
                         'shipyard.yml',
                         YAML.stringify(payload.data),
                         '[keptn]: Added shipyard containing the definition of each stage',
                         { encode: true });
  } catch (e) {
    console.log('Adding shipyard to master failed.');
    console.log(e.message);
  }
}

export = router;
