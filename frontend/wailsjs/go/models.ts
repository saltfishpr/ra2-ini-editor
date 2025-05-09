export namespace main {
	
	export class Property {
	    ukey: string;
	    key: string;
	    value: string;
	    comment: string;
	    desc?: string;
	
	    static createFrom(source: any = {}) {
	        return new Property(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ukey = source["ukey"];
	        this.key = source["key"];
	        this.value = source["value"];
	        this.comment = source["comment"];
	        this.desc = source["desc"];
	    }
	}
	export class Unit {
	    type: string;
	    id: number;
	    name: string;
	    ui_name: string;
	    properties: Property[];
	
	    static createFrom(source: any = {}) {
	        return new Unit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.id = source["id"];
	        this.name = source["name"];
	        this.ui_name = source["ui_name"];
	        this.properties = this.convertValues(source["properties"], Property);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

