export namespace common {
	
	export class YamlInfo {
	    dir: string;
	    env: {[key: string]: string};
	
	    static createFrom(source: any = {}) {
	        return new YamlInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir = source["dir"];
	        this.env = source["env"];
	    }
	}

}

export namespace main {
	
	export class Config {
	    index: number;
	    name: string;
	    command: string;
	    env: string;
	    dir: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.name = source["name"];
	        this.command = source["command"];
	        this.env = source["env"];
	        this.dir = source["dir"];
	    }
	}
	export class TypeConfig {
	    index: number;
	    type: string;
	    config: Config[];
	
	    static createFrom(source: any = {}) {
	        return new TypeConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.type = source["type"];
	        this.config = this.convertValues(source["config"], Config);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

