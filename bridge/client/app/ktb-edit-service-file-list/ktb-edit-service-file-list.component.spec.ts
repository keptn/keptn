import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEditServiceFileListComponent } from './ktb-edit-service-file-list.component';

describe('KtbEditServiceFileListComponent', () => {
  let component: KtbEditServiceFileListComponent;
  let fixture: ComponentFixture<KtbEditServiceFileListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({})
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEditServiceFileListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should get the link for github for a given stage', () => {
    // given
    component.stageName = 'dev';
    component.remoteUri = 'https://github.com/keptn/sockshop-upstream';
    component.serviceName = 'carts';

    // when
    const link = component.getGitRepositoryLink();

    // then
    expect(link).toEqual('https://github.com/keptn/sockshop-upstream/tree/dev/carts');
  });

  it('should get the link for bitbucket for a stage', () => {
    // given
    component.stageName = 'dev';
    component.remoteUri = 'https://bitbucket.org/keptn/sockshop-upstream';
    component.serviceName = 'carts';

    // when
    const link = component.getGitRepositoryLink();

    // then
    expect(link).toEqual('https://bitbucket.org/keptn/sockshop-upstream/src/dev/carts');
  });

  it('should get the link for azure for a stage', () => {
    // given
    component.stageName = 'dev';
    component.remoteUri = 'https://dev.azure.com/keptn/_git/sockshop-upstream';
    component.serviceName = 'carts';

    // when
    const link = component.getGitRepositoryLink();

    // then
    expect(link).toEqual('https://dev.azure.com/keptn/_git/sockshop-upstream?path=carts&version=GBdev');
  });

  it('should get the link for codeCommit for a stage', () => {
    // given
    component.stageName = 'dev';
    component.remoteUri = 'https://git-codecommit.eu-central-1.amazonaws.com/v1/repos/sockshop-upstream';
    component.serviceName = 'carts';

    // when
    const link = component.getGitRepositoryLink();

    // then
    expect(link).toEqual('https://eu-central-1.console.aws.amazon.com/codesuite/codecommit/repositories/sockshop-upstream/browse/refs/heads/dev');
  });

  it('should return the repository url when not github, bitbucket, azure or codeCommit', () => {
    // given
    component.stageName = 'dev';
    component.remoteUri = 'https://some-other-git-provider.com/keptn/keptn-upstream';
    component.serviceName = 'carts';

    // when
    const link = component.getGitRepositoryLink();

    // then
    expect(link).toEqual('https://some-other-git-provider.com/keptn/keptn-upstream');
  });

  it('should return an empty string if no remote URI is set', () => {
    // given
    component.stageName = 'dev';
    component.remoteUri = '';
    component.serviceName = 'carts';

    // when
    const link = component.getGitRepositoryLink();

    // then
    expect(link).toEqual('');
  });
});
