import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {KtbCopyToClipboardComponent} from './ktb-copy-to-clipboard.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbExpandableTileComponent', () => {
  let component: KtbCopyToClipboardComponent;
  let fixture: ComponentFixture<KtbCopyToClipboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents()
    .then(() => {
      fixture = TestBed.createComponent(KtbCopyToClipboardComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();
    });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
