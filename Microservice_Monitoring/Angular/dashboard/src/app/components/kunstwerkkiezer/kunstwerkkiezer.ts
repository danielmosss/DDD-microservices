import { Component } from '@angular/core';
import { Kunstwerken } from '../../models/types';
import { DataService } from '../../../data.service';

@Component({
  selector: 'app-kunstwerkkiezer',
  imports: [],
  templateUrl: './kunstwerkkiezer.html',
  styleUrl: './kunstwerkkiezer.scss',
})
export class Kunstwerkkiezer {
  constructor(private _dateservice: DataService){}

  kunstwerken: Kunstwerken[] = []

  ngOnInit(): void {
    this._dateservice.getKunstwerken().subscribe(kunstwerken => {
      this.kunstwerken = kunstwerken
    })
  }
}
