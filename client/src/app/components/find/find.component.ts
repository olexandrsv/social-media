import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatDialog } from '@angular/material/dialog';
import { ErrorComponent } from '../error/error.component';
import { PostsStoreService } from 'src/app/services/posts-store.service';
import { newUser } from 'src/app/objects/user';
import { ErrorWindow } from 'src/app/handlers/error';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatBadgeModule } from '@angular/material/badge';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';

export interface User{
	login: string
	id: number
}

export interface GetUserResponse{
	id: number
	login: string
	name: string
	surname: string
	interests: string
	bio: string
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, MatSelectModule, MatAutocompleteModule, MatBadgeModule,
	ReactiveFormsModule],
  selector: 'app-find',
  templateUrl: './find.component.html',
  styleUrls: ['./find.component.css']
})
export class FindComponent {
    options: User[] = [];
    text: string = '';
    login: string = '';
    firstName: string = '';
    secondName: string = '';
    bio: string = '';
    interests: string = '';
	followingID: number = -1;
	postsStore: PostsStoreService
    
    constructor(private http: HttpClient, public dialog: MatDialog, postsStore: PostsStoreService){
		this.postsStore = postsStore;
	}
    
    onButtonPressed(event: KeyboardEvent){
    	if(event.key !== "Enter"){			
	    	this.http.get<User[]>(`http://localhost:8081/users?info=`+this.text, {withCredentials: true}).subscribe({
				next: response => {
					console.log(response);
					const id = localStorage.getItem("id")
					this.options = response.filter(u => u.id.toString() !== id);
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
			});
		}else{
			this.onOptionChoosed();
		}
    }

	onOptionClicked(login: string){
		this.text = login;
		this.onOptionChoosed();
	}

	onOptionChoosed(){
		let id = -1;
		for (let i=0; i<this.options.length; i++){
			if (this.options[i].login == this.text){
				id = this.options[i].id;
			}
		}
		if (id === -1){
			return
		}
		this.http.get<GetUserResponse>(`http://localhost:8081/users/`+id, {withCredentials: true}).subscribe({
			next: response => {
				this.followingID = response.id;
				this.login = response.login;
				this.firstName = response.name;
				this.secondName = response.surname;
				this.bio = response.bio;
				this.interests = response.interests;
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
		});
	}
    
    async follow(){
		if (this.followingID === -1){
			return
		}
    			
		let err = await this.postsStore.subscribe(newUser(this.followingID, this.login))
		if (err){
			new ErrorWindow(err, this.dialog)
			return
		}
    }
}