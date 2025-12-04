import React from "react";
import Browser, { action } from "webextension-polyfill";

const CLASS_CONTAINER= "css-175oi2r r-18u37iz r-1h0z5md r-13awgt0"
const CLASS_BUTTON = "css-175oi2r r-1777fci r-bt1l66 r-bztko3 r-lrvibr r-1loqt21 r-1ny4l3l"
const CLASS_BUTTON_DIV = "css-146c3p1 r-bcqeeo r-1ttztb7 r-qvutc0 r-37j5jr r-a023e6 r-rjixqe r-16dba41 r-1awozwy r-6koalj r-1h0z5md r-o7ynqc r-clp7b1 r-3s2u2q"
const CLASS_ICON_DIV = "css-175oi2r r-xoduu5"
const CLASS_LABEL_DIV = "css-175oi2r r-xoduu5 r-1udh08x"
const CLASS_HALO_DIV_DEFAULT = "css-175oi2r r-xoduu5 r-1p0dtai r-1d2f490 r-u8s1d r-zchlnj r-ipm5af r-1niwhzg r-sdzlij r-xf4iuw r-o7ynqc r-6416eg r-1ny4l3l"
const CLASS_HALO_DIV_HOVER = "css-175oi2r r-xoduu5 r-1p0dtai r-1d2f490 r-u8s1d r-zchlnj r-ipm5af r-sdzlij r-xf4iuw r-o7ynqc r-6416eg r-1peqgm7 r-1ny4l3l"

// "transform: translate3d(0px, 0px, 0px); transition-property: transform; transition-duration: 0.3s;"
const STYLE_LABEL:React.CSSProperties = {
	transform: "translate3d(0px, 0px, 0px)",
	transitionProperty: "transofrm", 
	transitionDuration: "0.3s"
} 

const DISLIKE_POLL_RATE_MS = 2000

export function DislikeButton(props:{tweetInfo:TweetInfo}) {  
	const [count, setCount] = React.useState(0)
	const [dislike, setDislike] = React.useState(false)
	const [colorDefault] = React.useState("rgb(113, 118, 123)")
	const [colorHover] = React.useState("#c043faff")
	const [colorHoverHalo] = React.useState("#c043fa26")
	const [hover, setHover] = React.useState(false)
	

	const handleTweetDislikes = async (method: 'get' | 'upsert') => {
		const resp = await Browser.runtime.sendMessage<APIMessage, APIResponse>({
			action: method,
			data: props.tweetInfo
		})

		if("error" in resp) {
			console.error("Failed to like tweet: ", resp.error)
		} else if ("dislikeCount" in resp) {
			setCount(resp.dislikeCount)
			setDislike(resp.userDislike)
		}
	}

	React.useEffect(() => {
		// console.debug("Creating component for: ", props.tweetInfo)
		handleTweetDislikes('get')
	}, [])

	return <button data-testid="dislike" className={CLASS_BUTTON} onClick={() => handleTweetDislikes('upsert')} onMouseEnter={() => setHover(true)} onMouseLeave={() => setHover(false)}>
			<div className={CLASS_BUTTON_DIV} style={{color: hover || dislike ? colorHover : colorDefault}} dir="ltr" >
				{/* Icon */}
				<div className={CLASS_ICON_DIV}>
					<div className={hover ? CLASS_HALO_DIV_HOVER : CLASS_HALO_DIV_DEFAULT} style={{backgroundColor: hover ? colorHoverHalo : ""}} />
					<svg viewBox="0 0 24 24" transform="scale(1 -1)" aria-hidden="true" className="r-4qtqp9 r-yyyyoo r-dnmrzs r-bnwqim r-lrvibr r-m6rgpd r-1xvli5t r-1hdv0qi">
						<g>
						{dislike ? 
							<path d="M20.884 13.19c-1.351 2.48-4.001 5.12-8.379 7.67l-.503.3-.504-.3c-4.379-2.55-7.029-5.19-8.382-7.67-1.36-2.5-1.41-4.86-.514-6.67.887-1.79 2.647-2.91 4.601-3.01 1.651-.09 3.368.56 4.798 2.01 1.429-1.45 3.146-2.1 4.796-2.01 1.954.1 3.714 1.22 4.601 3.01.896 1.81.846 4.17-.514 6.67z"></path>
							:
							<path d="M16.697 5.5c-1.222-.06-2.679.51-3.89 2.16l-.805 1.09-.806-1.09C9.984 6.01 8.526 5.44 7.304 5.5c-1.243.07-2.349.78-2.91 1.91-.552 1.12-.633 2.78.479 4.82 1.074 1.97 3.257 4.27 7.129 6.61 3.87-2.34 6.052-4.64 7.126-6.61 1.111-2.04 1.03-3.7.477-4.82-.561-1.13-1.666-1.84-2.908-1.91zm4.187 7.69c-1.351 2.48-4.001 5.12-8.379 7.67l-.503.3-.504-.3c-4.379-2.55-7.029-5.19-8.382-7.67-1.36-2.5-1.41-4.86-.514-6.67.887-1.79 2.647-2.91 4.601-3.01 1.651-.09 3.368.56 4.798 2.01 1.429-1.45 3.146-2.1 4.796-2.01 1.954.1 3.714 1.22 4.601 3.01.896 1.81.846 4.17-.514 6.67z"></path>
						}
						</g>
					</svg>
				</div>

				{/* Label */}
				{count > 0 && 
					<div className={CLASS_LABEL_DIV}>
						<span data-testid="app-text-transition-container" style={STYLE_LABEL}>
							<span className="css-1jxf684 r-1ttztb7 r-qvutc0 r-poiln3 r-n6v787 r-1cwl3u0 r-1k6nrdp r-n7gxbd">
								<span className="css-1jxf684 r-bcqeeo r-1ttztb7 r-qvutc0 r-poiln3">
									{count}
								</span>
							</span>
						</span>
					</div>
				}
			</div>
		</button>
}


