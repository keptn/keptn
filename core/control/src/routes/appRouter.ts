import express = require('express');
import { CreateRequest } from '../types/createRequest';
import { FactoryBuilder } from '../factoryBuilders/factoryBuilder';
import { GitHubFactoryBuilder } from '../factoryBuilders/gitHubFactoryBuilder';

const router = express.Router();

router.post('/', async (request: express.Request, response: express.Response) => {

  const payload : CreateRequest = {
    data : {
      application: 'sockshop15',
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

  let factoryBuilder = GitHubFactoryBuilder.getInstance();

  await factoryBuilder.createRepository(gitHubOrgName, payload); 

  await factoryBuilder.initialCommit(gitHubOrgName, payload);

  await factoryBuilder.createBranchesForEachStages(gitHubOrgName, payload);

  await factoryBuilder.addShipyardToMaster(gitHubOrgName, payload);

  await factoryBuilder.setHook(gitHubOrgName, payload);

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

export = router;
