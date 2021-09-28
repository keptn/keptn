import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSelectableTileComponent } from './ktb-selectable-tile.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbSelectableTileComponent', () => {
  let component: KtbSelectableTileComponent;
  let fixture: ComponentFixture<KtbSelectableTileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSelectableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should add and remove the selected state', () => {
    const nativeElement = fixture.nativeElement;

    component.selected = false;
    fixture.detectChanges();
    expect(component.selected).toEqual(false);
    expect(nativeElement.getAttribute('class')).not.toContain('ktb-tile-selected');

    component.selected = true;
    fixture.detectChanges();
    expect(component.selected).toEqual(true);
    expect(nativeElement.getAttribute('class')).toContain('ktb-tile-selected');

    component.selected = false;
    fixture.detectChanges();
    expect(component.selected).toEqual(false);
    expect(nativeElement.getAttribute('class')).not.toContain('ktb-tile-selected');
  });
});
