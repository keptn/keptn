const GitHub = require('github-api');
const readline = require('readline');
const chalk = require('chalk');

class ShipYardOnboarder {
    constructor(gitHubUserName, gitHubPersonalAccessToken, gitHubOrganizationName) {
        this.gh = new GitHub({
            username: gitHubUserName,
            password: gitHubPersonalAccessToken,
            auth: 'basic'
        });
        this.gitHubOrganizationName = gitHubOrganizationName;
    }

    async initialize() {
        this.gitHubOrg = await gh.getOrganization(this.gitHubOrganizationName);
        this.configRepo = this.gh.getRepo(this.gitHubOrganizationName, 'keptn-config');
    }

    static async create(gitHubUserName, gitHubPersonalAccessToken, gitHubOrganizationName) {
        const so = new ShipYardOnboarder(gitHubUserName, gitHubPersonalAccessToken, gitHubOrganizationName);
        await so.initialize();
        return so;
    }

    async onboardShipyard() {
        if (await doesShipyardExist()) {
            // list the current content of the shipyard
            let shipyardDirContent = this.configRepo.getTree(this.shipyardDir.sha);
            if (shipyardDirContent.data.tree.length > 0) {
                console.log(chanlk.yellow(`The following shipyard files have been found in ${this.gitHubOrganizationName}/keptn-config:`));
                shipyardDirContent.data.tree.forEach(item => {
                    console.log(chalk.blue(item.path));
                });
            }
            await userP
        }
    }

    async doesShipyardExist() {
        await this.refreshTree();
        this.shipyardDir = this.tree.data.tree.find(item => item.path === 'shipyard');
        return shipyardDir !== undefined;
    }

    async refreshTree() {
        this.masterBranch = await this.configRepo.getBranch('master');
        this.tree = await this.configRepo.getTree(this.masterBranch.data.commit.sha);
    }
}

module.exports = ShipYardOnboarder;