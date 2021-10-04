import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsListComponent } from './ktb-service-settings-list.component';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbServiceSettingsListComponent', () => {
  let component: KtbServiceSettingsListComponent;
  let fixture: ComponentFixture<KtbServiceSettingsListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        { provide: DataService, useClass: DataServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(
              convertToParamMap({
                projectName: 'sockshop',
              })
            ),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbServiceSettingsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should contain services', () => {
    expect(component.dataSource.data.sort()).toEqual(['carts', 'carts-db']);
  });
});
