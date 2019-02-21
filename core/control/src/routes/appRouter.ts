import express = require('express');

const router = express.Router();

const GitHub = require('github-api');
const YAML = require('yamljs');

// Basic authentication
const gh = new GitHub({
  username: 'johannes-b',
  password: 'd74b5346409bdebe31ac4e9011602c17da62d7e5',
  auth: 'basic',
});

router.post('/', async (request: express.Request, response: express.Response) => {

  const payload = {
    data : {
      application: 'sockshop',
      stages: [
        {
          name: 'dev',
          next: 'staging',
        },
        {
          name: 'staging',
          next: 'production',
        },
        {
          name: 'production',
        },
      ],
    },
  };

  const repository = {
    name : payload.data.application,
  };

  const gitHubOrgName = 'keptn-test';

  // Create repository
  try {
    const organization = await gh.getOrganization(gitHubOrgName);
    //await organization.createRepo(repository);
  } catch (e) {
    console.log('Creating repository failed.');
    console.log(e.message);
  }

  // Initial commit
  try {
    const codeRepo = await gh.getRepo(gitHubOrgName, 'sockshop');
    /*
    await codeRepo.writeFile('master',
                             'README.md',
                             `# keptn takes care of your ${payload.data.application}`,
                             'Initial commit', { encode: true });
    */
  } catch (e) {
    console.log('Initial commit failed.');
    console.log(e.message);
  }

  // Create branches
  try {
    const codeRepo = await gh.getRepo(gitHubOrgName, 'sockshop');

    /*
    payload.data.stages.forEach(async stage =>
      await codeRepo.createBranch('master', stage.name),
    );
    */
  } catch (e) {
    console.log('Creating branches failed.');
    console.log(e.message);
  }

  // Add shipyard to master
  try {
    const codeRepo = await gh.getRepo(gitHubOrgName, 'sockshop');
    /*
    await codeRepo.writeFile('master',
                             'shipyard.yml',
                             YAML.stringify(payload.data),
                             'Added shipyard containing the definition of each stage',
                             { encode: true });
    */
  } catch (e) {
    console.log('Initial commit failed.');
    console.log(e.message);
  }

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
