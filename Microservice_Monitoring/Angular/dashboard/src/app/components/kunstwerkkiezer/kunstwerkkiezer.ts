import { Component, OnInit } from '@angular/core';
import { Kunstwerk, KunstwerkDHU } from '../../models/types';
import { DataService } from '../../../data.service';
import { Router, RouterModule } from '@angular/router';
import { DatePipe, NgIf, NgFor, AsyncPipe } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { Observable, BehaviorSubject, switchMap, filter } from 'rxjs'; // <-- Voeg RxJS operators toe

@Component({
  selector: 'app-kunstwerkkiezer',
  standalone: true,
  imports: [DatePipe, NgIf, NgFor, MatButtonModule, RouterModule, AsyncPipe],
  templateUrl: './kunstwerkkiezer.html',
  styleUrl: './kunstwerkkiezer.scss',
})
export class Kunstwerkkiezer implements OnInit {
  public kunstwerken$: Observable<Kunstwerk[]> | null = null;
  private _selectedKunstwerk = new BehaviorSubject<Kunstwerk | null>(null);
  public selectedKunstwerk$ = this._selectedKunstwerk.asObservable();
  public kunstwerkDHU$: Observable<KunstwerkDHU> | null = null;

  constructor(
    private _dataService: DataService,
  ) {}

  ngOnInit(): void {
    this.kunstwerken$ = this._dataService.getKunstwerken();
    this.kunstwerkDHU$ = this.selectedKunstwerk$.pipe(
      filter((kw): kw is Kunstwerk => kw !== null),
      switchMap((kw) => this._dataService.getKunstwerkDHU(kw.id))
    );
  }

  onSelect(kw: Kunstwerk): void {
    this._selectedKunstwerk.next(kw);
  }
}
