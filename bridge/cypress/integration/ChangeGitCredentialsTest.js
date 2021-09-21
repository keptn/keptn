/// <reference types="cypress" />

import BasePage from '../support/pageobjects/BasePage'

describe('Changing git credentials', () => {
    it('The test changes git credentials and makes sure they changed successfully', () => {
        const basePage = new BasePage()
        const DYNATRACE_PROJECT = 'dynatrace'
        const GIT_URL = 'https://git-repo.com'
        const GIT_USER = 'test-username'
        const GIT_TOKEN = 'test-token!§$%&/()='

        cy.visit('/')
        basePage.selectProject(DYNATRACE_PROJECT);
        basePage.declineAutomaticUpdate()
        basePage.gotoSettingsPage().
            inputGitUrl(GIT_URL).
            inputGitUsername(GIT_USER).
            inputGitToken(GIT_TOKEN).
            clickSaveChanges().
            waitUntilGitTokenIsSet()
        //  getErrorMessageText().contains('failed to set upstream: Error executing command git remote show origin: exit status 128')

    })
})
