import { Component, Input } from '@angular/core';
import { FlatTreeControl } from '@angular/cdk/tree';
import { DtTreeControl, DtTreeDataSource, DtTreeFlattener } from '@dynatrace/barista-components/core';

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
    this.treeControl = new DtTreeControl<TreeFlatEntry>(
      this.getEntryLevel,
      this.isEntryExpendable,
    );

    this.treeFlattener = new DtTreeFlattener(
      this.treeTransformer,
      this.getEntryLevel,
      this.isEntryExpendable,
      this.getEntryChildren,
    );
    this.treeDataSource = new DtTreeDataSource(
      this.treeControl,
      this.treeFlattener,
    );
  }

  public getGitRepositoryLink(): string {
    if (this.remoteUri) {
      if (this.remoteUri.includes('github.') || this.remoteUri.includes('gitlab.')) {
        return this.remoteUri + '/tree/' + this.stageName + '/' + this.serviceName;
      }
      if (this.remoteUri.includes('bitbucket.')) {
        return this.remoteUri + '/src/' + this.stageName + '/' + this.serviceName;
      }
      if (this.remoteUri.includes('azure.')) {
        return this.remoteUri + '?path=' + this.serviceName + '&version=GB' + this.stageName;
      }
      if (this.remoteUri.includes('git-codecommit.')) {
        const repoParts = this.remoteUri.split('/');
        const region = repoParts.find(part => part.includes('git-codecommit.'))?.split('.')[1];
        const repoName = repoParts[repoParts.length - 1];
        return 'https://' + region + '.console.aws.amazon.com/codesuite/codecommit/repositories/' + repoName + '/browse/refs/heads/' + this.stageName;
      }

      return this.remoteUri;
    }
    return '';
  }

  private getEntryLevel(entry: TreeFlatEntry): number {
    return entry.level;
  }

  private isEntryExpendable(entry: TreeFlatEntry): boolean {
    return entry.expandable;
  }

  private getEntryChildren(entry: TreeEntry): TreeEntry[] {
    if (!entry.children) {
      return [];
    }
    return entry.children;
  }

  private treeTransformer(node: TreeEntry, level: number): TreeFlatEntry {
    const flatNode = new TreeFlatEntry();
    flatNode.fileName = node.fileName;
    flatNode.level = level;
    flatNode.expandable = !!node.children;

    return flatNode;
  }
}

export interface TreeEntry {
  fileName: string;
  children?: TreeEntry[];
}

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
