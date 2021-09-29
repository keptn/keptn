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

@Directive({
  selector: '[ktbTreeListSelect]',
})
export class KtbTreeListSelectDirective implements OnInit {
  private overlayRef?: OverlayRef;
  private contentRef: ComponentRef<KtbTreeListSelectComponent> | undefined;

  @Input() data: SelectTreeNode[] = [];
  @Output() secret: EventEmitter<string> = new EventEmitter<string>();

  @HostListener('click')
  show(): void {
    const tooltipPortal: ComponentPortal<KtbTreeListSelectComponent> = new ComponentPortal(KtbTreeListSelectComponent);
    // Disable origin to prevent 'Host has already a portal attached' error
    this.elementRef.nativeElement.disabled = true;

    this.contentRef = this.overlayRef?.attach(tooltipPortal);
    if (this.contentRef) {
      this.contentRef.instance.data = this.data;
      this.contentRef.instance.closeDialog.subscribe(() => {
        this.close();
      });

      this.contentRef.instance.selectedSecret.subscribe(secret => {
        this.secret.emit(secret);
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
  private secretTreeFlattener: DtTreeFlattener<SelectTreeNode, SelectTreeFlatNode> = new DtTreeFlattener(this.secretTreeTransformer, this.getSecretLevel, this.isSecretExpandable, this.getSecretChildren);
  public secretTreeControl: FlatTreeControl<SelectTreeFlatNode> = new DtTreeControl<SelectTreeFlatNode>(this.getSecretLevel, this.isSecretExpandable);
  public secretDataSource: DtTreeDataSource<SelectTreeNode, SelectTreeFlatNode> = new DtTreeDataSource(this.secretTreeControl, this.secretTreeFlattener);

  @Input()
  set data(data: SelectTreeNode[]) {
    this.secretDataSource.data = data;
  }

  @Output() closeDialog: EventEmitter<void> = new EventEmitter<void>();
  @Output() selectedSecret: EventEmitter<string> = new EventEmitter<string>();

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
    const variable = `{{.${path}}}`;
    this.selectedSecret.emit(variable);
  }
}
