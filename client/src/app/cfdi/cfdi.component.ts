import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { TimbradoResponse } from 'src/model/models';
import { AppService } from 'src/services/app.service';

@Component({
  selector: 'app-cfdi',
  templateUrl: './cfdi.component.html',
  styleUrls: ['./cfdi.component.scss']
})
export class CfdiComponent implements OnInit {

  @ViewChild("cfdiFile") cfdiFile: ElementRef;
  responseTimbrado: TimbradoResponse;
  constructor(private appservice: AppService) { }

  ngOnInit(): void {
  }
  timbrar() {
    if (this.cfdiFile.nativeElement.files[0] == null) {
      alert('Seleccione el archivo cfdi.');
      return;
    }
    this.appservice.timbrar(this.cfdiFile.nativeElement.files[0]).subscribe(
      (r) => { this.responseTimbrado = r },
      (err) => {
        this.responseTimbrado = new TimbradoResponse();
        if (err.status) {
          this.responseTimbrado.StatusCode = err.status
        }
        if (err.error) {
          if (err.error.StatusCode) {
            this.responseTimbrado.StatusCode = err.error.StatusCode;
          }
          if (err.error.Message) {
            this.responseTimbrado.Message = err.error.Message;
          }else{
            this.responseTimbrado.Message = JSON.stringify(err.error, null, "\n")
          }
        } else {
          alert('Error: ' + JSON.stringify(err));
        }
        console.log(err);
      }
    )
  }
}
