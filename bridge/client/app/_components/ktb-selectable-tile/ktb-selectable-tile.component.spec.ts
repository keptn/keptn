import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSelectableTileComponent } from './ktb-selectable-tile.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('KtbSelectableTileComponent', () => {
  let component: KtbSelectableTileComponent;
  let fixture: ComponentFixture<KtbSelectableTileComponent>;

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
    fixture = TestBed.createComponent(KtbSelectableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
