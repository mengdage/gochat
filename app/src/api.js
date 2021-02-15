import axios from 'axios'

let api

function initApi(serverAddr) {
    api = axios.create({
        baseURL: 'http://' + serverAddr
    })
}

function getApi() {
    if (!api) {
        throw new Error('api is undefined. Call initApi first')
    }
    return api
}

export { initApi, getApi }