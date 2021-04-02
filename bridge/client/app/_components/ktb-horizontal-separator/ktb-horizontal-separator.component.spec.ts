import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {KtbHorizontalSeparatorComponent, KtbHorizontalSeparatorTitle} from './ktb-horizontal-separator.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {KtbExpandableTileComponent} from "../ktb-expandable-tile/ktb-expandable-tile.component";

describe('KtbHorizontalSeparatorComponent', () => {
  let component: KtbHorizontalSeparatorComponent;
  let fixture: ComponentFixture<KtbHorizontalSeparatorComponent>;

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
        fixture = TestBed.createComponent(KtbHorizontalSeparatorComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
