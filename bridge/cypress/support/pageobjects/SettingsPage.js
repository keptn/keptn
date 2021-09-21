/// <reference types="cypress" />

class SettingsPage{

    inputGitUrl(GIT_URL){
        cy.get('input[formcontrolname="gitUrl"]').type(GIT_URL)
        return this
    }

    inputGitUsername(GIT_USERNAME){
        cy.get('input[formcontrolname="gitUser"]').type(GIT_USERNAME)
        return this
    }

    inputGitToken(GIT_TOKEN){
        cy.get('input[formcontrolname="gitToken"]').type(GIT_TOKEN)
        return this
    }

    clickSaveChanges(){
        cy.get('.dt-button-primary > span.dt-button-label').contains('Save changes').click()
        return this
    }

    waitUntilGitTokenIsSet(){
        cy.get('dt-loading-spinner', { timeout: 3000 }).should('be.visible')
        cy.get('dt-loading-spinner', { timeout: 10000 }).should('not.be.visible')
        return this
    }

    getErrorMessageText(){
        return cy.get('.small')
    }

}

export default SettingsPage
