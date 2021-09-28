import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbCopyToClipboardComponent } from './ktb-copy-to-clipboard.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbExpandableTileComponent', () => {
  let component: KtbCopyToClipboardComponent;
  let fixture: ComponentFixture<KtbCopyToClipboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbCopyToClipboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
