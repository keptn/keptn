import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSelectableTileComponent } from './ktb-selectable-tile.component';
import {By} from "@angular/platform-browser";
import {Component} from "@angular/core";
import {AppModule} from "../../app.module";

describe('KtbSelectableTileComponent', () => {
  let component: SimpleKtbSelectableTileComponent;
  let fixture: ComponentFixture<SimpleKtbSelectableTileComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        SimpleKtbSelectableTileComponent
      ],
      imports: [
        AppModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SimpleKtbSelectableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });

  it('should add and remove the selected state', () => {
    let selectableTileDebugElement = fixture.debugElement.query(By.directive(KtbSelectableTileComponent));
    let selectableTileInstance = selectableTileDebugElement.componentInstance;
    let selectableTileNativeElement = selectableTileDebugElement.nativeElement;
    let testComponentInstace = fixture.debugElement.componentInstance;

    expect(selectableTileInstance.selected).toBe(false);

    testComponentInstace.isSelected = true;
    fixture.detectChanges();

    expect(selectableTileInstance.selected).toBe(true);
    expect(selectableTileNativeElement.classList).toContain('ktb-tile-selected');

    testComponentInstace.isSelected = false;
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
