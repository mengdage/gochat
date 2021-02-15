import axios, { AxiosInstance } from 'axios'

let api: AxiosInstance | null = null

interface IniApiOption {
    serverAddr: string
    authorization?: string
}
function initApi({serverAddr, authorization}:IniApiOption): void {
    api = axios.create({
        baseURL: 'http://' + serverAddr,
        headers: {
            Authorization: authorization
        }
    })
}

function getApi(): AxiosInstance {
    if (!api) {
        throw new Error('api is undefined. Call initApi first')
    }
    return api
}

export { initApi, getApi }