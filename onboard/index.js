#!/usr/bin/env node
const readline = require('readline');
const chalk = require('chalk');

const creds = require('./creds.json');

const GitHub = require('github-api');


const jenkinsApi = require('jenkins');

const program = require('commander');

const { base64encode, base64decode } = require('nodejs-base64');

const utils = require('./lib/utils.js');

const YAML = require('yamljs');

const decamelize = require('decamelize');

const mustache = require('mustache');

const uuidv1 = require('uuid/v1');

const camelCase = require('camelcase');

const STAGES = ['dev', 'staging', 'production'];

// basic auth
const gh = new GitHub({
   username: creds.githubUserName,
   password: creds.githubPersonalAccessToken,
   auth: 'basic'
});
const gitHubUser = gh.getUser();

let jenkins;
let codeRepos;

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  });

async function createJenkinsJobs(gitHubOrg) {
    let jobConfig = await utils.readFileContent('./job-template/config.xml');
    jobConfig = jobConfig.replace('GITHUB_ORGANIZATION_PLACEHOLDER', gitHubOrg);
    return new Promise((resolve, reject) => {
        jenkins.job.create(gitHubOrg, jobConfig, function(err) {
            // if (err) console.log(err);
            resolve();
          });
    });
}

async function setupHelmChart(configRepo, branchName, gitHubOrgName) {
    let createBlueGreen = branchName === 'staging' || branchName === 'production';
    console.log('setting up config repo helm chart');

    // get the directory structure of the repository containing the helm chart
    let templateTree= await getHelmChartTree();
    
    // download the values.yaml -> check if every app in the organization has an entry there
    const valuesYaml = await configRepo.getContents(branchName, 'helm-chart/values.yaml');
    let valuesYamlObj = YAML.parse(base64decode(valuesYaml.data.content));

    const chartYaml = await configRepo.getContents(branchName, 'helm-chart/Chart.yaml');
    let chartYamlObj = YAML.parse(base64decode(chartYaml.data.content));
    const chartName = chartYamlObj.name;  

    let istioIngressGatewayService;
    if (createBlueGreen) {
        // set up istio gateway for the organization (=app)
        // TODO read from github repo
        let gatewaySpec = await utils.readFileContent('istio-manifests/gateway.tpl');
        gatewaySpec = mustache.render(gatewaySpec, { gitHubOrg: gitHubOrgName });
        await configRepo.writeFile(branchName, `helm-chart/templates/istio-gateway.yaml`, gatewaySpec, `[keptn-onboard]: added istio gateway`, {encode: true});

        // get the istio ingress-gateway IP to configure VirtualService definitions
        istioIngressGatewayService = await utils.getK8sServiceUrl('istio-ingressgateway', 'istio-system');
    }

    let blueGreenValues = {};

    const appKeys = [];
    codeRepos.forEach(repo => appKeys.push(camelCase(repo.name)));

    for (let i = 0; i < appKeys.length; i++) {
        let appKey = appKeys[i];
        console.log(appKey);

        let stringifiedYaluesYaml;
        // microservice already defined in helm chart
        if (valuesYamlObj[appKey] !== undefined) {
            stringifiedYaluesYaml = YAML.stringify(valuesYamlObj[appKey], 100);
            if (appKey.indexOf('Green') < 0 && appKey.indexOf('Blue') < 0) {
                // use JSON.parse(JSON.stringify(...)) to avoid having the same reference for each array entry
                blueGreenValues[appKey] = YAML.parse(stringifiedYaluesYaml);
                
                if (Object.keys(blueGreenValues).find(val => val === `${appKey}Blue`)) {
                    // if there already is a Blue/green spec for this service, don't overwrite the existing one
                    continue;
                }
            }
        } else {
            // microservice not onboarded yet - create new (minimal) deployment yaml spec 
            stringifiedYaluesYaml = await onboardNewMicroservice(appKey);
            // refresh the directory structure
            templateTree = await getHelmChartTree();
        }
        if (createBlueGreen) {
            createBlueGreenValues(appKey, stringifiedYaluesYaml);
            await configRepo.writeFile(branchName, `helm-chart/values.yaml`, YAML.stringify(blueGreenValues, 100), `[keptn-onboard]: added blue/green values`, {encode: true});
        } else {
            continue;
        }

        // get the template for the microservice
        const appTemplates = getMicroserviceTemplates(appKey);

        // create blue/green yamls for each deployment/service
        for (let j = 0; j < appTemplates.length; j++) {
            let template = appTemplates[j];
            console.log(template.path);

            let decamelizedAppKey = decamelize(appKey, '-');

            let templateContentB64Enc = await configRepo.getContents(branchName, `helm-chart/templates/${template.path}`);
            let templateContent = base64decode(templateContentB64Enc.data.content);

            if (template.path.indexOf('-service.yaml') > 0) {
                // create istio destination rule
                await createIstioEntry(decamelizedAppKey, appKey);
                continue;
            }

            await createBlueGreenDeployment(appKey, decamelizedAppKey, templateContent, template);
        }
    }

    async function onboardNewMicroservice(appKey) {
        let valuesTemplate = await utils.readFileContent('service-template/values.tpl');
        let stringifiedYaluesYaml = mustache.render(valuesTemplate, { microServiceName: appKey });
        valuesYamlObj[appKey] = YAML.parse(stringifiedYaluesYaml);
        if (!createBlueGreen) {
            await configRepo.writeFile(branchName, `helm-chart/values.yaml`, YAML.stringify(valuesYamlObj, 100), `[keptn-onboard]: added entry for new app in values.yaml`, { encode: true });
            await configRepo.writeFile('master', `helm-chart/values.yaml`, YAML.stringify(valuesYamlObj, 100), `[keptn-onboard]: added entry for new app in values.yaml`, { encode: true });
        }
        else {
            blueGreenValues[appKey] = YAML.parse(stringifiedYaluesYaml);
        }
        let deploymentTemplate = await utils.readFileContent('service-template/deployment.tpl');
        let serviceTemplate = await utils.readFileContent('service-template/service.tpl');
        let cAppNameRegex = new RegExp('SERVICE_PLACEHOLDER_C', 'g');
        let decAppNameRegex = new RegExp('SERVICE_PLACEHOLDER_DEC', 'g');
        deploymentTemplate = deploymentTemplate.replace(cAppNameRegex, appKey);
        deploymentTemplate = deploymentTemplate.replace(decAppNameRegex, decamelize(appKey, '-'));
        serviceTemplate = serviceTemplate.replace(cAppNameRegex, appKey);
        serviceTemplate = serviceTemplate.replace(decAppNameRegex, decamelize(appKey, '-'));
        await configRepo.writeFile(branchName, `helm-chart/templates/${appKey}-deployment.yaml`, deploymentTemplate, `[keptn-onboard]: added deployment yaml template for new app: ${appKey}`, { encode: true });
        await configRepo.writeFile(branchName, `helm-chart/templates/${appKey}-service.yaml`, serviceTemplate, `[keptn-onboard]: added service yaml template for new app: ${appKey}`, { encode: true });
        if (!createBlueGreen) {
            await configRepo.writeFile('master', `helm-chart/templates/${appKey}-deployment.yaml`, deploymentTemplate, `[keptn-onboard]: added deployment yaml template for new app: ${appKey}`, { encode: true });
            await configRepo.writeFile('master', `helm-chart/templates/${appKey}-service.yaml`, serviceTemplate, `[keptn-onboard]: added service yaml template for new app: ${appKey}`, { encode: true });
        }
        return stringifiedYaluesYaml;
    }

    async function getHelmChartTree() {
        let branch = await configRepo.getBranch(branchName);
        let tree = await configRepo.getTree(branch.data.commit.sha);
        //get the contents of helm-chart/templates
        let helmChartTree = await configRepo.getTree(tree.data.tree.filter(item => item.path === 'helm-chart')[0].sha);
        let templateTree = await configRepo.getTree(helmChartTree.data.tree.filter(item => item.path === 'templates')[0].sha);
        return templateTree;
    }

    async function createBlueGreenDeployment(appKey, decamelizedAppKey, templateContent, template) {
        let replaceRegex = new RegExp(appKey, 'g');
        let tmpRegex = new RegExp('selector-' + decamelizedAppKey, 'g');
        let decamelizedAppNameRegex = new RegExp(decamelizedAppKey, 'g');
        let templateContentBlue = templateContent.replace(replaceRegex, `${appKey}Blue`);
        let tmpString = uuidv1();
        templateContentBlue = templateContentBlue.replace(tmpRegex, tmpString);
        templateContentBlue = templateContentBlue.replace(decamelizedAppNameRegex, `${decamelizedAppKey}-blue`);
        templateContentBlue = templateContentBlue.replace(new RegExp(tmpString, 'g'), 'selector-' + decamelizedAppKey);
        let templateContentGreen = templateContent.replace(replaceRegex, `${appKey}Green`);
        templateContentGreen = templateContentGreen.replace(tmpRegex, tmpString);
        templateContentGreen = templateContentGreen.replace(decamelizedAppNameRegex, `${decamelizedAppKey}-green`);
        templateContentGreen = templateContentGreen.replace(new RegExp(tmpString, 'g'), 'selector-' + decamelizedAppKey);
        let templateBluePathName = template.path.replace(replaceRegex, `${appKey}Blue`);
        let templateGreenPathName = template.path.replace(replaceRegex, `${appKey}Green`);
        await configRepo.writeFile(branchName, `helm-chart/templates/${templateBluePathName}`, templateContentBlue, `[keptn-onboard]: added blue version of ${appKey}`, { encode: true });
        await configRepo.writeFile(branchName, `helm-chart/templates/${templateGreenPathName}`, templateContentGreen, `[keptn-onboard]: added green version of ${appKey}`, { encode: true });
        // delete the original template
        await configRepo.deleteFile(branchName, `helm-chart/templates/${template.path}`);
    }

    async function createIstioEntry(decamelizedAppKey, appKey) {
        let destinationRuleTemplate = await utils.readFileContent('istio-manifests/destination_rule.tpl');
        destinationRuleTemplate = mustache.render(destinationRuleTemplate, {
            serviceName: decamelizedAppKey,
            chartName,
            environment: branchName
        });
        await configRepo.writeFile(branchName, `helm-chart/templates/istio-destination-rule-${appKey}.yaml`, destinationRuleTemplate, `[keptn-onboard]: added istio destination rule for ${appKey}`, { encode: true });
        // create istio virtual service
        let virtualServiceTemplate = await utils.readFileContent('istio-manifests/virtual_service.tpl');
        virtualServiceTemplate = mustache.render(virtualServiceTemplate, {
            gitHubOrg: gitHubOrgName,
            serviceName: decamelizedAppKey,
            chartName,
            environment: branchName,
            ingressGatewayIP: istioIngressGatewayService.ip
        });
        await configRepo.writeFile(branchName, `helm-chart/templates/istio-virtual-service-${appKey}.yaml`, virtualServiceTemplate, `[keptn-onboard]: added istio virtual service for ${appKey}`, { encode: true });
    }

    function getMicroserviceTemplates(appKey) {
        return templateTree.data.tree.filter(item => item.path.indexOf(appKey) === 0 &&
            (item.path.indexOf('yml') > -1 || item.path.indexOf('yaml') > -1) &&
            (item.path.indexOf('Blue') < 0) && (item.path.indexOf('Green') < 0));
    }

    function createBlueGreenValues(appKey, stringifiedYaluesYaml) {
        blueGreenValues[`${appKey}Blue`] = YAML.parse(stringifiedYaluesYaml);
        blueGreenValues[`${appKey}Green`] = YAML.parse(YAML.stringify(valuesYamlObj[appKey], 100));
        blueGreenValues[`${appKey}Blue`].image.tag = `${branchName}-stable`;
        if (blueGreenValues[`${appKey}Blue`].service) {
            blueGreenValues[`${appKey}Blue`].service.name = blueGreenValues[`${appKey}Blue`].service.name + '-blue';
        }
        if (blueGreenValues[`${appKey}Green`].service) {
            blueGreenValues[`${appKey}Green`].service.name = blueGreenValues[`${appKey}Green`].service.name + '-green';
        }
    }
}

