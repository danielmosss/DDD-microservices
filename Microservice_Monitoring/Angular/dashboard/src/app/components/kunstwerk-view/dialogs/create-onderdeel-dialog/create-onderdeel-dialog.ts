import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { CreateOnderdeelRequest, Onderdeel } from '../../../../models/types';

interface CreateOnderdeelDialogData {
  kunstwerkNaam: string;
  selectedOnderdeel: Onderdeel | null;
}

@Component({
  selector: 'app-create-onderdeel-dialog',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule, MatDialogModule],
  templateUrl: './create-onderdeel-dialog.html',
  styleUrls: ['./create-onderdeel-dialog.scss'],
})
export class CreateOnderdeelDialogComponent {
  private _fb = inject(FormBuilder);
  private _dialogRef = inject(MatDialogRef<CreateOnderdeelDialogComponent, CreateOnderdeelRequest>);
  public data = inject<CreateOnderdeelDialogData>(MAT_DIALOG_DATA);

  public form = this._fb.nonNullable.group({
    naam: ['', [Validators.required, Validators.maxLength(255)]],
    plaatsing: ['selected' as 'selected' | 'root'],
  });

  public submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const value = this.form.getRawValue();
    this._dialogRef.close({
      naam: value.naam.trim(),
      parentOnderdeelId: value.plaatsing === 'selected' ? this.data.selectedOnderdeel?.id ?? null : null,
    });
  }
}
