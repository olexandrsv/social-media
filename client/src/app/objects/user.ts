export function newUser(id: number, login: string): User{
    let user = new User()
    user.id = id
    user.login = login
    return user
}

export class User{
    id: number = 0;
    login: string = "";
    name: string = "";
    surname: string = "";

    constructor(){
        
    }
}