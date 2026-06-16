import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { DeleteDataRequest } from '../../../../models/types';

interface DeleteDataDialogData {
  title: string;
  targetLabel: string;
  targetType: 'onderdeel' | 'sensor';
  message: string;
}

@Component({
  selector: 'app-delete-data-dialog',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule, MatDialogModule],
  templateUrl: './delete-data-dialog.html',
  styleUrls: ['./delete-data-dialog.scss'],
})
export class DeleteDataDialogComponent {
  private _fb = inject(FormBuilder);
  private _dialogRef = inject(MatDialogRef<DeleteDataDialogComponent, DeleteDataRequest>);
  public data = inject<DeleteDataDialogData>(MAT_DIALOG_DATA);

  public form = this._fb.nonNullable.group({
    preserveSensorData: [true],
  });

  public submit(): void {
    this._dialogRef.close(this.form.getRawValue());
  }
}
