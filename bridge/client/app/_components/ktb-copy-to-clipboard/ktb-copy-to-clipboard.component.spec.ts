import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbCopyToClipboardComponent } from './ktb-copy-to-clipboard.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbCopyToClipboardModule } from './ktb-copy-to-clipboard.module';

describe('KtbExpandableTileComponent', () => {
  let component: KtbCopyToClipboardComponent;
  let fixture: ComponentFixture<KtbCopyToClipboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbCopyToClipboardModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbCopyToClipboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
