import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSshKeyInputComponent } from './ktb-ssh-key-input.component';

describe('KtbSshKeyInputComponent', () => {
  let component: KtbSshKeyInputComponent;
  let fixture: ComponentFixture<KtbSshKeyInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbSshKeyInputComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSshKeyInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
