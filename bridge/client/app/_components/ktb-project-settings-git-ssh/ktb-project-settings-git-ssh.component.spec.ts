import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitSshComponent } from './ktb-project-settings-git-ssh.component';

describe('KtbProjectSettingsGitSshComponent', () => {
  let component: KtbProjectSettingsGitSshComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitSshComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProjectSettingsGitSshComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectSettingsGitSshComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
