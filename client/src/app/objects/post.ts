import { GetPostResponse } from "../components/posts/posts.component";
import { UploadedFile } from "./file";
import { Image } from "./image";
import { Word } from "./word";

export class Post {
    id: string = "";
    userID: number = 0;
    text: string = "";
    images: Image[] = [];
    files: UploadedFile[] = [];
    deletedImages: Image[] = [];
    deletedFiles: UploadedFile[] = [];
    parsedText: Word[][] = [];
    isRead = false;

    constructor(){

    }

    copy(post: Post){
        this.id = post.id;
        this.userID = post.userID;
        this.text = post.text;
        this.images = post.images;
        this.files = post.files;
        this.parsedText = post.parsedText;
    }

    onImgSelected(event: any){
        console.log(event)
        let files: File[] = event.target.files;
        for (let i=0; i<files.length; i++){
            console.log("here")
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

    parseHttpResponse(r: GetPostResponse){
        this.id = r.id;
        this.userID = r.user_id;
        this.text = r.text;
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
        this.parseText()
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
        let char = "";
        let list = line.split(" ");
        for (let i = 0; i<list.length; i++){
            let text = list[i];
            if (text.endsWith(";")){
                char = text.slice(-1);
                text = text.slice(0, text.length - 1);
                console.log("char: " + char);
                console.log("text: " + text);
            }
            let color = this.choosedColor(text);
            let word = new Word(text, color);
            wordList.push(word);
            wordList.push(new Word(char, "black"));
            char = "";
        }
        return wordList
    }

    choosedColor(word: string): string {
        if (word === "static" || word === "void" || word === "public" || word == "int" ||
            word === "this" || word === "super" || word === "extends" || word === "implements" || word === "new" ||
            word === "package" || word === "import" || word === "try" || word === "catch" || word === "finally" ||
            word === "throw" || word === "throws" || word === "synchronized" || word === "volatile" ||
            word === "struct" || word === "type" || word === "float" || word === "double" || word === "bool" ||
            word === "func"
        ){
            return "blue"
        } else if (word === "String"){
            return "green"
        } else if (word === "class" || word === "if" || word === "else" || word === "for" || word === "while" || word === "return"){
            return "red"
        } else if (word.startsWith("\"") || word.endsWith("\"")){
            return "green"
        }
        return "black"
    }
}