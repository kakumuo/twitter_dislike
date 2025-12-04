import Browser from 'webextension-polyfill'

const BACKEND_HOST = "twitter-plus-fn-vsmsegopnq-ue.a.run.app"
const BACKEND_PORT = 8080
const BACKEND_PROTO = "http"
const BACKEND_PATH = "/tweet/dislike"
// const BACKEND_ENDPOINT = `${BACKEND_HOST}:${BACKEND_PORT}${BACKEND_PATH}`
const BACKEND_ENDPOINT = `${BACKEND_HOST}${BACKEND_PATH}`

Browser.runtime.onMessage.addListener(async (obj:any) => {
    const message:APIMessage = obj as APIMessage

    switch(message.action) {
        case 'get':
            return handleAPIGet(message)
        case 'upsert':
            return handleAPIPost(message)
    }
});

async function handleAPIGet(message:APIMessage){
    const query = Object.entries(message.data).map(([key, val]) => `${key}=${val}`).join("&")
    const path = `${BACKEND_PROTO}://${BACKEND_ENDPOINT}?${query}`
    
    // console.debug("Sending GET to:", path)

    const resp  = await fetch(path, {method: "get"})	
    const respJson  = await resp.json()
    // console.debug("API Response: ", respJson)   
    
    return respJson
}


async function handleAPIPost(message:APIMessage){
    const query = Object.entries(message.data).map(([key, val]) => `${key}=${val}`).join("&")
    const path = `${BACKEND_PROTO}://${BACKEND_ENDPOINT}?${query}`
    
    // console.debug("Sending POST to:", path)

    const resp  = await fetch(path, {method: "post"})	
    const respJson  = await resp.json()
    // console.debug("API Response: ", respJson)
    
    return respJson
}