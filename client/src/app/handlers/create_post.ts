import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { Post } from "../objects/post";
import { PostComponent, PostData } from "../components/post/post.component";
import { HttpClient, HttpHandler, HttpXhrBackend } from "@angular/common/http";
import { GetPostResponse } from "../components/posts/posts.component";

export class CreatePost{
    postData: PostData;

    constructor(public dialog: MatDialog, http: HttpClient, addPost: (post: Post) => void){
        this.postData = {
            title: "Create post",
            buttonTitle: "Create",
            btnOnClick: this.create(http, addPost),
        }
        const dialogRef = this.dialog.open(PostComponent, {
            data: this.postData,
            height: '400px',
            width: '600px',
        })
    }

   create(http: HttpClient, addPost: (post: Post) => void): (post: Post, dialogRef: MatDialogRef<PostComponent>)=> void{
    return (post: Post, dialogRef: MatDialogRef<PostComponent>) => {
        const form = post.generateCreateForm();

        http.post<GetPostResponse>(`http://localhost:8082/users/posts`, form, {withCredentials: true}).subscribe({
            next: result => {
                console.log(result);
                const post = new Post()
                post.parseHttpResponse(result)
                addPost(post);
                dialogRef.close();
            },
            error: err => {

            }
            })
        }
    }
}