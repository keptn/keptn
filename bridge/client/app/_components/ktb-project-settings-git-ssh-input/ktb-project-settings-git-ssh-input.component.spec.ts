import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitSshInputComponent } from './ktb-project-settings-git-ssh-input.component';

describe('KtbProjectSettingsGitSshInputComponent', () => {
  let component: KtbProjectSettingsGitSshInputComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitSshInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProjectSettingsGitSshInputComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectSettingsGitSshInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
