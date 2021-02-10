import {ComponentFixture, TestBed, fakeAsync} from '@angular/core/testing';
import {By} from "@angular/platform-browser";
import {Component} from "@angular/core";

import {KtbSelectableTileComponent} from './ktb-selectable-tile.component';

describe('KtbSelectableTileComponent', () => {
  let component: SimpleKtbSelectableTileComponent;
  let fixture: ComponentFixture<SimpleKtbSelectableTileComponent>;

  beforeEach(fakeAsync(() => {
    TestBed.configureTestingModule({
      declarations: [
        SimpleKtbSelectableTileComponent,
        KtbSelectableTileComponent,
      ],
      imports: [],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SimpleKtbSelectableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should add and remove the selected state', () => {
    let selectableTileDebugElement = fixture.debugElement.query(By.directive(KtbSelectableTileComponent));
    let selectableTileInstance = selectableTileDebugElement.componentInstance;
    let selectableTileNativeElement = selectableTileDebugElement.nativeElement;
    let testComponentInstance = fixture.debugElement.componentInstance;

    expect(selectableTileInstance.selected).toBe(false);

    testComponentInstance.isSelected = true;
    fixture.detectChanges();

    expect(selectableTileInstance.selected).toBe(true);
    expect(selectableTileNativeElement.classList).toContain('ktb-tile-selected');

    testComponentInstance.isSelected = false;
    fixture.detectChanges();

    expect(selectableTileInstance.selected).toBe(false);
    expect(selectableTileNativeElement.classList).not.toContain('ktb-tile-selected');
  });
});

/** Simple component for testing the KtbSelectableTileComponent */
@Component({
  template: `
    <div>
      <ktb-selectable-tile
        [error]="isError"
        [success]="isSuccess"
        [selected]="isSelected"
        (click)="onTileClicked($event)">
      </ktb-selectable-tile>
    </div>
  `,
})
class SimpleKtbSelectableTileComponent {
  isError = false;
  isSuccess = false;
  isSelected = false;

  onTileClicked: (event?: Event) => void = () => { this.isSelected = !this.isSelected; };
}
