import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { CommentComponent, CommentData } from "../components/comment/comment.component";
import { Comment } from "src/app/objects/comment";
import { HttpClient } from "@angular/common/http";
import { CommentModel } from "../components/posts/posts.component";

export class UpdatePostComment{
    commentData: CommentData;
    
    constructor(
        public dialog: MatDialog, 
        http: HttpClient,
        comment: Comment, 
        updateComment: (c: Comment) => void)
    {
        this.commentData = {
            title: "Update comment",
            buttonTitle: "Update",
            btnOnClick: this.update(http, updateComment),
            comment: comment,
        }
        const dialogRef = this.dialog.open(CommentComponent, {
            data: this.commentData,
            height: '400px',
            width: '600px',
        })
    }

    update(
        http: HttpClient,
        updateComment: (c: Comment) => void): (c: Comment, dialogRef: MatDialogRef<CommentComponent>)=> void
    {
        return (c: Comment, dialogRef: MatDialogRef<CommentComponent>)=>{
            const form = c.generateUpdateForm();
            
            http.put<CommentModel>(
                `http://localhost:8082/users/posts/comments/${c.id}`, 
                form, 
                {withCredentials: true}
            ).subscribe({
                next: result => {
                    console.log(result);
                    const comment = new Comment()
                    comment.parseHttpResponse(result)
                    updateComment(comment);
                    dialogRef.close();
                },
                error: err => {
    
                }
            })
        }
    }
}