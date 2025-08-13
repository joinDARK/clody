import type { IBase } from "./base";

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
    _links: ILinksApi
    data: IBase
}