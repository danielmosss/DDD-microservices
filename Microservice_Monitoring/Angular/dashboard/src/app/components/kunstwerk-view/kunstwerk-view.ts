import { ChangeDetectorRef, Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KunstwerkTreeView, Onderdeel, SensorDetailResponse } from '../../models/types';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { map, Observable, take, tap, timer } from 'rxjs';
import { DataService } from '../../../data.service';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-kunstwerk-view',
  standalone: true,
  imports: [CommonModule, RouterModule, MatButtonModule],
  templateUrl: './kunstwerk-view.html',
  styleUrls: ['./kunstwerk-view.scss'],
})
export class KunstwerkViewComponent implements OnInit {
  public kunstwerkId: number = 0;
  private activatedRoute = inject(ActivatedRoute);
  public treeView$: Observable<KunstwerkTreeView> | null = null;
  public refreshDataTimer$: Observable<number> | null = null;

  private _cdr = inject(ChangeDetectorRef);

  constructor(private _dataService: DataService) {}

  ngOnInit(): void {
    this.kunstwerkId = parseInt(this.activatedRoute.snapshot.paramMap.get('kunstwerkid') ?? '0');
    this.treeView$ = this._dataService.getKunstwerkTree(this.kunstwerkId).pipe(
      tap((treeData) => {
        if (treeData && treeData.onderdelen && treeData.onderdelen.length > 0) {
          this.selectOnderdeel(treeData.onderdelen[0]);
        }
      })
    );

  }

  selectedOnderdeel: Onderdeel | null = null;

  selectOnderdeel(onderdeel: Onderdeel) {
    this.selectedOnderdeel = onderdeel;
    this.fetchAndAttachSensors(onderdeel);
    for (let index = 0; index < onderdeel.onderdelen?.length; index++) {
      const subOnderdeel = onderdeel.onderdelen[index];
      this.fetchAndAttachSensors(subOnderdeel);
    }
    setTimeout(() => {
      this.refreshFunction()
      this._cdr.detectChanges();
    }, 250);
  }

  private fetchAndAttachSensors(onderdeel: Onderdeel): void {
    if (!onderdeel.sensoren || onderdeel.sensoren.length === 0) {
      return;
    }

    this._dataService.getKunstwerkenSensorenBulk(this.kunstwerkId, onderdeel.sensoren).subscribe({
      next: (sensorDetailsArray) => {
        onderdeel.sensoren_details = [];
        onderdeel.sensoren_details = sensorDetailsArray;
        this._cdr.detectChanges();
      },
      error: (err) => console.error(`Failed to fetch sensors for onderdeel ${onderdeel.id}`, err),
    });
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
  }
}
