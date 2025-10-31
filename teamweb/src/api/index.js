import axios from "axios";
import qs from 'qs';

// const baseURL = '/v1/';
// const baseURL = 'http://172.20.10.9:8008/v1/';
const baseURL = 'http://localhost:9527/v1/stzb/';
// const baseURL = 'https://stzb.888718.xyz/stzb/';
// const baseURL = 'http://38.247.27.48:2888/stzb/';
// const baseURL = '/stzb/';
// const baseURL = 'http://192.168.31.251:8233/';
// const baseURL = 'http://localhost:8008/v1/';
// const baseURL = 'http://104.167.247.127:8008/v1/';
// const baseURL = 'http://192.168.1.101:8008/v1/';
// const skey = '352b0d7b04ae0d43';

// 创建 Axios 实例
const api = axios.create({
    baseURL:baseURL,
    timeout: 15000
});

// 添加请求拦截器
// api.interceptors.request.use(
//     (config) => {
//         if (localStorage.getItem('user')) {
//             let user = JSON.parse(localStorage.getItem('user'));
//             let token = user['token'];
//             config.headers['Bkmtoken'] = token;
//         }
//         return config;
//     },
//     (error) => {
//         return Promise.reject(error);
//     }
// );

export function reportList(data){
    return api.get('report/list?'+qs.stringify(data));
}

export function getPlayerTeam(data){
    return api.get('player/team/get?'+qs.stringify(data));
}

// export async function aitest(data) {
//     const response = await fetch(`${baseURL}ai/chat/completions`, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/x-www-form-urlencoded'
//         },
//         body: qs.stringify(data)
//     });

//     if (!response.ok) {
//         throw new Error(`HTTP error! status: ${response.status}`);
//     }

//     return response.body.getReader();
// }

// export async function novelChat(data) {
//     const response = await fetch(`${baseURL}ai/chat/completions`, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/x-www-form-urlencoded'
//         },
//         body: qs.stringify(data)
//     });

//     if (!response.ok) {
//         throw new Error(`HTTP error! status: ${response.status}`);
//     }

//     return response.body.getReader();
// }