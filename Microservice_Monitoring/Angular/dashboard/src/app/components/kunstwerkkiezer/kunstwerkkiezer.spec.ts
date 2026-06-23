import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { Kunstwerkkiezer } from './kunstwerkkiezer';
import { DataService } from '../../../data.service';

describe('Kunstwerkkiezer', () => {
  let component: Kunstwerkkiezer;
  let fixture: ComponentFixture<Kunstwerkkiezer>;
  const dataServiceStub = {
    getKunstwerken: () => of([]),
    getKunstwerkDHU: () => of(null),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Kunstwerkkiezer],
      providers: [provideRouter([]), { provide: DataService, useValue: dataServiceStub }],
    }).compileComponents();

    fixture = TestBed.createComponent(Kunstwerkkiezer);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
