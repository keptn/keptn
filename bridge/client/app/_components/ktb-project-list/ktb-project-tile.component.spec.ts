import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectListModule } from './ktb-project-list.module';
import { KtbProjectTileComponent } from './ktb-project-tile.component';

describe('KtbProjectTileComponent', () => {
  let component: KtbProjectTileComponent;
  let fixture: ComponentFixture<KtbProjectTileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectListModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
