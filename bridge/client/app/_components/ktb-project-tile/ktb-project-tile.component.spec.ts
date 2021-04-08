import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbProjectTileComponent } from './ktb-project-tile.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbProjectListComponent} from "../ktb-project-list/ktb-project-list.component";

describe('KtbProjectTileComponent', () => {
  let component: KtbProjectTileComponent;
  let fixture: ComponentFixture<KtbProjectTileComponent>;

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
        fixture = TestBed.createComponent(KtbProjectTileComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
