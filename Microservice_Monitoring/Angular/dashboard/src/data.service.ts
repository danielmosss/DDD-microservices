import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { environment } from './environments/environment.local'
import { Router } from '@angular/router';
import {
  CreateOnderdeelRequest,
  CreateSensorRequest,
  DeleteDataRequest,
  Kunstwerk,
  KunstwerkDHU,
  KunstwerkTreeView,
  Onderdeel,
  SensorConfiguratie,
  SensorConfiguratieBron,
  SensorDetailResponse,
  SensorType,
  UpdateSensorConfiguratieRequest,
} from './app/models/types';

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

  public createOnderdeel(kunstwerkId: number, request: CreateOnderdeelRequest) {
    return this.http.post<Onderdeel>(this._SecureApi + `/v1/kunstwerken/${kunstwerkId}/onderdelen`, request, { headers: this.getCustomHeaders() });
  }

  public createSensorForOnderdeel(kunstwerkId: number, onderdeelId: number, request: CreateSensorRequest) {
    return this.http.post<SensorDetailResponse>(this._SecureApi + `/v1/kunstwerken/${kunstwerkId}/onderdelen/${onderdeelId}/sensoren`, request, { headers: this.getCustomHeaders() });
  }

  public deleteOnderdeel(kunstwerkId: number, onderdeelId: number, request: DeleteDataRequest) {
    return this.http.delete<{ deleted: boolean; preserveSensorData: boolean }>(this._SecureApi + `/v1/kunstwerken/${kunstwerkId}/onderdelen/${onderdeelId}`, { headers: this.getCustomHeaders(), body: request });
  }

  public deleteSensor(sensorId: number, request: DeleteDataRequest) {
    return this.http.delete<{ deleted: boolean; preserveSensorData: boolean }>(this._SecureApi + `/v1/sensoren/${sensorId}`, { headers: this.getCustomHeaders(), body: request });
  }

  public getSensorTypes() {
    return this.http.get<SensorType[]>(this._SecureApi + `/v1/sensortypes`, { headers: this.getCustomHeaders() });
  }

  public getSensorConfiguratieBronnen(kunstwerkId: number, sensorTypeId?: number) {
    let params = new HttpParams();
    if (sensorTypeId) {
      params = params.set('sensorTypeId', sensorTypeId);
    }

    return this.http.get<SensorConfiguratieBron[]>(this._SecureApi + `/v1/kunstwerken/${kunstwerkId}/sensor-configuratie-bronnen`, { headers: this.getCustomHeaders(), params });
  }

  public updateSensorConfiguratie(sensorId: number, request: UpdateSensorConfiguratieRequest) {
    return this.http.put<SensorConfiguratie>(this._SecureApi + `/v1/sensoren/${sensorId}/configuratie`, request, { headers: this.getCustomHeaders() });
  }

  public getCustomHeaders(): HttpHeaders {
    var headers = new HttpHeaders()
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json')
    return headers;
  }
}
