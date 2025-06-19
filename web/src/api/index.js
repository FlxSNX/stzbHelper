import axios from "axios";
import qs from 'qs';

// const baseURL = '/v1/';
const baseURL = 'http://localhost:9527/v1';

const api = axios.create({
    baseURL:baseURL,
    timeout: 1000 * 180
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

export function ApiGetGroupWu(){
    return api.get("getGroupWu");
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

export function ApiDelTaskReport(id){
    return api.get(`deleteTaskReport/${id}`);
}