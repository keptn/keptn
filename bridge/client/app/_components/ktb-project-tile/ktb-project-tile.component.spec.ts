import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectTileComponent } from './ktb-project-tile.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

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
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
