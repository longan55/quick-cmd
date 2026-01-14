export namespace main {
	
	export class Collection {
	    id: number;
	    name: string;
	    description?: string;
	    searchCount?: number;
	    os?: string[];
	    commandIds?: number[];
	    createdAt?: string;
	    updatedAt?: string;
	    deletedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Collection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.searchCount = source["searchCount"];
	        this.os = source["os"];
	        this.commandIds = source["commandIds"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.deletedAt = source["deletedAt"];
	    }
	}
	export class Command {
	    id: number;
	    name: string;
	    content?: string;
	    description?: string;
	    copyCount?: number;
	    searchCount?: number;
	    os?: string[];
	    tagIDs?: number[];
	    collectionIDs?: number[];
	    createdAt?: string;
	    updatedAt?: string;
	    deletedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Command(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.content = source["content"];
	        this.description = source["description"];
	        this.copyCount = source["copyCount"];
	        this.searchCount = source["searchCount"];
	        this.os = source["os"];
	        this.tagIDs = source["tagIDs"];
	        this.collectionIDs = source["collectionIDs"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.deletedAt = source["deletedAt"];
	    }
	}
	export class SortOption {
	    name?: string;
	    create_time?: string;
	    copy_counts?: string;
	
	    static createFrom(source: any = {}) {
	        return new SortOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.create_time = source["create_time"];
	        this.copy_counts = source["copy_counts"];
	    }
	}
	export class Option {
	    name: string;
	    os: string[];
	    type: string;
	    id: number;
	    sort: SortOption;
	
	    static createFrom(source: any = {}) {
	        return new Option(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.os = source["os"];
	        this.type = source["type"];
	        this.id = source["id"];
	        this.sort = this.convertValues(source["sort"], SortOption);
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
	export class Response {
	    code: number;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class SortIndex {
	    creatTimeAsc: boolean;
	    idAsc: boolean;
	    nameAsc: boolean;
	    sortValueAsc: boolean;
	    copyCountAsc: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SortIndex(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.creatTimeAsc = source["creatTimeAsc"];
	        this.idAsc = source["idAsc"];
	        this.nameAsc = source["nameAsc"];
	        this.sortValueAsc = source["sortValueAsc"];
	        this.copyCountAsc = source["copyCountAsc"];
	    }
	}
	
	export class Status {
	    os: string[];
	    sortIndex: SortIndex;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.sortIndex = this.convertValues(source["sortIndex"], SortIndex);
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
	export class Tag {
	    id: number;
	    name: string;
	    description?: string;
	    searchCount?: number;
	    os?: string[];
	    commandIds?: number[];
	    createdAt?: string;
	    updatedAt?: string;
	    deletedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Tag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.searchCount = source["searchCount"];
	        this.os = source["os"];
	        this.commandIds = source["commandIds"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.deletedAt = source["deletedAt"];
	    }
	}

}

