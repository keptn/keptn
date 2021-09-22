import { Component, Directive, ElementRef, EventEmitter, HostListener, Input, OnInit } from '@angular/core';
import { DtTreeControl, DtTreeDataSource, DtTreeFlattener } from '@dynatrace/barista-components/core';
import { FlatTreeControl } from '@angular/cdk/tree';
import { Overlay, OverlayPositionBuilder, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';

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

@Directive({
  selector: '[ktbTreeListSelect]',
})
export class KtbTreeListSelectDirective implements OnInit {
  private overlayRef?: OverlayRef;

  @Input() data: SelectTreeNode[] = [];

  @HostListener('click')
  show(): void {
    const tooltipPortal = new ComponentPortal(KtbTreeListSelectComponent);

    // @ts-ignore
    const contentRef = this.overlayRef.attach(tooltipPortal);
    contentRef.instance.data = this.data;
    contentRef.instance.closeDialog.subscribe(() => {
      this.close();
    });
  }

  constructor(private overlay: Overlay, private overlayPositionBuilder: OverlayPositionBuilder, private elementRef: ElementRef) {
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
    this.overlayRef?.detach();
  }
}


@Component({
  selector: 'ktb-tree-list-select',
  templateUrl: './ktb-tree-list-select.component.html',
  styleUrls: ['./ktb-tree-list-select.component.scss'],
})
export class KtbTreeListSelectComponent implements OnInit {
  private secretTreeFlattener: DtTreeFlattener<SelectTreeNode, SelectTreeFlatNode> = new DtTreeFlattener(this.secretTreeTransformer, this.getSecretLevel, this.isSecretExpandable, this.getSecretChildren);
  public secretTreeControl: FlatTreeControl<SelectTreeFlatNode> = new DtTreeControl<SelectTreeFlatNode>(this.getSecretLevel, this.isSecretExpandable);
  public secretDataSource: DtTreeDataSource<SelectTreeNode, SelectTreeFlatNode> = new DtTreeDataSource(this.secretTreeControl, this.secretTreeFlattener);
  public closeDialog: EventEmitter<void> = new EventEmitter<void>();

  @Input()
  set data(data: SelectTreeNode[]) {
    this.secretDataSource.data = data;
  }

  constructor() {
  }

  ngOnInit(): void {
  }

  private getSecretLevel(node: SelectTreeFlatNode): number {
    return node.level;
  }

  private isSecretExpandable(node: SelectTreeFlatNode): boolean {
    return node.expandable;
  }

  private getSecretChildren(node: SelectTreeNode): SelectTreeNode[] | undefined {
    return node.keys;
  }

  private secretTreeTransformer(node: SelectTreeNode, level: number): SelectTreeFlatNode {
    const flatNode = new SelectTreeFlatNode();
    flatNode.name = node.name;
    flatNode.level = level;
    flatNode.expandable = !!node.keys;
    flatNode.path = node.path || undefined;
    return flatNode;
  }

  public selectSecret(path: string): void {
    // TODO determine last focused field and paste
  }
}
