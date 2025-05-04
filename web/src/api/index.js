import axios from "axios";
import qs from 'qs';

// const baseURL = '/v1/';
const baseURL = 'http://localhost:9527/v1';

const api = axios.create({
    baseURL:baseURL,
    timeout: 5000
});

export function ApiGetTeamUser(){
    return api.get("getTeamUser");
}

export function ApiGetTeamGroup(){
    return api.get("getTeamGroup");
}

export function ApiCreateTask(data){
    return api.post("createTask",qs.stringify(data,{ arrayFormat: 'repeat' }));
}