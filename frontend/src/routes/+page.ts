import axios from "axios";
import type { PageLoad } from './$types';
import type { AxiosResponse } from "axios";
import { error } from '@sveltejs/kit';
import type { IBaseApi } from '@share/interfaces/api'

export const load: PageLoad<IBaseApi> = async ({params}) => {
    let res: AxiosResponse;
    try {
        res = await axios.get('http://localhost:8080/api/bases/1', {
            params: {
                WithTables: 1
            }
        });
        console.debug("axios.get.data: ", res.data);
    } catch (err) {
        throw error(500, "Internal Sever Error");
    }
    return res.data;
}