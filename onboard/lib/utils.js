const fs = require('fs');

const exec = require('child_process').exec;

const utils = {
    readFileContent: async function(filePath) {
        return String(fs.readFileSync(filePath));
    },
    execCmd: async function(cmd) {
        return new Promise((resolve, reject) => {
            exec(cmd, function(err, stdout, stderr) {
                if (err) {
                    reject(stderr);
                }
                resolve(stdout);
            });   
        });
    },
    getK8sServiceUrl: async function(serviceName, namespace) {
        let services = await this.execCmd(`kubectl get svc -n ${namespace} -o json`);
        return new Promise((resolve, reject) => {
            services = JSON.parse(services);
            let service = services.items.filter(svc => svc.metadata.name === serviceName)[0];
            const httpPorts = service.spec.ports.filter(port => port.name === 'http');
            let port;
            if (httpPorts.length > 0) {
                port = httpPorts[0].port
            } 
            resolve({
                ip: service.status.loadBalancer.ingress[0].ip,
                port
            });       
        });
    },
    userPrompt: async function (question) {
        return new Promise((resolve, reject) => {
            rl.question(chalk.green(question), (answer) => {
                rl.close();
                resolve(answer);
            });
        });
    }
}

module.exports = utils;