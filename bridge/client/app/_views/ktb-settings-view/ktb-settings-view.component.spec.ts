import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { KtbSettingsViewModule } from './ktb-settings-view.module';
import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbSettingsViewComponent', () => {
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSettingsViewModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSettingsViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have 5 entries in the submenu', () => {
    // given
    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu.submenu .dt-menu-item'));

    // then
    expect(menuItems).toBeTruthy();
    expect(menuItems.length).toEqual(5);
    expect(menuItems[0].nativeElement.textContent.trim()).toEqual('Project');
    expect(menuItems[1].nativeElement.textContent.trim()).toEqual('Services');
    expect(menuItems[2].nativeElement.textContent.trim()).toEqual('Integrations');
    expect(menuItems[3].nativeElement.textContent.trim()).toEqual('Secrets');
    expect(menuItems[4].nativeElement.textContent.trim()).toEqual('Common use cases');
  });
});
