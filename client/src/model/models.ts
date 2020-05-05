export class SatCertificate {
    serialNumber: string;
    emisor: Issuer;
    receptor: Issuer;
    constructor() {
        this.serialNumber = "";
        this.emisor = new Issuer();
        this.receptor = new Issuer;
    }
}

export class Issuer {
    o: string;
    ou: string;
    cn: string;
    constructor() {
        this.o = "";
        this.ou = "";
        this.cn = "";
    }
}

export class WS {
    id: string;
    name: string;
}

export class WSAuth {
    ws: WS;
    usuario: string;
    password: string;
}

export class TimbradoResponse {
    StatusCode: string;
    Message: string;
    CFDI: string;
}