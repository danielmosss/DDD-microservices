import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { environment } from './environments/environment.local'
import { Router } from '@angular/router';
import { Kunstwerken } from './app/models/types';

@Injectable({
  providedIn: 'root'
})
export class DataService {
  private _hostname = environment.apiUrl;
  private _SecureApi = this._hostname + "/api";

  constructor(private http: HttpClient, private _router: Router) { }

  public getKunstwerken() {
    return this.http.get<Array<Kunstwerken>>(this._SecureApi + `/v1/kunstwerken`, { headers: this.getCustomHeaders() });
  }

  public getCustomHeaders(): HttpHeaders {
    var headers = new HttpHeaders()
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json')
    return headers;
  }
}
