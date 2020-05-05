import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { delay, map } from 'rxjs/operators';
import { SatCertificate, WS, WSAuth, TimbradoResponse } from 'src/model/models';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class AppService {
  private header = new Headers({
    'Content-Type': 'application/json'
  });
  constructor(private http: HttpClient) { }

  // getCert obtiene el certificado que esta registrado en go.
  getCert(): Observable<SatCertificate> {
    return this.http.get<SatCertificate>(`/api/cert`);
  }

  // uploadCert carga el certificado para sellar los archivos.
  uploadCert(f: File): Observable<SatCertificate> {
    const formData: FormData = new FormData();
    formData.append('file', f, f.name);
    return this.http.post<SatCertificate>(`/api/cert`, formData);
  }

  // uploadKey carga la llave del certificado.
  uploadKey(f: File, p: string): Observable<void> {
    const formData: FormData = new FormData();
    formData.append('file', f, f.name);
    formData.append('pass', p);
    return this.http.post<void>(`/api/cert/key`, formData);
  }

  // catalogoWS regresa el catalogo de servicios web disponilbes para timbrar.
  catalogoWS(): Observable<WS[]> {
    return this.http.get<WS[]>(`/api/webservices`);
  }

  getInfoWS(): Observable<WSAuth> {
    return this.http.get<WSAuth>(`/api/webservices/ws`);
  }

  setInfoWS(ws: string, usuario: string, password: string): Observable<void> {
    return this.http.post<void>(`/api/webservices/ws`, {
      'ws': { 'id': ws },
      'usuario': usuario,
      'password': password,
    });
  }

  timbrar(f: File): Observable<TimbradoResponse> {
    const formData: FormData = new FormData();
    formData.append('file', f, f.name);
    return this.http.post<TimbradoResponse>(`/api/timbrar`, formData);
  }
}
