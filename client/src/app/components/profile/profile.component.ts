import { Component } from '@angular/core';
import { HttpClient, HttpErrorResponse } from "@angular/common/http";
import { CookieService } from 'ngx-cookie-service';
import { catchError, retry, throwError } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { ErrorComponent } from '../error/error.component';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { MatInputModule } from '@angular/material/input';


export interface GetUserResponse {
	login: string
	name: string
	surname: string
	interests: string
	bio: string
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, RouterModule],
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent {
	isDisabled: boolean = true;
	name: string = '';
	surname: string = '';
	bio: string = '';
	interests: string = '';
	login :string = '';
	
	constructor(
		private http: HttpClient, 
		private cookieService: CookieService,
		public dialog: MatDialog,
	){
		const id = localStorage.getItem("id");
		http.get<GetUserResponse>('http://localhost:8081/users/'+id, {withCredentials: true}).subscribe({
			next: response => {
				this.name = response.name;
				this.surname = response.surname;
				this.bio = response.bio;
				this.interests = response.interests;
				this.login = response.login;
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
	
	updateInfo(){
		const putData = new FormData();
		putData.append('name',  this.name);
		putData.append('surname',  this.surname);
		putData.append('bio',  this.bio);
		putData.append('interests',  this.interests);

		const options = {
			withCredentials: true
		};
		
		this.http.put('http://localhost:8081/users', putData, options).subscribe({
			next: response => {

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

	private handleError(error: HttpErrorResponse) {
		return throwError(() => new Error(error.message))
	}
	
	changeMode(){
		console.log(this.isDisabled);
		this.isDisabled = !this.isDisabled;
	}
}
