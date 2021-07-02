import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsEditProjectComponent } from './ktb-project-settings-edit-project.component';

describe('KtbProjectSettingsEditProjectComponent', () => {
  let component: KtbProjectSettingsEditProjectComponent;
  let fixture: ComponentFixture<KtbProjectSettingsEditProjectComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProjectSettingsEditProjectComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectSettingsEditProjectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
