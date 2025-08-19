import axios from "axios";
import type { PageLoad } from "./$types";
import type { AxiosResponse } from "axios";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ params }) => {
  const id: string = params.id;
  let res: AxiosResponse;
  try {
    res = await axios.get(`http://localhost:8080/api/tables/${id}`);
    console.debug("[axios] data: ", res.data);
  } catch(e) {
    throw error(500, "Internal Server Error")
  }
  return res.data
};
