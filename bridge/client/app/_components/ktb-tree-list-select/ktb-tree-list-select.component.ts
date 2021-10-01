import { Component, ComponentRef, Directive, ElementRef, EventEmitter, HostListener, Input, OnInit, Output } from '@angular/core';
import { DtTreeControl, DtTreeDataSource, DtTreeFlattener } from '@dynatrace/barista-components/core';
import { FlatTreeControl } from '@angular/cdk/tree';
import { Overlay, OverlayPositionBuilder, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { NavigationStart, Router } from '@angular/router';
import { filter } from 'rxjs/operators';

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
}

export type TreeListSelectOptions = {
  headerText: string;
  emptyText: string;
};

@Directive({
  selector: '[ktbTreeListSelect]',
})
export class KtbTreeListSelectDirective implements OnInit {
  private overlayRef?: OverlayRef;
  private contentRef: ComponentRef<KtbTreeListSelectComponent> | undefined;

  @Input() data: SelectTreeNode[] = [];
  @Input() options: TreeListSelectOptions = {headerText: '', emptyText: ''};
  @Output() selected: EventEmitter<string> = new EventEmitter<string>();

  @HostListener('click')
  show(): void {
    const tooltipPortal: ComponentPortal<KtbTreeListSelectComponent> = new ComponentPortal(KtbTreeListSelectComponent);
    // Disable origin to prevent 'Host has already a portal attached' error
    this.elementRef.nativeElement.disabled = true;

    this.contentRef = this.overlayRef?.attach(tooltipPortal);
    if (this.contentRef) {
      this.contentRef.instance.data = this.data;
      this.contentRef.instance.options = this.options;
      this.contentRef.instance.closeDialog.subscribe(() => {
        this.close();
      });

      this.contentRef.instance.selected.subscribe(selected => {
        this.selected.emit(selected);
        this.close();
      });
    }
  }

  constructor(private overlay: Overlay, private overlayPositionBuilder: OverlayPositionBuilder, private elementRef: ElementRef, private router: Router) {
    // Close when navigation happens - to keep the overlay on the UI
    this.router.events.pipe(filter(event => event instanceof NavigationStart)).subscribe(() => {
      this.close();
    });
  }

  public ngOnInit(): void {
    const positionStrategy = this.overlayPositionBuilder
      .flexibleConnectedTo(this.elementRef)
      .withPositions([{
        originX: 'start',
        originY: 'bottom',
        overlayX: 'start',
        overlayY: 'top',
        offsetY: 10,
        offsetX: -20,
      }]);


    this.overlayRef = this.overlay.create({positionStrategy, width: '400px', height: '200px'});
  }

  public close(): void {
    this.elementRef.nativeElement.disabled = false;
    this.overlayRef?.detach();
  }
}


@Component({
  selector: 'ktb-tree-list-select',
  templateUrl: './ktb-tree-list-select.component.html',
  styleUrls: ['./ktb-tree-list-select.component.scss'],
})
export class KtbTreeListSelectComponent {
  private treeFlattener: DtTreeFlattener<SelectTreeNode, SelectTreeFlatNode> = new DtTreeFlattener(this.treeTransformer, this.getNodeLevel, this.isNodeExpandable, this.getNodeChildren);
  public treeControl: FlatTreeControl<SelectTreeFlatNode> = new DtTreeControl<SelectTreeFlatNode>(this.getNodeLevel, this.isNodeExpandable);
  public dataSource: DtTreeDataSource<SelectTreeNode, SelectTreeFlatNode> = new DtTreeDataSource(this.treeControl, this.treeFlattener);

  @Input()
  set data(data: SelectTreeNode[]) {
    this.dataSource.data = data;
  }

  @Input() options: TreeListSelectOptions = {headerText: '', emptyText: ''};

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

  public selectValue(path: string): void {
    this.selected.emit(path);
  }
}
