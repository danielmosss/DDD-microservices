import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { KunstwerkViewComponent } from './kunstwerk-view';
import { DataService } from '../../../data.service';

describe('Kusntwerkview', () => {
  let component: KunstwerkViewComponent;
  let fixture: ComponentFixture<KunstwerkViewComponent>;
  const dataServiceStub = {
    getKunstwerkTree: () => of({
      kunstwerkdetail: {
        kunstwerk: {
          id: 1,
          beheerIdentifier: 'KW-1',
          naam: 'Kunstwerk 1',
          geolocation: null,
          kunstwerkTypeId: null,
          beschrijving: null,
          deleted: false,
          lastsenddhupdate: null,
        },
        KunstwerkType: { id: 1, beschrijving: 'Brug', naam: 'Brug' },
      },
      onderdelen: [],
      losseSensoren: [],
      sensoren_details: [],
    }),
    getKunstwerkenSensorenBulk: () => of([]),
    createOnderdeel: () => of({ id: 2, naam: 'Nieuw', kunstwerkId: 1, parentOnderdeelId: null, onderdelen: [], sensoren: [], sensoren_details: [] }),
    createSensorForOnderdeel: () => of({ id: 10 }),
    getSensorTypes: () => of([]),
    getSensorConfiguratieBronnen: () => of([]),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KunstwerkViewComponent],
      providers: [
        provideRouter([]),
        { provide: DataService, useValue: dataServiceStub },
        { provide: ActivatedRoute, useValue: { snapshot: { paramMap: { get: () => '1' } } } },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KunstwerkViewComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
