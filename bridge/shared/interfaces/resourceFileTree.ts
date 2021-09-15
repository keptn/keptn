export interface FileTree {
  stageName: string;
  tree: TreeEntry[];
}

export interface TreeEntry {
  fileName: string;
  children?: TreeEntry[];
}
