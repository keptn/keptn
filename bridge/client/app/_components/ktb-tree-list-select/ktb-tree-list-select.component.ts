import {
  Component,
  ComponentRef,
  Directive,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  OnInit,
  Output,
  TemplateRef,
} from '@angular/core';
import { DtTreeControl, DtTreeDataSource, DtTreeFlattener } from '@dynatrace/barista-components/core';
import { FlatTreeControl } from '@angular/cdk/tree';
import { ComponentPortal } from '@angular/cdk/portal';
import { OverlayService } from '../../_directives/overlay-service/overlay.service';
import { KtbOverlay } from '../_abstract/ktb-overlay';

export interface SelectTreeNode {
  name: string;
  path?: string;
  keys?: SelectTreeNode[];
}

export class SelectTreeFlatNode implements SelectTreeNode {
  name!: string;
  level!: number;
  path?: string;
  expandable!: boolean;
  expanded = false;
}

export type TreeListSelectOptions = {
  headerText: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  emptyTemplate: TemplateRef<any> | null;
  hintText: string;
};

@Directive({
  selector: '[ktbTreeListSelect]',
})
export class KtbTreeListSelectDirective extends KtbOverlay implements OnInit, OnDestroy {
  private contentRef: ComponentRef<KtbTreeListSelectComponent> | undefined;

  @Input() data: SelectTreeNode[] = [];
  @Input() options: TreeListSelectOptions = { headerText: '', emptyTemplate: null, hintText: '' };
  @Output() selected: EventEmitter<string> = new EventEmitter<string>();

  constructor(protected elementRef: ElementRef, protected overlayService: OverlayService) {
    super(elementRef, overlayService, '400px', '200px');
  }

  ngOnInit(): void {
    this.onInit();
  }

  ngOnDestroy(): void {
    this.onDestroy();
  }

  @HostListener('click')
  show(): void {
    const treeListSelectPortal: ComponentPortal<KtbTreeListSelectComponent> = new ComponentPortal(
      // eslint-disable-next-line @typescript-eslint/no-use-before-define
      KtbTreeListSelectComponent
    );
    // Disable origin to prevent 'Host has already a portal attached' error
    this.elementRef.nativeElement.disabled = true;

    this.contentRef = this.overlayRef?.attach(treeListSelectPortal);
    if (this.contentRef) {
      this.contentRef.instance.data = this.data;
      this.contentRef.instance.options = this.options;
      this.contentRef.instance.closeDialog.subscribe(() => {
        this.close();
      });

      this.contentRef.instance.selected.subscribe((selected) => {
        this.selected.emit(selected);
        this.close();
      });
    }
  }
}

@Component({
  selector: 'ktb-tree-list-select',
  templateUrl: './ktb-tree-list-select.component.html',
  styleUrls: ['./ktb-tree-list-select.component.scss'],
})
export class KtbTreeListSelectComponent {
  private treeFlattener: DtTreeFlattener<SelectTreeNode, SelectTreeFlatNode> = new DtTreeFlattener(
    this.treeTransformer,
    this.getNodeLevel,
    this.isNodeExpandable,
    this.getNodeChildren
  );
  public treeControl: FlatTreeControl<SelectTreeFlatNode> = new DtTreeControl<SelectTreeFlatNode>(
    this.getNodeLevel,
    this.isNodeExpandable
  );
  public dataSource: DtTreeDataSource<SelectTreeNode, SelectTreeFlatNode> = new DtTreeDataSource(
    this.treeControl,
    this.treeFlattener
  );

  @Input()
  set data(data: SelectTreeNode[]) {
    this.dataSource.data = data;
  }

  @Input() options: TreeListSelectOptions = { headerText: '', emptyTemplate: null, hintText: '' };

  @Output() closeDialog: EventEmitter<void> = new EventEmitter<void>();
  @Output() selected: EventEmitter<string> = new EventEmitter<string>();

  private getNodeLevel(node: SelectTreeFlatNode): number {
    return node.level;
  }

  private isNodeExpandable(node: SelectTreeFlatNode): boolean {
    return node.expandable;
  }

  private getNodeChildren(node: SelectTreeNode): SelectTreeNode[] | undefined {
    return node.keys;
  }

  private treeTransformer(node: SelectTreeNode, level: number): SelectTreeFlatNode {
    const flatNode = new SelectTreeFlatNode();
    flatNode.name = node.name;
    flatNode.level = level;
    flatNode.expandable = !!node.keys;
    flatNode.path = node.path || undefined;
    return flatNode;
  }

  public handleClick(row: SelectTreeFlatNode): void {
    if (row.expandable) {
      if (row.expanded) {
        this.treeControl.collapse(row);
      } else {
        this.treeControl.expand(row);
      }
      row.expanded = !row.expanded;
    } else {
      this.selected.emit(row.path);
    }
  }
}
