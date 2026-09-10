import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { HttpClient } from '@angular/common/http';
import { Comment } from "src/app/objects/comment";
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';

export interface CommentData{
  comment?: Comment
  title: string
  buttonTitle: string
  btnOnClick: (c: Comment, dialogRef: MatDialogRef<CommentComponent>) => void;
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, RouterModule],
  selector: 'app-comment',
  templateUrl: './comment.component.html',
  styleUrls: ['./comment.component.css']
})

export class CommentComponent {
  comment: Comment = new Comment();
  
  commentData: CommentData;

  constructor(
    public dialogRef: MatDialogRef<CommentComponent>,
    private http: HttpClient,  
    @Inject(MAT_DIALOG_DATA) public data: CommentData
  ){
    this.commentData = data;
    if (data.comment) {
      this.comment = data.comment;
    }
  }

  onImgSelected(event: any, input: HTMLInputElement){
    this.comment.onImgSelected(event);
    input.value = '';
  }

  onFileSelected(event: any, input: HTMLInputElement){
    this.comment.onFileSelected(event);
    input.value = '';
  }

  removeImage(i: number){
    this.comment.onImgRemoved(i);
  }

  removeFile(i: number){
    this.comment.onFileRemoved(i);
  }

  onClick(){
    this.comment.parseText();
    this.commentData.btnOnClick(this.comment, this.dialogRef);
  }
}
