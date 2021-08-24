import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServicesListComponent } from './ktb-services-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { Service } from '../../_models/service';
import { ServiceMock } from '../../_models/service.mock';

describe('KtbServicesListComponent', () => {
  let component: KtbServicesListComponent;
  let fixture: ComponentFixture<KtbServicesListComponent>;
  let service: Service;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbServicesListComponent);
    component = fixture.componentInstance;
    service = Service.fromJSON(ServiceMock);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should return an image string (<image>:<version>)', () => {
    // when
    const imageString = component.getImageText(service);

    // then
    expect(imageString).toEqual('mongo:4.2.2');
  });

  it('should return an image string (<image>) when no version is set', () => {
    // given
    service.deployedImage = 'docker.io/mongo';

    // when
    const imageString = component.getImageText(service);

    // then
    expect(imageString).toEqual('mongo:docker.io/mongo');
  });

  it('should return an empty string when Service.deployedImage is not set', () => {
    // given
    service.deployedImage = undefined;

    // when
    const imageString = component.getImageText(service);

    // then
    expect(imageString).toEqual('');
  });

  it('should display the service name and the deployed image', () => {
    // given
    component.services = [service];
    fixture.detectChanges();

    // when
    const elem = fixture.nativeElement.querySelectorAll('.dt-table-column-serviceName p a span')[1];

    // then
    expect(elem).toBeTruthy();
    expect(elem.textContent).toEqual('mongo:4.2.2');
  });

  it('should return a link to the service', () => {
    // given
    service.stage = 'dev';

    // when
    const link = component.getServiceLink(service);

    // then
    expect(link).toEqual(['service', 'carts-db', 'context', 'ff8a3e69-7e5c-48ec-b668-4e96a006a505', 'stage', 'dev']);
  });
});
