import type { ITable } from "./table";

export interface IBase {
    id: number,
    name: string,
    description?: string,
    created_at?: Date,
    tables?: ITable[],
}
