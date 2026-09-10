import { MatDialog } from "@angular/material/dialog";
import { ErrorComponent } from "../components/error/error.component";

export class ErrorWindow{

    constructor(err: any, dialog: MatDialog){
        const dialogRef = dialog.open(ErrorComponent, {
            data: {
                code: err.status,
                message: err.error,
                statusText: err.statusText,
            },
            height: '400px',
            width: '800px',
        })
    }
}