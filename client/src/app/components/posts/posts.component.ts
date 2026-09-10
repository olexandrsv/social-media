import { Component, Inject } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import {MatDialog, MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import { Post } from 'src/app/objects/post';
import { CreatePost } from 'src/app/handlers/create_post';
import { UpdatePost } from 'src/app/handlers/update_post';
import { ErrorComponent } from '../error/error.component';
import { CreatePostComment } from 'src/app/handlers/create_post_comment';
import { Comment } from 'src/app/objects/comment';
import { UpdatePostComment } from 'src/app/handlers/update_post_comment';
import { CreateCommentComment } from 'src/app/handlers/create_comment_comment';
import { ErrorWindow } from 'src/app/handlers/error';
import { WebsocketService } from 'src/app/services/websocket.service';
import { Tone } from 'src/app/objects/tone';
import { BiResult, ClientService, Result } from 'src/app/services/client.service';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatBadgeModule } from '@angular/material/badge';
import { DecimalPipe } from '@angular/common';


export interface GetPostResponse{
	id: string
	user_id: number
	text: string
	files: string[]
	images: string[]
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, RouterModule],
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.css']
})
export class PostsComponent {
	text: string = '';
	images: string = '';
	selectedImgs = [];
	selectedFiles = [];
	
	posts: Post[] = [];
	
	constructor(public http: HttpClient, public dialog: MatDialog){
		const options = {
			withCredentials: true
		};	
		const id = localStorage.getItem("id");
		http.get<GetPostResponse[]>(`http://localhost:8082/users/${id}/posts`, options).subscribe((response) => {
			for (let i=0; i<response.length; i++){
				const post = new Post();
				post.parseHttpResponse(response[i]);
				post.parseText();
				this.posts.push(post);
			}

			console.log(response)
			console.log(this.posts);
		});
	}

	deletePost(i: number){
		const id = this.posts[i].id;
		this.http.delete(`http://localhost:8082/users/posts/${id}`, {withCredentials: true}).subscribe({
			next: response => {
				this.posts = this.posts.filter((el, index) => {
					return index !== i
				})
			},
			error: err => {
				const dialogRef = this.dialog.open(ErrorComponent, {
					data: {
						code: err.status,
						message: err.error,
						statusText: err.statusText,
					},
					height: '400px',
					width: '800px',
				})
			}
		})
		
	}

	updatePost(i: number){
		let post = new Post();
		post.copy(this.posts[i]);
		new UpdatePost(this.dialog, this.http, post, (post: Post) => {
			post.parseText();
			this.posts[i] = post;
		});
	}

	createPost(){
		new CreatePost(this.dialog, this.http, (post: Post)=>{
			post.parseText();
			this.addPost(post);
		});
	}

	addPost(post: Post){
		this.posts.push(post);
	}
	
	post(){
		const postData = new FormData();
		postData.append('text',  this.text);
		for (let file of this.selectedImgs) {
        	postData.append('[]images', file);
    	}
		for (let file of this.selectedFiles) {
        	postData.append('[]files', file);
    	}

		const options = {
			withCredentials: true
		};

		this.http.post('http://localhost:8082/post', postData, options).subscribe((response: any) => {
			console.log(response);
			if (this.posts == null){
				this.posts = [];
			}
			this.posts.push(response);
			console.log(this.posts);
		});
	}
	
	onFileSelected(event: any){
		this.selectedFiles = event.target.files;
	}
	
	onImgSelected(event: any){
		this.selectedImgs = event.target.files;
	}
	
	openDialog(i: number): void {
	    const dialogRef = this.dialog.open(CommentsWindow, {
	        data: {postId: this.posts[i].id},
			height: '600px',
  			width: '700px',
	    });
	
	    dialogRef.afterClosed().subscribe(result => {
	      console.log('result');
	    });
    }
}

export interface GetCommentsReq{
	tone: ToneModel
	comments: CommentModel[]
}

export interface ToneModel{
	positive_percentage: number
	negative_percentage: number
}

export interface CommentModel{
	id: string
	user_id: number
	user_name: string
	user_surname: string
	text: string
	files: string[]
	images: string[]
}

enum Parent{
	Post,
	Comment
}

