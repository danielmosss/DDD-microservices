import { ChangeDetectorRef, Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KunstwerkTreeView, Onderdeel, SensorDetailResponse } from '../../models/types';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { catchError, forkJoin, map, Observable, of, take, tap, timer } from 'rxjs';
import { DataService } from '../../../data.service';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { CreateOnderdeelDialogComponent } from './dialogs/create-onderdeel-dialog/create-onderdeel-dialog';
import { CreateSensorDialogComponent } from './dialogs/create-sensor-dialog/create-sensor-dialog';
import { DeleteDataDialogComponent } from './dialogs/delete-data-dialog/delete-data-dialog';

@Component({
  selector: 'app-kunstwerk-view',
  standalone: true,
  imports: [CommonModule, RouterModule, MatButtonModule, MatDialogModule],
  templateUrl: './kunstwerk-view.html',
  styleUrls: ['./kunstwerk-view.scss'],
})
export class KunstwerkViewComponent implements OnInit {
  public kunstwerkId: number = 0;
  private activatedRoute = inject(ActivatedRoute);
  public treeView$: Observable<KunstwerkTreeView> | null = null;
  public refreshDataTimer$: Observable<number> | null = null;
  public currentTree: KunstwerkTreeView | null = null;

  private _cdr = inject(ChangeDetectorRef);
  private _dialog = inject(MatDialog);

  public isSavingOnderdeel = false;
  public isSavingSensor = false;
  public isPreparingSensorDialog = false;
  public isDeleting = false;
  public managementError: string | null = null;
  public managementSuccess: string | null = null;

  constructor(private _dataService: DataService) {}

  ngOnInit(): void {
    this.kunstwerkId = parseInt(this.activatedRoute.snapshot.paramMap.get('kunstwerkid') ?? '0');
    this.loadTree();
  }

  selectedOnderdeel: Onderdeel | null = null;

  private loadTree(selectOnderdeelId?: number): void {
    this.treeView$ = this._dataService.getKunstwerkTree(this.kunstwerkId).pipe(
      tap((treeData) => {
        this.currentTree = treeData;
        const onderdeelToSelect = selectOnderdeelId
          ? this.findOnderdeelById(treeData.onderdelen, selectOnderdeelId)
          : treeData.onderdelen?.[0];

        if (onderdeelToSelect) {
          this.selectOnderdeel(onderdeelToSelect);
        }
      })
    );
  }

  selectOnderdeel(onderdeel: Onderdeel) {
    this.selectedOnderdeel = onderdeel;
    const onderdelenToRefresh = [onderdeel, ...(onderdeel.onderdelen ?? [])];
    const sensorRequests = onderdelenToRefresh.map((subOnderdeel) => {
      subOnderdeel.sensoren_details = [];
      return this.fetchAndAttachSensors(subOnderdeel);
    });

    forkJoin(sensorRequests).subscribe(() => {
      this.refreshFunction();
    });
  }

  public openCreateOnderdeelDialog(): void {
    const dialogRef = this._dialog.open(CreateOnderdeelDialogComponent, {
      width: '520px',
      data: {
        kunstwerkNaam: this.currentTree?.kunstwerkdetail.kunstwerk.naam ?? '',
        selectedOnderdeel: this.selectedOnderdeel,
      },
    });

    this.managementError = null;
    this.managementSuccess = null;

    dialogRef.afterClosed().subscribe((request) => {
      if (!request || this.isSavingOnderdeel) {
        return;
      }

      this.isSavingOnderdeel = true;
      this.managementError = null;
      this.managementSuccess = null;

      this._dataService.createOnderdeel(this.kunstwerkId, request).subscribe({
        next: (onderdeel) => {
          this.managementSuccess = `Onderdeel "${onderdeel.naam}" is aangemaakt.`;
          this.loadTree(onderdeel.id);
        },
        error: () => {
          this.managementError = 'Onderdeel aanmaken is niet gelukt.';
        },
        complete: () => {
          this.isSavingOnderdeel = false;
        },
      });
    });
  }

  public openCreateSensorDialog(): void {
    if (!this.selectedOnderdeel || this.isSavingSensor || this.isPreparingSensorDialog) {
      return;
    }

    const selectedOnderdeel = this.selectedOnderdeel;
    this.isPreparingSensorDialog = true;
    this.managementError = null;
    this.managementSuccess = null;

    forkJoin({
      sensorTypes: this._dataService.getSensorTypes(),
      configuratieBronnen: this._dataService.getSensorConfiguratieBronnen(this.kunstwerkId),
    }).subscribe({
      next: ({ sensorTypes, configuratieBronnen }) => {
        this.isPreparingSensorDialog = false;
        const dialogRef = this._dialog.open(CreateSensorDialogComponent, {
          width: '680px',
          data: {
            kunstwerkNaam: this.currentTree?.kunstwerkdetail.kunstwerk.naam ?? '',
            selectedOnderdeel,
            sensorTypes,
            configuratieBronnen,
          },
        });

        dialogRef.afterClosed().subscribe((request) => {
          if (!request || this.isSavingSensor) {
            return;
          }

          this.isSavingSensor = true;
          this._dataService.createSensorForOnderdeel(this.kunstwerkId, selectedOnderdeel.id, request).subscribe({
            next: (sensor) => {
              this.managementSuccess = `Sensor #${sensor.id} is aangemaakt voor ${selectedOnderdeel.naam}.`;
              this.loadTree(selectedOnderdeel.id);
            },
            error: () => {
              this.managementError = 'Sensor aanmaken is niet gelukt.';
            },
            complete: () => {
              this.isSavingSensor = false;
            },
          });
        });
      },
      error: () => {
        this.isPreparingSensorDialog = false;
        this.managementError = 'Sensorformulier voorbereiden is niet gelukt.';
      },
    });
  }

