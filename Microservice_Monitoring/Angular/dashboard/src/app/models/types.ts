export interface SensorConfiguratie {
  id: number;
  sensorId: number;
  minValue: number | null;
  maxValue: number | null;
  margePercentage: number | null;
}

export interface Kunstwerk {
  id: number;
  beheerIdentifier: string;
  naam: string;
  geolocation: string | null;
  kunstwerkTypeId: number | null;
  beschrijving: string | null;
  deleted: boolean;
  lastsenddhupdate: string | null; // time.Time komt als ISO string uit JSON
}

export interface KunstwerkDHU {
  kunstwerkId: number;
  status: 'healthy' | 'warning' | 'critical' | 'offline';
  aantalSensoren: number;
  aantalActieveSensoren: number;
  aantalAfwijkendeSensoren: number;
  aantalAfwijkingen: number;
}

// --- Toegevoegd voor de Boomstructuur ---

export interface Onderdeel {
  id: number;
  naam: string;
  kunstwerkId: number;
  parentOnderdeelId: number | null;

  // De geneste children voor de UI:
  onderdelen: Onderdeel[];
  sensoren: number[];
  sensoren_details: SensorDetailResponse[];
}

export interface KunstwerkType {
  id: number;
  beschrijving: string;
  naam: string;
}

export interface KunstwerkDetail {
  kunstwerk: Kunstwerk;
  KunstwerkType: KunstwerkType;
}

export interface KunstwerkTreeView {
  kunstwerkdetail: KunstwerkDetail;
  onderdelen: Onderdeel[];
  losseSensoren: number[];
  sensoren_details: SensorDetailResponse[];
}

export interface SensorDetailResponse {
  afwijking: {
    gemeten_waarde: number;
    id: number;
    is_warning: boolean;
    kunstwerkId: number;
    metingId: number;
    norm_marge_percentage: number;
    norm_max_waarde: number;
    norm_min_waarde: number;
    sensorId: number;
    time: string;
  };
  id: number;
  laatsteMeting: {
    id: number;
    inspectieId: string;
    isHandmatig: boolean;
    kunstwerkId: number;
    sensorId: number;
    time: string;
    waarde: number;
  };
  sensorConfiguratie: {
    id: number;
    marge_percentage: number;
    max_value: number;
    min_value: number;
    sensor_id: number;
  };
  sensorType: { drempel_is_range: boolean; eenheid: string; id: number; naam: string };
  status: string;
}
