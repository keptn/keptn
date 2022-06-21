import { Component, Input } from '@angular/core';
import { FlatTreeControl } from '@angular/cdk/tree';
import { DtTreeControl, DtTreeDataSource, DtTreeFlattener } from '@dynatrace/barista-components/core';
import { TreeEntry } from '../../../../shared/interfaces/resourceFileTree';

export class TreeFlatEntry {
  fileName: string;
  level: number;
  expandable: boolean;

  constructor() {
    this.fileName = '';
    this.level = -1;
    this.expandable = false;
  }
}

@Component({
  selector: 'ktb-edit-service-file-list',
  templateUrl: './ktb-edit-service-file-list.component.html',
  styleUrls: [],
})
export class KtbEditServiceFileListComponent {
  public treeControl: FlatTreeControl<TreeFlatEntry>;
  private readonly treeFlattener: DtTreeFlattener<TreeEntry, TreeFlatEntry>;
  public treeDataSource: DtTreeDataSource<TreeEntry, TreeFlatEntry>;

  @Input() public stageName = '';
  @Input() public serviceName = '';
  @Input() public remoteUri: string | undefined = '';

  @Input()
  set treeData(treeData: TreeEntry[]) {
    this.treeDataSource.data = treeData;
  }

  constructor() {
    this.treeControl = new DtTreeControl<TreeFlatEntry>(this.getLevel, this.isExpandable);

    this.treeFlattener = new DtTreeFlattener(this.treeTransformer, this.getLevel, this.isExpandable, this.getChildren);
    this.treeDataSource = new DtTreeDataSource(this.treeControl, this.treeFlattener);
  }

  public getGitRepositoryLink(): string {
    let uri = '';
    if (this.remoteUri) {
      uri = this.remoteUri;

      // Remove .git from the end of the URL
      if (uri.endsWith('.git')) {
        uri = uri.slice(0, -4);
      }

      if (uri.includes('github.com/') || uri.includes('gitlab.com/')) {
        return this.getGithubGitlabUrl(uri);
      }
      if (uri.includes('bitbucket.org/')) {
        return this.getBitbucketUrl(uri);
      }
      if (uri.includes('dev.azure.com/')) {
        return this.getAzureUrl(uri);
      }
      if (uri.includes('git-codecommit.')) {
        return this.getCodeCommitUrl(uri);
      }
    }
    return uri;
  }

  private getGithubGitlabUrl(uri: string): string {
    return uri + '/tree/' + this.stageName + '/' + this.serviceName;
  }

  private getBitbucketUrl(uri: string): string {
    uri = uri.replace(/https:\/\/(.*)@/, 'https://');
    return uri + '/src/' + this.stageName + '/' + this.serviceName;
  }

  private getAzureUrl(uri: string): string {
    uri = uri.replace(/https:\/\/.*@dev.azure.com\//, 'https://dev.azure.com/');
    return uri;
  }

  private getCodeCommitUrl(uri: string): string {
    const repoParts = uri.split('/');
    const region = repoParts.find((part) => part.includes('git-codecommit.'))?.split('.')[1];
    const repoName = repoParts[repoParts.length - 1];
    return (
      'https://' +
      region +
      '.console.aws.amazon.com/codesuite/codecommit/repositories/' +
      repoName +
      '/browse/refs/heads/' +
      this.stageName
    );
  }

  private getLevel(entry: TreeFlatEntry): number {
    return entry.level;
  }

  private isExpandable(entry: TreeFlatEntry): boolean {
    return entry.expandable;
  }

  private getChildren(entry: TreeEntry): TreeEntry[] {
    return entry.children || [];
  }

  private treeTransformer(node: TreeEntry, level: number): TreeFlatEntry {
    const flatNode = new TreeFlatEntry();
    flatNode.fileName = node.fileName;
    flatNode.level = level;
    flatNode.expandable = !!node.children;

    return flatNode;
  }
}
