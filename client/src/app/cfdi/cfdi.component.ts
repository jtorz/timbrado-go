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
  responseTimbrado:TimbradoResponse;
  constructor(private appservice:AppService) { }

  ngOnInit(): void {
  }
  timbrar(){
    if(this.cfdiFile.nativeElement.files[0] == null){
      alert('Seleccione el archivo cfdi.');
      return;
    }
    this.appservice.timbrar(this.cfdiFile.nativeElement.files[0]).subscribe(
      (r)=>{this.responseTimbrado = r},
      (err) => {
        alert('Error: ' + JSON.stringify(err));
        console.log(err);
      }
    )
  }
}
