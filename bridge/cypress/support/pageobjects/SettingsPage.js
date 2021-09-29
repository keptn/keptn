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

    getErrorMessageText(){
        return cy.get('.small')
    }

    clickDeleteProjectButton(){
        cy.get('span.dt-button-label').contains('Delete this project').click()
        return this
    }

    typeProjectNameToDelete(projectName){
        const projectInputLoc = 'input[placeholder=proj_pattern]'
        cy.get(projectInputLoc.replace("proj_pattern", projectName)).click().type(projectName)
        return this
    }

    submitDelete(){
        cy.get('span.dt-button-label').contains('I understand the consequences, delete this project').click()
    }


}

export default SettingsPage
