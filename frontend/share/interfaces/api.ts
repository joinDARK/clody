import type { IBase } from "./base";
import type { ITable } from "./table";

export interface ILinksApi {
    create: {
        href: string;
        method: string;
    };
    delete: {
        href: string;
        method: string;
    };
    self: {
        href: string;
        params: string;
    };
    update: {
        href: string;
        method: string;
    };
}

export interface IBaseApi {
    _links: ILinksApi;
    data: IBase;
}

export interface ITableApi {
    _links: ILinksApi;
    data: ITable;
}
