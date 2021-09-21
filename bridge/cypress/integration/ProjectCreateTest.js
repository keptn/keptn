/// <reference types="cypress" />

import BasePage from '../support/pageobjects/BasePage'


describe('create new project', () => {
    it('The test checks the UI of create new project ', () => {
        const basePage = new BasePage()
        const DYNATRACE_PROJECT = 'dynatrace'
        const GIT_USERNAME = 'carpe-github-username'
        const PROJECT_NAME = 'test-project-bycypress-001'
        const GIT_REMOTE_URL = 'https://git-repo.com'
        const GIT_TOKEN = 'testtoken§$%&/()='

        cy.visit('/')
        basePage.selectProject(DYNATRACE_PROJECT)
        basePage.declineAutomaticUpdate()

        basePage.clickMainHeaderKeptn()
        basePage.clickCreatNewProjectButton().
            inputProjectName(PROJECT_NAME).
            inputGitUrl(GIT_REMOTE_URL).
            inputGitUsername(GIT_USERNAME).
            inputGitToken(GIT_TOKEN)
    })
})
