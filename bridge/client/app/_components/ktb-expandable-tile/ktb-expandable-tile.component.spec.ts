import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {KtbExpandableTileComponent, KtbExpandableTileHeader} from './ktb-expandable-tile.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbExpandableTileComponent', () => {
  let component: KtbExpandableTileComponent;
  let fixture: ComponentFixture<KtbExpandableTileComponent>;

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
    fixture = TestBed.createComponent(KtbExpandableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
