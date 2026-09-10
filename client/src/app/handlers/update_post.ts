import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { PostComponent, PostData } from "../components/post/post.component";
import { HttpClient } from "@angular/common/http";
import { Post } from "../objects/post";
import { GetPostResponse } from "../components/posts/posts.component";

export class UpdatePost{
    postData: PostData;
    dialogRef: MatDialogRef<PostComponent, any>;
    
    constructor(public dialog: MatDialog, http: HttpClient, post: Post, updatePost: (p: Post)=>void){
        this.postData = {
            post: post,
            title: "Update post",
            buttonTitle: "Update",
            btnOnClick: this.update(http, updatePost),
        }
        this.dialogRef = this.dialog.open(PostComponent, {
            data: this.postData,
            height: '400px',
            width: '600px',
        })
    }

    update(http: HttpClient, updatePost: (p: Post)=>void): (post:Post, dialogRef: MatDialogRef<PostComponent>) => void {
        return (post:Post, dialogRef: MatDialogRef<PostComponent>) => {
            const form = post.generateUpdateForm();
            http.put<GetPostResponse>(`http://localhost:8082/users/posts/`+post.id, form, {withCredentials: true}).subscribe({
                next: result => {
                    console.log(result);
                    const post = new Post();
                    post.parseHttpResponse(result);
                    dialogRef.close();
                    updatePost(post);
                },
                error: err => {
                    
                }
            })
        }
    }
}