import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import { AppModule } from '../../app.module';
import { By } from '@angular/platform-browser';

describe('KtbSettingsViewComponent', () => {
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSettingsViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have 1 entry in the submenu', () => {
    // given
    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu.submenu .dt-menu-item'));

    // then
    expect(menuItems).toBeTruthy();
    expect(menuItems.length).toEqual(2);
    expect(menuItems[0].nativeElement.textContent.trim()).toEqual('Project');
    expect(menuItems[1].nativeElement.textContent.trim()).toEqual('Services');
  });
});
