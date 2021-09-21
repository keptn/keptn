import BasePage from '../support/pageobjects/BasePage'
import SecretsPage from '../support/pageobjects/SecretsPage'

describe('Keptn Secrets adding deleting test', () => {
    it('The test goes over the pages and does validation', () => {
        const basePage = new BasePage()
        const secretsPage = new SecretsPage()
        const SECRET_NAME = "testsecretname";
        const SECRET_KEY = "secretkey";
        const SECRET_VALUE = "secretvalue!@#$%^&*(!@#$%^&*()";
        const DYNATRACE_PROJECT = 'dynatrace';
         
        cy.visit('/')
        //basePage.login('claus.keptn-dev@ruxitlabs.com', 'Labpass12345')
        basePage.selectProject(DYNATRACE_PROJECT);

        basePage.goToUniformPage().goToSecretsPage()
        secretsPage.addSecret(SECRET_NAME, SECRET_KEY, SECRET_VALUE)
        secretsPage.deleteSecret(SECRET_NAME)

    })
})