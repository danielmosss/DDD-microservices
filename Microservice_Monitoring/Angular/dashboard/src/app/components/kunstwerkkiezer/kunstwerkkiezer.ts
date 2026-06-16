import { Component, OnInit } from '@angular/core';
import { Kunstwerk } from '../../models/types';
import { DataService } from '../../../data.service';
import { Router, RouterModule } from '@angular/router';
import { DatePipe, NgIf, NgFor, AsyncPipe } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-kunstwerkkiezer',
  standalone: true,
  imports: [DatePipe, NgIf, NgFor, MatButtonModule, RouterModule, AsyncPipe],
  templateUrl: './kunstwerkkiezer.html',
  styleUrl: './kunstwerkkiezer.scss',
})
export class Kunstwerkkiezer implements OnInit {
  public selectedKunstwerk: Kunstwerk | null = null;
public kunstwerken$: Observable<Kunstwerk[]> | null = null;
  constructor(
    private _dataService: DataService,
    private _router: Router,
  ) {}


  ngOnInit(): void {
    this.kunstwerken$ = this._dataService.getKunstwerken();
  }

  onSelect(kw: Kunstwerk): void {
    this.selectedKunstwerk = kw;
  }

  navigate(kwId: number): void {
    this._router.navigate([`kunstwerk/${kwId}`]);
  }
}