  public openDeleteOnderdeelDialog(): void {
    if (!this.selectedOnderdeel || this.isDeleting) {
      return;
    }

    const selectedOnderdeel = this.selectedOnderdeel;
    const dialogRef = this._dialog.open(DeleteDataDialogComponent, {
      width: '560px',
      data: {
        title: `Onderdeel "${selectedOnderdeel.naam}" verwijderen`,
        targetLabel: selectedOnderdeel.naam,
        targetType: 'onderdeel',
        message: 'Alle onderliggende onderdelen en sensoren worden meegenomen.',
      },
    });

    dialogRef.afterClosed().subscribe((request) => {
      if (!request || this.isDeleting) {
        return;
      }

      this.isDeleting = true;
      this.managementError = null;
      this.managementSuccess = null;

      this._dataService.deleteOnderdeel(this.kunstwerkId, selectedOnderdeel.id, request).subscribe({
        next: () => {
          this.managementSuccess = request.preserveSensorData
            ? `Onderdeel "${selectedOnderdeel.naam}" is gemarkeerd als verwijderd.`
            : `Onderdeel "${selectedOnderdeel.naam}" en alle gekoppelde data zijn definitief verwijderd.`;
          this.loadTree(request.preserveSensorData ? selectedOnderdeel.id : undefined);
        },
        error: () => {
          this.managementError = 'Onderdeel verwijderen is niet gelukt.';
        },
        complete: () => {
          this.isDeleting = false;
        },
      });
    });
  }

  public openDeleteSensorDialog(sensor: SensorDetailResponse, onderdeel: Onderdeel): void {
    if (this.isDeleting) {
      return;
    }

    const dialogRef = this._dialog.open(DeleteDataDialogComponent, {
      width: '560px',
      data: {
        title: `Sensor #${sensor.id} verwijderen`,
        targetLabel: `Sensor #${sensor.id}`,
        targetType: 'sensor',
        message: 'Deze sensor is gekoppeld aan metingen, afwijkingen en configuratie.',
      },
    });

    dialogRef.afterClosed().subscribe((request) => {
      if (!request || this.isDeleting) {
        return;
      }

      this.isDeleting = true;
      this.managementError = null;
      this.managementSuccess = null;

      this._dataService.deleteSensor(sensor.id, request).subscribe({
        next: () => {
          this.managementSuccess = request.preserveSensorData
            ? `Sensor #${sensor.id} is gemarkeerd als verwijderd.`
            : `Sensor #${sensor.id} en alle gekoppelde data zijn definitief verwijderd.`;
          this.loadTree(onderdeel.id);
        },
        error: () => {
          this.managementError = 'Sensor verwijderen is niet gelukt.';
        },
        complete: () => {
          this.isDeleting = false;
        },
      });
    });
  }

  private findOnderdeelById(onderdelen: Onderdeel[], onderdeelId: number): Onderdeel | null {
    for (const onderdeel of onderdelen ?? []) {
      if (onderdeel.id === onderdeelId) {
        return onderdeel;
      }

      const child = this.findOnderdeelById(onderdeel.onderdelen, onderdeelId);
      if (child) {
        return child;
      }
    }

    return null;
  }

  private fetchAndAttachSensors(onderdeel: Onderdeel): Observable<SensorDetailResponse[]> {
    if (!onderdeel.sensoren || onderdeel.sensoren.length === 0) {
      return of([]);
    }

    return this._dataService.getKunstwerkenSensorenBulk(this.kunstwerkId, onderdeel.sensoren).pipe(
      tap((sensorDetailsArray) => {
        onderdeel.sensoren_details = [];
        onderdeel.sensoren_details = sensorDetailsArray;
      }),
      catchError((err) => {
        console.error(`Failed to fetch sensors for onderdeel ${onderdeel.id}`, err);
        return of([]);
      })
    );
  }

  private refreshFunction(){
    this.refreshDataTimer$ = timer(0,1000).pipe(
      map(tick => 10 - tick),
      take(11),
      tap(value =>{
        if (value === 0 ){
          if (this.selectedOnderdeel) this.selectOnderdeel(this.selectedOnderdeel)
        }
      })
    )
    this._cdr.detectChanges();
  }
}
