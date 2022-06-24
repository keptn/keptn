import { KtbUserComponent } from './ktb-user.component';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { EndSessionData } from '../../../../../shared/interfaces/end-session-data';
import { KtbAppHeaderModule } from '../ktb-app-header.module';

describe('ktbUserComponentTest', () => {
  let component: KtbUserComponent;
  let fixture: ComponentFixture<KtbUserComponent>;
  let httpMock: HttpTestingController;
  const locationAssignMock = mockLocation();

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbAppHeaderModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbUserComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    locationAssignMock.mockClear();
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeDefined();
  });

  it('should send POST to logout', () => {
    const submitForm = { target: { submit: (): void => {} } };
    const submitSpy = jest.spyOn(submitForm.target, 'submit');
    component.logout(submitForm);
    httpMock.expectOne('./oauth/logout').flush({
      id_token_hint: '',
      end_session_endpoint: '',
      post_logout_redirect_uri: '',
      state: '',
    } as EndSessionData);
    expect(submitSpy).toHaveBeenCalled();
  });

  it('should redirect to root if no data for logout is returned', () => {
    const submitForm = { target: { submit: (): void => {} } };
    component.logout(submitForm);
    httpMock.expectOne('./oauth/logout').flush(null);
    expect(locationAssignMock).toBeCalledWith('http://localhost/logoutsession');
  });
});

function mockLocation(): jest.Mock<unknown, unknown[]> {
  const locationAssignMock = jest.fn();
  /* eslint-disable @typescript-eslint/ban-ts-comment */
  // @ts-ignore
  delete window.location;
  // @ts-ignore
  window.location = { assign: locationAssignMock };
  /* eslint-enable @typescript-eslint/ban-ts-comment */
  return locationAssignMock;
}
