import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbExpandableTileComponent } from './ktb-expandable-tile.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbExpandableTileModule } from './ktb-expandable-tile.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('KtbExpandableTileComponent', () => {
  let component: KtbExpandableTileComponent;
  let fixture: ComponentFixture<KtbExpandableTileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbExpandableTileModule, BrowserAnimationsModule, HttpClientTestingModule],
    }).compileComponents();

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
    expect(showMoreButton.getAttribute('class')).not.toContain('dt-show-more-disabled');
    expect(expandablePanel.getAttribute('class')).toContain('dt-expandable-panel-opened');
    expect(showMoreButton.disabled).toBe(false);
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
    expect(showMoreButton.getAttribute('class')).toContain('dt-show-more-disabled');
    expect(expandablePanel.getAttribute('class')).not.toContain('dt-expandable-panel-opened');
    expect(showMoreButton.disabled).toBe(true);
  });
});