@Component({
	imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
		MatInputModule, MatSelectModule, MatAutocompleteModule, MatBadgeModule,
		ReactiveFormsModule, DecimalPipe],
	selector: 'comment',
	templateUrl: './comment.html',
	styleUrls: ['./comment.css']
})

export class CommentsWindow {
	text: string = '';
	selectedImgs = [];
	selectedFiles = [];
	postId: string = "";
	commentId: string = "";
	parent: Parent = Parent.Post;
	comments: Comment[] = [];
	tone: Tone = new Tone(0, 0);
	
	constructor(
	  private http: HttpClient,
	  public dialogRef: MatDialogRef<CommentsWindow>,
	  public dialog: MatDialog,
	  public client: ClientService,
	  @Inject(MAT_DIALOG_DATA) public data: any,
	) {
		this.init(data)
	}

	async init(data: any){
		if (data.postId){
			this.postId = data.postId;
		}

		if (data.commentId){
			this.commentId = data.commentId;
			this.parent = Parent.Comment;
		}

		let result: BiResult<Comment[], Tone>
		if (this.parent === Parent.Comment){
			result = await this.client.getCommentComments(this.commentId)
		} else {
			result = await this.client.getPostComments(this.postId)
		}
		const [comments, tone, err] = result
		if (err){
			new ErrorWindow(err, this.dialog)
			return
		}
		this.tone = tone
		this.comments = comments		
	}

	deleteComment(i: number){
		if (this.parent === Parent.Post){
			this.deletePostComment(i);
		} else if (this.parent === Parent.Comment){
			this.deleteCommentComment(i);
		}
	}

	updateComment(i: number){
		const comment = new Comment();
		comment.copy(this.comments[i]);
		this.updatePostComment(i, comment);
	}

	createComment(){
		if (this.parent == Parent.Post){
			this.createPostComment();
		} else if (this.parent == Parent.Comment){
			this.createCommentComment();
		}
	}

	createPostComment(){
		new CreatePostComment(this.dialog, this.http, this.postId, (c: Comment) => {
			c.parseText();
			this.comments.push(c);
		})
	}

	createCommentComment(){
		new CreateCommentComment(this.dialog, this.http, this.commentId, (c: Comment) => {
			c.parseText();
			this.comments.push(c);
		})
	}

	async deletePostComment(i: number){
		const comment = this.comments[i]
		const err = await this.client.deletePostComment(this.postId, comment.id)
		if (err){
			new ErrorWindow(err, this.dialog)
		}
		this.comments = this.comments.filter((value, index) => {
			return index !== i
		})
	}

	async deleteCommentComment(i: number){
		const comment = this.comments[i]
		const err = await this.client.deleteCommentComment(this.postId, comment.id)
		if (err){
			new ErrorWindow(err, this.dialog)
		}
		this.comments = this.comments.filter((value, index) => {
			return index !== i
		})
	}

	updatePostComment(i: number, comment: Comment){
		new UpdatePostComment(this.dialog, this.http, comment, (c: Comment) => {
			c.parseText();
			this.comments[i] = c;
		})
	}

	openComments(i: number){
		const dialogRef = this.dialog.open(CommentsWindow, {
	        data: {commentId: this.comments[i].id},
			height: '600px',
  			width: '700px',
	    });
	}
	
	onNoClick(): void {
	  this.dialogRef.close();
	}
	
	post(){
		// const postData = new FormData();
		// postData.append('postId', this.postId);
		// postData.append('text',  this.text);
		// for (let file of this.selectedImgs) {
        // 	postData.append('images[]', file);
    	// }
		// for (let file of this.selectedFiles) {
        // 	postData.append('files[]', file);
    	// }
		
		// this.http.post('http://localhost:8080/comment', postData).subscribe((response: any) => {
		// 	console.log(this.posts);
		// 	if(this.posts == null){
		// 		this.posts = [];
		// 		console.log("here1");
		// 	}
		// 	this.posts.push(response);
		// 	console.log(this.posts);
		// });
		
	}
	
	onFileSelected(event: any){
		this.selectedFiles = event.target.files;
	}
	
	onImgSelected(event: any){
		this.selectedImgs = event.target.files;
	}
}
