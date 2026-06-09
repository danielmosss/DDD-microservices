export interface SensorConfiguratie {
  id: number;
  sensorId: number;
  minValue: number | null;
  maxValue: number | null;
  margePercentage: number | null;
}

export interface Sensor {
  id: number;
  kunstwerkId: number;
  onderdeelId: number | null;
  geolocation: string | null;
  sensorTypeId: number;
  lastAnalyzedMetingId: number | null;
  sensorConfiguratie: SensorConfiguratie;
  laatsteMeting?: number; // voor updates live
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

// --- Toegevoegd voor de Boomstructuur ---

export interface Onderdeel {
  id: number;
  naam: string;
  kunstwerkId: number;
  parentOnderdeelId: number | null; // null als het direct onder het kunstwerk valt

  // De geneste children voor de UI:
  onderdelen: Onderdeel[];
  sensoren: Sensor[];
}

export interface KunstwerkTreeView {
  kunstwerk: Kunstwerk;
  onderdelen: Onderdeel[]; // De hoogste niveau onderdelen
  losseSensoren: Sensor[]; // Sensoren die direct aan het kunstwerk hangen, niet aan een onderdeel
}