async function setupConfigRepoBranch(configRepo, branch) {
    try {
        await configRepo.createBranch('master', branch);
        console.log(chalk.green(`${branch} branch created.`));
    } catch (e) {
        console.log(chalk.yellow(`${branch} branch already exists.`));
    }
}

async function setupConfigRepo(gitHubOrgName, eventBrokerUri) {
    gitHubOrg = await gh.getOrganization(gitHubOrgName);
    let repos = await gitHubOrg.getRepos();

    let configRepo = repos.data.filter(repo => repo.name.indexOf('-config') > -1)[0];
    // TODO: create config repo if it does not exist
    configRepo = await gh.getRepo(gitHubOrgName, configRepo.name);
    await configRepo.updateRepository({name: 'keptn-config'});
    configRepo = await gh.getRepo(gitHubOrgName, 'keptn-config');
    

    await setupConfigRepoBranch(configRepo, 'dev');
    await setupConfigRepoBranch(configRepo, 'staging');
    await setupConfigRepoBranch(configRepo, 'production');

    await setupHelmChart(configRepo, 'dev', gitHubOrgName);
    await setupHelmChart(configRepo, 'staging', gitHubOrgName);
    await setupHelmChart(configRepo, 'production', gitHubOrgName);   

    try {
        await configRepo.createHook({
            name: 'web',
            events: ['push'],
            config: {
                url: `http://${eventBrokerUri}/github`,
                content_type: 'json'   
            }
        });
        console.log(chalk.green(`Webhook created: http://${eventBrokerUri}/github`));
    } catch (e) {
        console.log(chalk.yellow('Webhook already exists'));
    }
}

