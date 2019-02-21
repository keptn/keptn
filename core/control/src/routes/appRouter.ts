import express = require('express');

const router = express.Router();

const GitHub = require('github-api');
const YAML = require('yamljs');

// Basic authentication
const gh = new GitHub({
  username: 'johannes-b',
  password: '**',
  auth: 'basic',
});

router.post('/', async (request: express.Request, response: express.Response) => {

  const payload = {
    data : {
      application: 'sockshop1',
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

  const gitHubOrgName = 'keptn-test';

  await createRepository(gitHubOrgName, payload); 

  await initialCommit(gitHubOrgName, payload);

  await createEmptyBranches(gitHubOrgName, payload);

  await addShipyardToMaster(gitHubOrgName, payload);

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

async function createRepository(gitHubOrgName, payload) {
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

async function initialCommit(gitHubOrgName, payload) {
  try {
    const codeRepo = await gh.getRepo(gitHubOrgName, payload.data.application);
    
    await codeRepo.writeFile('master',
                            'README.md',
                            `# keptn takes care of your ${payload.data.application}`,
                            'Initial commit', { encode: true });
  } catch (e) {
    console.log('Initial commit failed.');
    console.log(e.message);
  }
}

async function createEmptyBranches(gitHubOrgName, payload) {
  try {
    const codeRepo = await gh.getRepo(gitHubOrgName, payload.data.application);

    payload.data.stages.forEach(async stage =>
      await codeRepo.createBranch('master', stage.name),
    );
  } catch (e) {
    console.log('Creating branches failed.');
    console.log(e.message);
  }
}

async function addShipyardToMaster(gitHubOrgName, payload) {
  try {
    const codeRepo = await gh.getRepo(gitHubOrgName, payload.data.application);
    await codeRepo.writeFile('master',
                              'shipyard.yml',
                              YAML.stringify(payload.data),
                              'Added shipyard containing the definition of each stage',
                              { encode: true });
  } catch (e) {
    console.log('Adding shipyard to master failed.');
    console.log(e.message);
  }
}

export = router;
