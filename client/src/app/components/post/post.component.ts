import { HttpClient } from '@angular/common/http';
import { Component, Inject } from '@angular/core';
import { Post } from 'src/app/objects/post';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { MatInputModule } from '@angular/material/input';


export interface PostData{
  post?: Post
  title: string
  buttonTitle: string
  btnOnClick: (post: Post, dialogRef: MatDialogRef<PostComponent>) => void;
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, RouterModule],
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.css']
})
export class PostComponent {
  post: Post = new Post();
  postData: PostData;

  constructor(
    public dialogRef: MatDialogRef<PostComponent>,
    private http: HttpClient,  
    @Inject(MAT_DIALOG_DATA) public data: PostData
  ){
    this.postData = data;
    if (data.post) {
      this.post = data.post;
    }
  }

  onImgSelected(event: any, input: HTMLInputElement){
    this.post.onImgSelected(event);
    input.value = '';
  }

  onFileSelected(event: any, input: HTMLInputElement){
    this.post.onFileSelected(event);
    input.value = '';
  }

  removeImage(i: number){
    this.post.onImgRemoved(i);
  }

  removeFile(i: number){
    this.post.onFileRemoved(i);
  }

  onClick(){
    this.postData.btnOnClick(this.post, this.dialogRef);
  }
}
