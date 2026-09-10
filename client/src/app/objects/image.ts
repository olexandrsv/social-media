export class Image{
    name: string = "";
    img?: File;
    url: String | ArrayBuffer | null = null;

    constructor(){
       
    }

    processFile(img: File){
        this.name = img.name;
        this.img = img;

        const reader = new FileReader();
        reader.onload = () => {
            this.url = reader.result;
        }
        reader.readAsDataURL(img);
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
        if (this.img){
            return false
        }
        return true
    }
}