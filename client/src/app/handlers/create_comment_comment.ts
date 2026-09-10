import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { CommentComponent, CommentData } from "../components/comment/comment.component";
import { Comment } from "src/app/objects/comment";
import { HttpClient } from "@angular/common/http";
import { CommentModel } from "../components/posts/posts.component";

export class CreateCommentComment{
    commentData: CommentData;
    
    constructor(
        public dialog: MatDialog, 
        http: HttpClient, 
        parentID: string,
        addComment: (c: Comment) => void)
    {
        this.commentData = {
            title: "Create comment",
            buttonTitle: "Create",
            btnOnClick: this.create(http, parentID, addComment),
        }
        const dialogRef = this.dialog.open(CommentComponent, {
            data: this.commentData,
            height: '400px',
            width: '600px',
        })
    }

    create(http: HttpClient, parentID: string, addComment: (c: Comment) => void): 
        (c: Comment, dialogRef: MatDialogRef<CommentComponent>)=> void 
    {
        return (c: Comment, dialogRef: MatDialogRef<CommentComponent>)=>{
            const form = c.generateCreateForm();
            
            http.post<CommentModel>(
                `http://localhost:8082/users/posts/comments/${parentID}/comments`, 
                form, 
                {withCredentials: true}
            ).subscribe({
                next: result => {
                    console.log(result);
                    const comment = new Comment()
                    comment.parseHttpResponse(result)
                    addComment(comment);
                    dialogRef.close();
                },
                error: err => {
    
                }
            })
        }
    }
}