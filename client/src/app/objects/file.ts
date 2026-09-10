export class UploadedFile{
    file?: File;
    url: String | null = null;
    name: string = "";

    constructor(){

    }

    processFile(file: File){
        this.name = file.name;
        this.file = file;
        this.url = URL.createObjectURL(file)
    }

    processPostURL(url: string){
        this.processURL("posts/"+url)
    }

    processMessageURL(url: string){
        this.processURL("messages/"+url)
    }

    private processURL(url: string){
        let splitted = url.split("/")
        this.name = splitted[splitted.length-1]
        this.url = `http://localhost:8082/upload/`+url;
    }

    getURL(): string{
        return this.url+""
    }

    alreadyUploaded(): boolean{
        if (this.file){
            return false
        }
        return true
    }
}