async function setupJenkinsFile(deployRepo, branchName, keptnOperatorService) {
    // get the Jenkinsfile
    let jenkinsFileContent = await deployRepo.getContents(branchName, 'Jenkinsfile');
    jenkinsFileContent = base64decode(jenkinsFileContent.data.content);

    const placeholderRegex = new RegExp('KEPTN_OPERATOR_PLACEHOLDER', 'g');
    jenkinsFileContent = jenkinsFileContent.replace(placeholderRegex, `http://${keptnOperatorService.ip}:${keptnOperatorService.port}/jenkins`);

    // push the changes to the repo
    await deployRepo.writeFile(branchName, 'Jenkinsfile', jenkinsFileContent, `[keptn-onboard]: inserted keptn operator URL: http://${keptnOperatorService.ip}:${keptnOperatorService.port}/jenkins`, {encode: true});
}

async function setupDeployRepo(gitHubOrgName, keptnOperatorService) {
    gitHubOrg = await gh.getOrganization(gitHubOrgName);
    let repos = await gitHubOrg.getRepos();
    if (repos.data.filter(repo => repo.name.toLowerCase() === 'keptn-deploy').length > 0) {
        console.log(chalk.yellow(`'deploy' repo already exists.`));
    } else {
        console.log(`forking 'keptn-deploy' repo from https://github.com/dynatrace-innovationlab/keptn-deploy into ${gitHubOrgName} organization.`);
        await utils.execCmd(`mkdir tmp_repo && cd tmp_repo && git clone -q https://github.com/dynatrace-innovationlab/keptn-deploy && cd keptn-deploy && hub fork --org=${gitHubOrgName} && cd ../.. && rm -rf tmp_repo`);
        console.log(chalk.green(`'keptn-deploy' repo successfully cloned.`));
    }

    // replace the keptn operator placeholder in the pipelines with the deployed operator in the cluster
    let deployRepo = await gh.getRepo(gitHubOrgName, 'keptn-deploy');

    await setupJenkinsFile(deployRepo, 'dev', keptnOperatorService);
    await setupJenkinsFile(deployRepo, 'staging', keptnOperatorService);
    await setupJenkinsFile(deployRepo, 'production', keptnOperatorService);
}