/*
    const CLASS_FOOTER= "css-175oi2r r-18u37iz r-1h0z5md r-13awgt0"
    const CLASS_BUTTON = "css-175oi2r r-1777fci r-bt1l66 r-bztko3 r-lrvibr r-1loqt21 r-1ny4l3l"
    const CLASS_BUTTON_DIV = "css-146c3p1 r-bcqeeo r-1ttztb7 r-qvutc0 r-37j5jr r-a023e6 r-rjixqe r-16dba41 r-1awozwy r-6koalj r-1h0z5md r-o7ynqc r-clp7b1 r-3s2u2q"
    const CLASS_ICON_DIV = "css-175oi2r r-xoduu5"
    const CLASS_LABEL_DIV = "css-175oi2r r-xoduu5 r-1udh08x"

    console.debug("parent tweet", parentTweetDiv)
    const dislikes = 0

    const footer = parentTweetDiv.querySelector("div[aria-label*='likes']")
    const container = document.createElement('div'); 
    const button = document.createElement('button'); 
    
    const button_div = document.createElement('div');
    button_div.dir = "ltr"
    button_div.style = "color: rgb(113, 118, 123);"

    const label_div = document.createElement('div'); 
    label_div.className = CLASS_LABEL_DIV
    label_div.innerHTML = `
        <span data-testid="app-text-transition-container" style="transform: translate3d(0px, 0px, 0px); transition-property: transform; transition-duration: 0.3s;">
            <span class="css-1jxf684 r-1ttztb7 r-qvutc0 r-poiln3 r-n6v787 r-1cwl3u0 r-1k6nrdp r-n7gxbd">
                <span class="css-1jxf684 r-bcqeeo r-1ttztb7 r-qvutc0 r-poiln3">${dislikes}</span>
            </span>
        </span>
    `

    const icon_div = document.createElement('div'); 
    icon_div.className = CLASS_ICON_DIV
    icon_div.innerHTML = `<svg viewBox="0 0 24 24" transform="scale(1 -1)" aria-hidden="true" class="r-4qtqp9 r-yyyyoo r-dnmrzs r-bnwqim r-lrvibr r-m6rgpd r-1xvli5t r-1hdv0qi"><g><path d="M16.697 5.5c-1.222-.06-2.679.51-3.89 2.16l-.805 1.09-.806-1.09C9.984 6.01 8.526 5.44 7.304 5.5c-1.243.07-2.349.78-2.91 1.91-.552 1.12-.633 2.78.479 4.82 1.074 1.97 3.257 4.27 7.129 6.61 3.87-2.34 6.052-4.64 7.126-6.61 1.111-2.04 1.03-3.7.477-4.82-.561-1.13-1.666-1.84-2.908-1.91zm4.187 7.69c-1.351 2.48-4.001 5.12-8.379 7.67l-.503.3-.504-.3c-4.379-2.55-7.029-5.19-8.382-7.67-1.36-2.5-1.41-4.86-.514-6.67.887-1.79 2.647-2.91 4.601-3.01 1.651-.09 3.368.56 4.798 2.01 1.429-1.45 3.146-2.1 4.796-2.01 1.954.1 3.714 1.22 4.601 3.01.896 1.81.846 4.17-.514 6.67z"></path></g></svg>`

*/