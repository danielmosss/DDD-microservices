import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';

import { DataService } from './data.service';

describe('DataService', () => {
  let service: DataService;
  let httpTesting: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting(), provideRouter([])],
    });
    service = TestBed.inject(DataService);
    httpTesting = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpTesting.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should create an onderdeel', () => {
    service.createOnderdeel(12, { naam: 'Fundering', parentOnderdeelId: 5 }).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/kunstwerken/12/onderdelen');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual({ naam: 'Fundering', parentOnderdeelId: 5 });
    request.flush({ id: 30, naam: 'Fundering', kunstwerkId: 12, parentOnderdeelId: 5, onderdelen: [], sensoren: [], sensoren_details: [] });
  });

  it('should create a sensor for an onderdeel with configuration', () => {
    service.createSensorForOnderdeel(12, 30, {
      sensorTypeId: 3,
      geolocation: null,
      configuratieBronSensorId: 44,
      configuratie: { minValue: 1, maxValue: 10, margePercentage: 5 },
    }).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/kunstwerken/12/onderdelen/30/sensoren');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual({
      sensorTypeId: 3,
      geolocation: null,
      configuratieBronSensorId: 44,
      configuratie: { minValue: 1, maxValue: 10, margePercentage: 5 },
    });
    request.flush({ id: 44 });
  });

  it('should delete an onderdeel with data preservation choice', () => {
    service.deleteOnderdeel(12, 30, { preserveSensorData: true }).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/kunstwerken/12/onderdelen/30');
    expect(request.request.method).toBe('DELETE');
    expect(request.request.body).toEqual({ preserveSensorData: true });
    request.flush({ deleted: true, preserveSensorData: true });
  });

  it('should delete a sensor with data preservation choice', () => {
    service.deleteSensor(44, { preserveSensorData: false }).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/sensoren/44');
    expect(request.request.method).toBe('DELETE');
    expect(request.request.body).toEqual({ preserveSensorData: false });
    request.flush({ deleted: true, preserveSensorData: false });
  });

  it('should get sensor types', () => {
    service.getSensorTypes().subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/sensortypes');
    expect(request.request.method).toBe('GET');
    request.flush([]);
  });

  it('should get sensor configuration copy sources', () => {
    service.getSensorConfiguratieBronnen(12, 3).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/kunstwerken/12/sensor-configuratie-bronnen?sensorTypeId=3');
    expect(request.request.method).toBe('GET');
    request.flush([]);
  });

  it('should update sensor configuration', () => {
    service.updateSensorConfiguratie(44, { minValue: 1, maxValue: 10, margePercentage: 5 }).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/api/v1/sensoren/44/configuratie');
    expect(request.request.method).toBe('PUT');
    expect(request.request.body).toEqual({ minValue: 1, maxValue: 10, margePercentage: 5 });
    request.flush({ id: 1, sensorId: 44, minValue: 1, maxValue: 10, margePercentage: 5 });
  });
});
