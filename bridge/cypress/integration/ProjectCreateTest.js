/// <reference types="cypress" />
import BasePage from '../support/pageobjects/BasePage'
import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage'
import createprojjson from '../fixtures/create.project.request.body.json'


describe('sample mock test', () => {
    it('test', () => {
        const basePage = new BasePage()
        const newProjectCreatePage = new NewProjectCreatePage()
        const GIT_USERNAME = 'carpe-github-username'
        const PROJECT_NAME = 'test-project-bycypress-001'
        const GIT_REMOTE_URL = 'https://git-repo.com'
        const GIT_TOKEN = 'testtoken'
        cy.fixture('get.project.json').as("initProjectJSON");
        cy.fixture('metadata.json').as("initmetadata")
  
        cy.intercept('GET', 'api/v1/metadata', {fixture: 'metadata.json'}).as("metadataCmpl")
        cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {fixture: 'get.project.json'}).as("initProjects")
        cy.intercept('POST', 'api/controlPlane/v1/project',  {
            statusCode: 200,
            body: ""
        }).as('createProjectUrl')
        cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', {fixture: 'project.sequences.json'})
        cy.intercept('GET', 'api/project/testproject?approval=true&remediation=true', {
            statusCode: 200
        }).as('projectApproval')

        cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
            statusCode: 200
        }).as('disableUpstreamSync')


        cy.visit('/')
        cy.wait('@metadataCmpl')
        basePage.declineAutomaticUpdate()
        basePage.clickCreatNewProjectButton().
        inputProjectName(PROJECT_NAME).
        inputGitUrl(GIT_REMOTE_URL).
        inputGitUsername(GIT_USERNAME).
        inputGitToken(GIT_TOKEN)

        cy.fixture('shipyard.yaml').then(fileContent => {
            cy.get('input[id="shipyard-file-input"]').attachFile({
                fileContent: fileContent.toString(),
                fileName: 'shipyard.yaml'
            });
        });

        newProjectCreatePage.clickCreateProject()

        cy.wait('@createProjectUrl', {timeout: 20000})
        cy.get('@createProjectUrl').its('request.body').should('deep.equal', createprojjson)

    })
})
