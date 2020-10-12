import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {KtbMarkdownComponent} from './ktb-markdown.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbExpandableTileComponent', () => {
  let component: KtbMarkdownComponent;
  let fixture: ComponentFixture<KtbMarkdownComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbMarkdownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
