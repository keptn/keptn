import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitComponent } from './ktb-project-settings-git.component';

describe('KtbProjectSettingsGitComponent', () => {
  let component: KtbProjectSettingsGitComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProjectSettingsGitComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectSettingsGitComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
