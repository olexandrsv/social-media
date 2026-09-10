import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router, RouterModule } from '@angular/router';
import { SharedService } from '../../services/shared.service';
import { WebsocketService } from '../../services/websocket.service';
import { CookieService } from 'ngx-cookie-service';
import { MatDialog } from '@angular/material/dialog';
import { ErrorComponent } from '../error/error.component';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';


export interface Response {
    token: string
	id: number
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, RouterModule],
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
	hide = true;
	login: string = "";
	name: string = "";
	surname: string = "";
	password: string = "";
	confirmedPassword: string = "";
	
	constructor(
		private http: HttpClient, 
		private router: Router, 
		private sharedService: SharedService,
		private websocket: WebsocketService,
		private cookieService: CookieService, 
		public dialog: MatDialog,
	){}

	onSubmit(){
		if (this.password != this.confirmedPassword){
			const dialogRef = this.dialog.open(ErrorComponent, {
				data: {
					code: "",
					message: "Login and Confirmed Login should be equal",
					statusText: "Wrong Confirmed Login",
				},
				height: '400px',
				width: '800px',
			})
		}
		const postData = new FormData();
		const login = this.login;
		postData.append('login',  login);
		postData.append('name', this.name);
		postData.append('surname',  this.surname);
		postData.append('password',  this.password);
		
		this.http.post<Response>('http://localhost:8081/users', postData).subscribe({
			next: response => {
				console.log(response);
				this.saveData(response.token, login, response.id);
				
				this.websocket.start();
				this.sharedService.start();
				this.sharedService.bol = true;
				this.router.navigate(['/']);
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
		this.clearData();
	}

	saveData(token: string, login: string, id: number){
		const expiryDate = new Date();
		expiryDate.setHours(expiryDate.getHours() + 1);
		this.cookieService.set('token', token, expiryDate, '/', 'localhost', true, 'None');

		localStorage.setItem("login", login);
		localStorage.setItem("id", id.toString());
	}

	clearData(){
		this.login = "";
		this.name = "";
		this.surname = "";
		this.password = "";
		this.confirmedPassword = "";
	}
}
