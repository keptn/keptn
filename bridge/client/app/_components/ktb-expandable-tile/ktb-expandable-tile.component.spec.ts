import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';

import {KtbExpandableTileComponent} from './ktb-expandable-tile.component';
import {AppModule} from '../../app.module';
import { KtbDangerZoneComponent } from '../ktb-danger-zone/ktb-danger-zone.component';
import { DeleteType } from '../../_interfaces/delete';

describe('KtbExpandableTileComponent', () => {
  let component: KtbExpandableTileComponent;
  let fixture: ComponentFixture<KtbExpandableTileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbDangerZoneComponent],
      imports: [AppModule],
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbExpandableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should be expandable by default', () => {
    // given
    fixture.detectChanges();
    const showMoreButton = fixture.nativeElement.querySelector('.dt-show-more');
    const expandablePanel = fixture.nativeElement.querySelector('.dt-expandable-panel');

    // when
    showMoreButton.click();
    fixture.detectChanges();

    // then
    // expect(component).not.toHaveClass('ktb-tile-disabled');
    expect(showMoreButton).not.toHaveClass('dt-show-more-disabled');
    expect(showMoreButton.disabled).toBeFalse();
    expect(expandablePanel).toHaveClass('dt-expandable-panel-opened');
  });

  it('should not be expandable if disabled', () => {
    // given
    component.disabled = true;
    fixture.detectChanges();
    const showMoreButton = fixture.nativeElement.querySelector('.dt-show-more');
    const expandablePanel = fixture.nativeElement.querySelector('.dt-expandable-panel');

    // when
    showMoreButton.click();
    fixture.detectChanges();

    // then
    // expect(component).toHaveClass('ktb-tile-disabled');
    expect(showMoreButton).toHaveClass('dt-show-more-disabled');
    expect(showMoreButton.disabled).toBeTrue();
    expect(expandablePanel).not.toHaveClass('dt-expandable-panel-opened');
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