async function setupCodeRepos(gitHubOrgName, eventBrokerUri) {
    gitHubOrg = await gh.getOrganization(gitHubOrgName);
    let repos = await gitHubOrg.getRepos();
    
    codeRepos = repos.data.filter(repo => repo.name.indexOf('-config') < 0 && repo.name.indexOf('deploy') < 0);

    for (let i = 0; i < codeRepos.length; i++) {
        repo = codeRepos[i];
        codeRepo = await gh.getRepo(gitHubOrgName, repo.name);
        try {        
            await codeRepo.createHook({
                name: 'web',
                events: ['pull_request'],
                config: {
                    url: `http://${eventBrokerUri}/github`,
                    content_type: 'json'   
                }
            });
            console.log(chalk.green(`Webhook created: http://${eventBrokerUri}/github`));
        } catch (e) {
            console.log(chalk.yellow('Webhook already exists'));
        }

        // check if version file exists
        let masterBranch = await codeRepo.getBranch('master');
        let tree = await codeRepo.getTree(masterBranch.data.commit.sha);
        let versionFile = tree.data.tree.filter(file => file.path.toLowerCase() === 'version');
        // if no version file exists,  create one with version 0.1.0
        if (versionFile.length === 0) {
            console.log(`Adding version file to repo ${repo.name}`);
            codeRepo.writeFile('master', 'version', '0.1.0', '[keptn-onboard]: added version file', {encode: true});
        }
    }
}

async function main(gitHubOrg) {
    if (gitHubOrg === undefined) {
        gitHubOrg = await utils.userPrompt('GitHub Organization: ');
    }   
    console.log(chalk.yellow(`Onboarding ${gitHubOrg}`));

    let jenkinsUrl = await utils.getK8sServiceUrl('jenkins', 'cicd');
    let jenkinsConnectionString = `http://${creds.jenkinsUser}:${creds.jenkinsPassword}@${jenkinsUrl.ip}:${jenkinsUrl.port}`
    
    jenkins = jenkinsApi({ baseUrl: jenkinsConnectionString });
    console.log (jenkins.baseUrl);
    try {
        createJenkinsJobs(gitHubOrg);
    } catch(e) {
        console.log(e);
    }

    const istioIngressGatewayService = await utils.getK8sServiceUrl('istio-ingressgateway', 'istio-system');
    const eventBrokerUri = `event-broker.keptn.${istioIngressGatewayService.ip}.xip.io`;

    await setupCodeRepos(gitHubOrg, eventBrokerUri);

    await setupConfigRepo(gitHubOrg, eventBrokerUri);
    
    await setupDeployRepo(gitHubOrg, eventBrokerUri);
    
    return;
}

(async () => {
    program
        .option('-o', '--organization <gitHubOrg>')
        .action(async (gitHubOrg) => {
            await main(gitHubOrg);
            process.exit(0);
        })
        .parse(process.argv);
})().catch(e => {
    console.log(e);
    process.exit(-1);
});

