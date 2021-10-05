import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEditServiceFileListComponent } from './ktb-edit-service-file-list.component';
import { AppModule } from '../app.module';
import { By } from '@angular/platform-browser';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbEditServiceFileListComponent', () => {
  let component: KtbEditServiceFileListComponent;
  let fixture: ComponentFixture<KtbEditServiceFileListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEditServiceFileListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show the files for one stage', () => {
    // given
    const tree = [
      {
        fileName: 'helm',
        children: [
          {
            fileName: 'carts',
            children: [
              {
                fileName: 'templates',
                children: [
                  {
                    fileName: 'deployment.yaml',
                  },
                  {
                    fileName: 'service.yaml',
                  },
                ],
              },
              {
                fileName: 'Chart.yaml',
              },
              {
                fileName: 'values.yaml',
              },
            ],
          },
        ],
      },
      {
        fileName: 'metadata.yaml',
      },
    ];
    component.stageName = 'dev';
    component.serviceName = 'carts';
    component.remoteUri = 'https://github.com/keptn/sockshop-upstream';
    component.treeData = tree;

    // when
    fixture.detectChanges();

    // then
    const stageElem = fixture.nativeElement.querySelector('div > div.bold');

    for (let i = 0; i <= 4; i++) {
      const toggles = fixture.debugElement.queryAll(By.css('.dt-tree-table-toggle:enabled'));
      toggles[toggles.length - 1].nativeElement.click();
      fixture.detectChanges();
    }

    const tableRowElements = fixture.nativeElement.querySelectorAll('.dt-tree-table-row');

    expect(component.treeDataSource.data).toEqual(tree);
    expect(stageElem).toBeTruthy();
    expect(stageElem.textContent).toEqual('dev');

    expect(tableRowElements).toBeTruthy();
    expect(tableRowElements.length).toEqual(8);
    expect(tableRowElements[0].textContent.trim()).toEqual('helm');
    expect(tableRowElements[0].getAttribute('aria-expanded')).toEqual('true');
    expect(tableRowElements[1].textContent.trim()).toEqual('carts');
    expect(tableRowElements[1].getAttribute('aria-expanded')).toEqual('true');
    expect(tableRowElements[2].textContent.trim()).toEqual('templates');
    expect(tableRowElements[2].getAttribute('aria-expanded')).toEqual('true');
    expect(tableRowElements[3].textContent.trim()).toEqual('deployment.yaml');
    expect(tableRowElements[4].textContent.trim()).toEqual('service.yaml');
    expect(tableRowElements[5].textContent.trim()).toEqual('Chart.yaml');
    expect(tableRowElements[6].textContent.trim()).toEqual('values.yaml');
    expect(tableRowElements[7].textContent.trim()).toEqual('metadata.yaml');
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
    expect(link).toEqual(
      'https://eu-central-1.console.aws.amazon.com/codesuite/codecommit/repositories/sockshop-upstream/browse/refs/heads/dev'
    );
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
