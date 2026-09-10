import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

export interface ErrorMessage {
	code: string
  message: string
  statusText: string
}

@Component({
  selector: 'app-error',
  templateUrl: './error.component.html',
  styleUrls: ['./error.component.css']
})
export class ErrorComponent {
  code: string = "";
  messages: string[] = [];
  statusText: string = "";

  constructor(
    public dialogRef: MatDialogRef<ErrorComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ErrorMessage,
  ){
    this.code = data.code;
    this.statusText = data.statusText;
    if (data.message){
      this.messages = data.message.split("\n");
    }
  }
}
