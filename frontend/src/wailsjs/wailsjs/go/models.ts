export namespace models {
	
	export class MeasurementDTO {
	    id: number;
	    timestamp: string;
	    download_speed: number;
	    upload_speed: number;
	    latency: number;
	
	    static createFrom(source: any = {}) {
	        return new MeasurementDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = source["timestamp"];
	        this.download_speed = source["download_speed"];
	        this.upload_speed = source["upload_speed"];
	        this.latency = source["latency"];
	    }
	}

}

