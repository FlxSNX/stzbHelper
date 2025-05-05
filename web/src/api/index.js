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

export function ApiGetTaskList(){
    return api.get("getTaskList");
}

export function ApiDelTask(id){
    return api.get(`deleteTask/${id}`);
}

export function ApiEnableGetReport(data){
    return api.post("enable/getReport",qs.stringify(data,{ arrayFormat: 'repeat' }));
}

export function ApiGetReportNumByTaskId(id){
    return api.get(`getReportNumByTaskId/${id}`);
}

export function ApiStatisticsReport(id){
    return api.get(`statisticsReport/${id}`);
}

export function ApiGetTask(id){
    return api.get(`getTask/${id}`);
}