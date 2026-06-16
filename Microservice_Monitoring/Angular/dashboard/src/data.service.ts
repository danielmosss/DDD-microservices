import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { environment } from './environments/environment.local'
import { Router } from '@angular/router';
import { Kunstwerk, KunstwerkDHU, KunstwerkTreeView, SensorDetailResponse } from './app/models/types';

@Injectable({
  providedIn: 'root'
})
export class DataService {
  private _hostname = environment.apiUrl;
  private _SecureApi = this._hostname + "/api";

  constructor(private http: HttpClient, private _router: Router) { }

  public getKunstwerken() {
    return this.http.get<Array<Kunstwerk>>(this._SecureApi + `/v1/kunstwerken`, { headers: this.getCustomHeaders() });
  }

  public getKunstwerkDHU(kunstwerkId: number) {
    return this.http.get<KunstwerkDHU>(this._SecureApi + `/v1/kunstwerken/${kunstwerkId}/dailyhealthupdate`, { headers: this.getCustomHeaders() });
  }

  public getKunstwerkTree(kunstwerkId: number) {
    return this.http.get<KunstwerkTreeView>(this._SecureApi + `/frontend/kunstwerken/${kunstwerkId}/tree`, { headers: this.getCustomHeaders() });
  }

  public getKunstwerkenSensorenBulk(kunstwerkId: number, sensoren: Array<number>){
    return this.http.post<SensorDetailResponse[]>(this._SecureApi + `/frontend/kunstwerken/${kunstwerkId}/sensoren/bulk-actueel`, { sensorIds: sensoren } , { headers: this.getCustomHeaders() });
  }

  public getCustomHeaders(): HttpHeaders {
    var headers = new HttpHeaders()
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json')
    return headers;
  }
}
