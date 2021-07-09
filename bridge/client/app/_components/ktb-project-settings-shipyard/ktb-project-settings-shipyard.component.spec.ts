import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsShipyardComponent } from './ktb-project-settings-shipyard.component';

describe('KtbProjectSettingsEditProjectComponent', () => {
  let component: KtbProjectSettingsShipyardComponent;
  let fixture: ComponentFixture<KtbProjectSettingsShipyardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProjectSettingsShipyardComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectSettingsShipyardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
