import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitHttpsComponent } from './ktb-project-settings-git-https.component';

describe('KtbProjectSettingsGitHttpsComponent', () => {
  let component: KtbProjectSettingsGitHttpsComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitHttpsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProjectSettingsGitHttpsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectSettingsGitHttpsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
