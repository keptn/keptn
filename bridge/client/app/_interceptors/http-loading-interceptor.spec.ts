import {HttpLoadingInterceptor} from './http-loading-interceptor';
import {HttpStateService} from "../_services/http-state.service";

describe('HttpLoadingInterceptor', () => {
  it('should create an instance', () => {
    expect(new HttpLoadingInterceptor(new HttpStateService())).toBeTruthy();
  });
});
