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
} from '@angular/core';
import { DtTreeControl, DtTreeDataSource, DtTreeFlattener } from '@dynatrace/barista-components/core';
import { FlatTreeControl } from '@angular/cdk/tree';
import { OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { Router } from '@angular/router';
import { OverlayService } from '../../_directives/overlay-service/overlay.service';
import { Subject } from 'rxjs';

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
  emptyText: string;
  hintText: string;
};

@Directive({
  selector: '[ktbTreeListSelect]',
})
export class KtbTreeListSelectDirective implements OnInit, OnDestroy {
  private overlayRef?: OverlayRef;
  private contentRef: ComponentRef<KtbTreeListSelectComponent> | undefined;
  private unsubscribe$: Subject<void> = new Subject();

  @Input() data: SelectTreeNode[] = [];
  @Input() options: TreeListSelectOptions = { headerText: '', emptyText: '', hintText: '' };
  @Output() selected: EventEmitter<string> = new EventEmitter<string>();

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

  constructor(private elementRef: ElementRef, private router: Router, private overlayService: OverlayService) {
    overlayService.registerNavigationEvent(this.unsubscribe$, this.close.bind(this));
  }

  public ngOnInit(): void {
    this.overlayRef = this.overlayService.initOverlay('400px', '200px', true, this.elementRef, this.close.bind(this));
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public close(): void {
    this.overlayService.closeOverlay(this.overlayRef, this.elementRef);
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

  @Input() options: TreeListSelectOptions = { headerText: '', emptyText: '', hintText: '' };

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
