import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {KtbMarkdownComponent} from './ktb-markdown.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbKeptnServicesListComponent} from "../ktb-keptn-services-list/ktb-keptn-services-list.component";

describe('KtbExpandableTileComponent', () => {
  let component: KtbMarkdownComponent;
  let fixture: ComponentFixture<KtbMarkdownComponent>;

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
        fixture = TestBed.createComponent(KtbMarkdownComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
