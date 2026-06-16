import { CommonModule } from '@angular/common';
import { Component, computed, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { startWith } from 'rxjs';
import { CreateSensorRequest, Onderdeel, SensorConfiguratieBron, SensorType } from '../../../../models/types';

interface CreateSensorDialogData {
  kunstwerkNaam: string;
  selectedOnderdeel: Onderdeel;
  sensorTypes: SensorType[];
  configuratieBronnen: SensorConfiguratieBron[];
}

@Component({
  selector: 'app-create-sensor-dialog',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule, MatDialogModule],
  templateUrl: './create-sensor-dialog.html',
  styleUrls: ['./create-sensor-dialog.scss'],
})
export class CreateSensorDialogComponent {
  private _fb = inject(FormBuilder);
  private _dialogRef = inject(MatDialogRef<CreateSensorDialogComponent, CreateSensorRequest>);
  public data = inject<CreateSensorDialogData>(MAT_DIALOG_DATA);

  public error: string | null = null;
  public form = this._fb.group({
    sensorTypeId: [this.data.sensorTypes[0]?.id ?? null, [Validators.required, Validators.min(1)]],
    geolocation: [''],
    configuratieBronSensorId: [null as number | null],
    minValue: [null as number | null, [Validators.required]],
    maxValue: [null as number | null],
    margePercentage: [null as number | null, [Validators.min(0)]],
  });

  private sensorTypeId = toSignal(this.form.controls.sensorTypeId.valueChanges.pipe(startWith(this.form.controls.sensorTypeId.value)));
  public selectedSensorType = computed(() => this.data.sensorTypes.find((sensorType) => sensorType.id === Number(this.sensorTypeId())) ?? null);
  public configuratieBronnenVoorType = computed(() => {
    const selectedSensorType = this.selectedSensorType();
    return selectedSensorType
      ? this.data.configuratieBronnen.filter((bron) => bron.sensorType.id === selectedSensorType.id)
      : this.data.configuratieBronnen;
  });

  public constructor() {
    this.form.controls.sensorTypeId.valueChanges.subscribe(() => {
      this.form.patchValue({ configuratieBronSensorId: null, minValue: null, maxValue: null, margePercentage: null }, { emitEvent: false });
      this.error = null;
    });
  }

  public applyConfiguratieBron(): void {
    const bronSensorId = Number(this.form.controls.configuratieBronSensorId.value);
    const bron = this.data.configuratieBronnen.find((item) => item.sensorId === bronSensorId);
    if (!bron) {
      return;
    }

    this.form.patchValue({
      sensorTypeId: bron.sensorType.id,
      minValue: bron.sensorConfiguratie.min_value,
      maxValue: bron.sensorConfiguratie.max_value,
      margePercentage: bron.sensorConfiguratie.marge_percentage,
    });
  }

  public submit(): void {
    this.error = null;
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.error = 'Vul de verplichte velden in.';
      return;
    }

    const value = this.form.getRawValue();
    const sensorType = this.selectedSensorType();
    const minValue = this.toNumberOrNull(value.minValue);
    const maxValue = this.toNumberOrNull(value.maxValue);
    const margePercentage = this.toNumberOrNull(value.margePercentage);

    if (!sensorType || !value.sensorTypeId) {
      this.error = 'Kies een sensortype.';
      return;
    }

    if (sensorType.drempel_is_range) {
      if (minValue === null || maxValue === null) {
        this.error = 'Range-sensoren hebben een minimale en maximale waarde nodig.';
        return;
      }
      if (minValue >= maxValue) {
        this.error = 'De minimale waarde moet lager zijn dan de maximale waarde.';
        return;
      }
    }

    if (!sensorType.drempel_is_range && minValue === null) {
      this.error = 'Niet-range sensoren hebben een normwaarde nodig.';
      return;
    }

    if (margePercentage !== null && margePercentage < 0) {
      this.error = 'De marge mag niet negatief zijn.';
      return;
    }

    const geolocation = value.geolocation?.trim() || null;
    this._dialogRef.close({
      sensorTypeId: Number(value.sensorTypeId),
      geolocation,
      configuratieBronSensorId: this.toNumberOrNull(value.configuratieBronSensorId),
      configuratie: {
        minValue,
        maxValue: sensorType.drempel_is_range ? maxValue : null,
        margePercentage,
      },
    });
  }

  private toNumberOrNull(value: number | string | null | undefined): number | null {
    if (value === null || value === undefined || value === '') {
      return null;
    }

    const parsedValue = Number(value);
    return Number.isNaN(parsedValue) ? null : parsedValue;
  }
}
