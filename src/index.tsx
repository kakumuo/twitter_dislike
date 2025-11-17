import React from "react";
import { createRoot } from "react-dom/client";
import {DislikeButton} from "./components/ComponentMain";


const SELECT_TIMELINE_PARENT = 'css-175oi2r'
const timelineParent = document.querySelector<HTMLDivElement>(`div[class*="${SELECT_TIMELINE_PARENT}"`)
const timelineButtonMap:WeakMap<HTMLDivElement, HTMLDivElement> = new WeakMap()

const observer = new MutationObserver((mutationsList) => {
  mutationsList.forEach((mutation) => {
    mutation.addedNodes.forEach((node) => {
        const element:HTMLDivElement = node as HTMLDivElement
        let targetLikeButton = element.querySelector<HTMLButtonElement>('button[data-testid="like"]')
        targetLikeButton = targetLikeButton ? targetLikeButton : element.querySelector<HTMLButtonElement>('button[data-testid="unlike"]')
            
        if(targetLikeButton && element.classList.contains("css-175oi2r") && element.getAttribute("data-testid") == "cellInnerDiv") {
            const tweetInfo = getTweetInfo(element)
            console.debug(element)
            addDislikeButton(targetLikeButton, tweetInfo)
        }
    });

    // TODO: find way to safely destroy object to prevent polling when out of screen; below code does not trigger
    // mutation.removedNodes.forEach((node) => {
    //     const element:HTMLDivElement = node as HTMLDivElement
    //     if(element.classList.contains("css-175oi2r") && element.classList.length == 1) {
    //         const targetButtonElement = timelineButtonMap.get(element)
    //         if(targetButtonElement) {
    //             targetButtonElement.remove()                
    //         }
	// 	}
    // })
  });
});

if(timelineParent) {
    observer.observe(timelineParent, { attributes: true, childList: true, subtree: true })
}

function addDislikeButton(likeButton:HTMLButtonElement, tweetInfo:TweetInfo) {
    const elementId = `__${tweetInfo.tweetId}`

    const container = document.createElement('div'); 
    container.id = elementId
    container.className = "css-175oi2r r-18u37iz r-1h0z5md r-13awgt0"

    ;(likeButton.parentElement as HTMLDivElement).insertAdjacentElement('afterend', container)

    const rootContainer = document.querySelector(`#${elementId}`)
    if(!rootContainer) {
        throw new Error("Cannot find elment container w/ id: " + elementId)
    }

    const root = createRoot(rootContainer)
    root.render(<DislikeButton tweetInfo={tweetInfo} />)

    return container
}

function getTweetInfo(parentTweetDiv:HTMLDivElement) {
    const info:TweetInfo = {
        ownerId: "", 
        profileId: "", 
        tweetId: ""
    }

    const profileId = getProfileId()

    const anchorElement = parentTweetDiv.querySelector<HTMLAnchorElement>("a[href*='/status/']")
    if (anchorElement) {
        const hrefElements = anchorElement.href.split("/")
        info.ownerId = hrefElements[3]
        info.tweetId = hrefElements[5]
        info.profileId = profileId
    } 

    return info
}

function getProfileId(){
	const profileAnchor = document.querySelector<HTMLAnchorElement>("a[aria-label='Profile'][role='link']") 
    if(profileAnchor) {
        const arr = profileAnchor.href.split('/')
        return arr[arr.length - 1]
    }
	
    return ""
}

