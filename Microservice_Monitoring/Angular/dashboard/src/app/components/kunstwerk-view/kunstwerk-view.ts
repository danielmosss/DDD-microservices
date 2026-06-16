import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KunstwerkTreeView, Onderdeel } from '../../models/types';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-kunstwerk-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './kunstwerk-view.html',
  styleUrls: ['./kunstwerk-view.scss'],
})
export class KunstwerkViewComponent implements OnInit{
  public kunstwerkId: number = 0;
  private activatedRoute = inject(ActivatedRoute)

  ngOnInit(): void {
    console.log(this.activatedRoute)
    this.kunstwerkId = parseInt(this.activatedRoute.snapshot.paramMap.get("kunstwerkid") ?? '0')
    console.log(this.kunstwerkId)
  }
  // Dit is je Mock Data
  treeView: KunstwerkTreeView = {
    kunstwerk: {
      id: 1,
      beheerIdentifier: 'BRUG-001',
      naam: 'Beweegbare Brug Complex',
      geolocation: '52.3702, 4.8952',
      kunstwerkTypeId: 1,
      beschrijving: 'Grote basculebrug met meerdere complexe aandrijfsystemen',
      deleted: false,
      lastsenddhupdate: new Date().toISOString(),
    },
    // Sensoren die op het hoogste niveau zitten (niet gekoppeld aan een specifiek onderdeel)
    losseSensoren: [
      {
        id: 10,
        kunstwerkId: 1,
        onderdeelId: null,
        geolocation: 'Dak bedieningshuis',
        sensorTypeId: 1,
        lastAnalyzedMetingId: null,
        laatsteMeting: 18.5, // Bijv. Buitentemperatuur
        sensorConfiguratie: {
          id: 1,
          sensorId: 10,
          minValue: -20,
          maxValue: 40,
          margePercentage: null,
        },
      },
    ],
    onderdelen: [
      // --- HOOFDONDERDEEL 1: Statische structuur ---
      {
        id: 101,
        naam: 'Pijler Noord',
        kunstwerkId: 1,
        parentOnderdeelId: null,
        sensoren: [
          {
            id: 501,
            kunstwerkId: 1,
            onderdeelId: 101,
            geolocation: 'Waterlijn',
            sensorTypeId: 4,
            lastAnalyzedMetingId: null,
            laatsteMeting: -0.02, // Bijv. Verzakking in mm
            sensorConfiguratie: {
              id: 11,
              sensorId: 501,
              minValue: -0.5,
              maxValue: 0.5,
              margePercentage: 5,
            },
          },
        ],
        onderdelen: [
          {
            id: 201,
            naam: 'Fundering Beton',
            kunstwerkId: 1,
            parentOnderdeelId: 101,
            sensoren: [
              {
                id: 502,
                kunstwerkId: 1,
                onderdeelId: 201,
                geolocation: null,
                sensorTypeId: 5,
                lastAnalyzedMetingId: null,
                laatsteMeting: 450, // Bijv. Druk in kPa
                sensorConfiguratie: {
                  id: 12,
                  sensorId: 502,
                  minValue: 0,
                  maxValue: 1000,
                  margePercentage: 10,
                },
              },
            ],
            onderdelen: [],
          },
        ],
      },
      // --- HOOFDONDERDEEL 2: Het complexe, diep geneste aandrijfsysteem ---
      {
        id: 102,
        naam: 'Aandrijfsysteem Bascule',
        kunstwerkId: 1,
        parentOnderdeelId: null,
        sensoren: [
          {
            id: 503,
            kunstwerkId: 1,
            onderdeelId: 102,
            geolocation: 'Machinekamer',
            sensorTypeId: 6,
            lastAnalyzedMetingId: null,
            laatsteMeting: 42.1, // Bijv. Omgevingstemperatuur
            sensorConfiguratie: {
              id: 13,
              sensorId: 503,
              minValue: 0,
              maxValue: 60,
              margePercentage: 5,
            },
          },
        ],
        onderdelen: [
          {
            id: 202,
            naam: 'Hoofdmotor',
            kunstwerkId: 1,
            parentOnderdeelId: 102,
            sensoren: [
              {
                id: 504,
                kunstwerkId: 1,
                onderdeelId: 202,
                geolocation: 'Motorhuis',
                sensorTypeId: 7,
                lastAnalyzedMetingId: null,
                laatsteMeting: 1450, // RPM
                sensorConfiguratie: {
                  id: 14,
                  sensorId: 504,
                  minValue: 0,
                  maxValue: 1500,
                  margePercentage: 2,
                },
              },
              {
                id: 505,
                kunstwerkId: 1,
                onderdeelId: 202,
                geolocation: 'Koeling',
                sensorTypeId: 6,
                lastAnalyzedMetingId: null,
                laatsteMeting: 82.5, // Temperatuur
                sensorConfiguratie: {
                  id: 15,
                  sensorId: 505,
                  minValue: 0,
                  maxValue: 95,
                  margePercentage: 10,
                },
              },
            ],
            // Sub-sub-onderdelen van de motor
            onderdelen: [
              {
                id: 301,
                naam: 'Aandrijfas',
                kunstwerkId: 1,
                parentOnderdeelId: 202,
                sensoren: [
                  {
                    id: 506,
                    kunstwerkId: 1,
                    onderdeelId: 301,
                    geolocation: null,
                    sensorTypeId: 8,
                    lastAnalyzedMetingId: null,
                    laatsteMeting: 0.04, // Trilling/Vibratie
                    sensorConfiguratie: {
                      id: 16,
                      sensorId: 506,
                      minValue: 0,
                      maxValue: 0.1,
                      margePercentage: 5,
                    },
                  },
                ],
                onderdelen: [],
              },
              {
                id: 302,
                naam: 'Remmenset',
                kunstwerkId: 1,
                parentOnderdeelId: 202,
                sensoren: [],
                // Sub-sub-sub-onderdelen (4 niveaus diep!)
                onderdelen: [
                  {
                    id: 401,
                    naam: 'Remblok Links',
                    kunstwerkId: 1,
                    parentOnderdeelId: 302,
                    sensoren: [
                      {
                        id: 507,
                        kunstwerkId: 1,
                        onderdeelId: 401,
                        geolocation: null,
                        sensorTypeId: 9,
                        lastAnalyzedMetingId: null,
                        laatsteMeting: 14.2, // Slijtage indicator
                        sensorConfiguratie: {
                          id: 17,
                          sensorId: 507,
                          minValue: 5,
                          maxValue: 30,
                          margePercentage: null,
                        },
                      },
                    ],
                    onderdelen: [],
                  },
                  {
                    id: 402,
                    naam: 'Remblok Rechts',
                    kunstwerkId: 1,
                    parentOnderdeelId: 302,
                    sensoren: [
                      {
                        id: 508,
                        kunstwerkId: 1,
                        onderdeelId: 402,
                        geolocation: null,
                        sensorTypeId: 9,
                        lastAnalyzedMetingId: null,
                        laatsteMeting: 13.8, // Slijtage indicator
                        sensorConfiguratie: {
                          id: 18,
                          sensorId: 508,
                          minValue: 5,
                          maxValue: 30,
                          margePercentage: null,
                        },
                      },
                    ],
                    onderdelen: [],
                  },
                ],
              },
            ],
          },
        ],
      },
    ],
  };

  selectedOnderdeel: Onderdeel | null = null;

  selectOnderdeel(onderdeel: Onderdeel) {
    this.selectedOnderdeel = onderdeel;
    console.log('Nu geselecteerd:', onderdeel.naam);
  }

  getSensorStatus(val: number | undefined): string {
    if (!val) return 'Geen data';
    return val > 100 ? 'Warning' : 'OK';
  }
}
