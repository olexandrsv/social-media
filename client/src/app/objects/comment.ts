import { CommentModel } from "../components/posts/posts.component";
import { UploadedFile } from "./file";
import { Image } from "./image";
import { Word } from "./word";

export class Comment{
    id: string = "";
    parentID: string = "";
    userID: number = 0;
    userName: string = "";
    userSurname: string = "";
    text: string = "";
    images: Image[] = [];
    files: UploadedFile[] = [];
    deletedImages: Image[] = [];
    deletedFiles: UploadedFile[] = [];
    ownComment: boolean = false;
    parsedText: Word[][] = [];

    constructor(){

    }

    copy(comment: Comment){
        this.id = comment.id;
        this.userID = comment.userID;
        this.userName = comment.userName;
        this.userSurname = comment.userSurname;
        this.text = comment.text;
        this.images = comment.images;
        this.files = comment.files;
        this.parsedText = comment.parsedText;
    }

    onImgSelected(event: any){
        console.log(event)
        let files: File[] = event.target.files;
        for (let i=0; i<files.length; i++){
            const file = files[i];
            let image = new Image();
            image.processFile(file);
            this.addImage(image)
        }
    }
    
    onFileSelected(event: any){
        let files: File[] = event.target.files;
        for (let i=0; i<files.length; i++){
            let uploadedFile = new UploadedFile();
            uploadedFile.processFile(files[i]);
            this.addFile(uploadedFile);
        }
    }
    
    onImgRemoved(i: number){
        const img = this.images[i];
        if (img.alreadyUploaded()){
            this.deletedImages.push(img);
        }
        this.removeImage(i);
    }
    
    onFileRemoved(i: number){
        const file = this.files[i];
        if (file.alreadyUploaded()){
            this.deletedFiles.push(file);
        }
        this.removeFile(i);
    }

    private addImage(image: Image){
        this.images.push(image);
    }

    private addFile(file: UploadedFile){
        this.files.push(file);
    }

    private removeImage(i: number){
        this.images = this.images.filter((item, index) => {
            return index !== i;
        })
    }

    private removeFile(i: number){
        this.files = this.files.filter((item, index) => {
            return index !== i;
        })
    }

    getName(): string{
        if (this.ownComment){
            return "You"
        }
        if(this.userName !== "" || this.userSurname !== ""){
            return this.userName+" "+this.userSurname
        }
        return this.userID.toString()
    }

    generateCreateForm(): FormData{
        const form = new FormData();
        form.append("text", this.text)
        for(let i=0; i<this.images.length; i++){
            form.append("images[]", this.images[i].img!)
        }
        for(let i=0; i<this.files.length; i++){
            form.append("files[]", this.files[i].file!)
        }
        return form
    }

    generateUpdateForm(): FormData{
        const form = new FormData();
        form.append("text", this.text);
        for(let i=0; i<this.images.length; i++){
            const img = this.images[i];
            if (!img.alreadyUploaded()){
                form.append("images[]", img.img!)
            }
        }
        for(let i=0; i<this.files.length; i++){
            const file = this.files[i];
            if (!file.alreadyUploaded()){
                form.append("files[]", file.file!)
            }
        }
        for (let i=0; i<this.deletedImages.length; i++){
            form.append("deletedImages[]", this.deletedImages[i].name)
        }
        for (let i=0; i<this.deletedFiles.length; i++){
            form.append("deletedFiles[]", this.deletedFiles[i].name)
        }
        return form
    }

    parseHttpResponse(r: CommentModel){
        this.id = r.id;
        this.userID = r.user_id;
        this.text = r.text;
        this.userName = r.user_name;
        this.userSurname = r.user_surname;
        if (r.images){
            for (let i=0; i<r.images.length; i++){
                let image = new Image();
                image.processPostURL(r.images[i]);
                this.addImage(image);
            }
        }
        if (r.files){
            for (let i=0; i<r.files.length; i++){
                let file = new UploadedFile();
                file.processPostURL(r.files[i]);
                this.addFile(file);
            }
        }
        const id = localStorage.getItem("id")+""
        if (this.userID.toString() === id){
            this.ownComment = true;
        }
    }

    parseText(){
        this.parsedText = [];
        let list = this.text.split("\r\n")
        let code = false;
        for (let i = 0; i<list.length; i++){
            if (list[i].includes("-->")){
                code = !code;
                continue;
            }
            if (code){
                let map = this.parseLine(list[i]);
                this.parsedText.push(map);
            } else {
                let map: Word[] = [];
                map.push(new Word(list[i], "black"));
                this.parsedText.push(map);
            }
        }
        console.log(this.parsedText);
    }

    parseLine(line: string): Word[]{
        let wordList: Word[] = [];
        let list = line.split(" ");
        for (let i = 0; i<list.length; i++){
            let text = list[i];            
            let color = this.choosedColor(list[i]);
            let word = new Word(text, color);
            wordList.push(word);
        }
        return wordList
    }

    choosedColor(word: string): string {
        if (word === "static" || word === "void" || word === "public" || word == "int"){
            return "blue"
        } else if (word === "String"){
            return "green"
        }
        return "black"
    }
}