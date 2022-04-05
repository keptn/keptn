import { ErrorPage } from '../support/pageobjects/ErrorPage';
import { LogoutPage } from '../support/pageobjects/LogoutPage';
import { interceptDashboard } from '../support/intercept';

describe('Test error pages', () => {
  const errorPage = new ErrorPage();

  it('should show internal error if status is 500', () => {
    errorPage.visit('500').isInternalError();
  });

  it('should show internal error if status is not provided', () => {
    errorPage.visit().isInternalError();
  });

  it('should show insufficient permission if status is 403', () => {
    errorPage.visit('403').isInsufficientPermissionError();
  });

  it('should show internal error if status is not a valid number', () => {
    errorPage.visit('asdf').isInternalError();
  });

  it('should show internal error if status is neither 500 nor 403', () => {
    errorPage.visit('400').isInternalError();
  });

  it('should redirect to dashboard', () => {
    interceptDashboard();
    errorPage.visit().clickLocation();
    cy.location('pathname').should('eq', '/dashboard');
  });
});

describe('Test logout page', () => {
  const logoutPage = new LogoutPage();
  it('should be logout page and redirect to dashboard', () => {
    interceptDashboard();
    logoutPage.visit().assertHeaderText('You haven been logged out').clickLocation();
    cy.location('pathname').should('eq', '/dashboard');
  });
});
