import axios from "axios";
import type { PageLoad } from './$types';
import type { AxiosResponse } from "axios";
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({params}) => {
    let res: AxiosResponse;
    try {
        res = await axios.get('http://localhost:8080/ping');
        console.debug(res.data);
    } catch (err) {
        throw error(404, "Response not found");
    }
    return res.data;
}