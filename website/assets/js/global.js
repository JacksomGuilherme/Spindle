const tabState = {
    playlists: { page: 0 },
    albums: { page: 0 },
    artists: {
        cursorStack: [],
        currentPage: 0
    }
}

function openTab(tabId, e) {
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active')
    })

    e.classList.add('active')

    const container = document.getElementById("tab-content")

    let url = `/tab/content?tab=${tabId}`

    if (tabId === "artists") {
        url += `&after=`
        tabState.artists.cursorStack = []
        tabState.artists.currentPage = 0
    } else {
        tabState[tabId].page = 0
        url += `&page=0`
    }

    fetch(url)
        .then(res => res.text())
        .then(data => {
            container.innerHTML = data
            if (e.id === "artists") {
                updateArtistPage()
            }
        })
}

function changePage(value, nextCursor) {
    const activeTabId = document.querySelector(".tab-btn.active").id
    const container = document.getElementById("tab-content")

    let url = `/tab/content?tab=${activeTabId}`

    if (activeTabId !== "artists") {
        tabState[activeTabId].page = value
        url += `&page=${value}`
    }
    else {
        const state = tabState.artists

        if (value === "next") {
            state.cursorStack.push(nextCursor)
            state.currentPage++
            url += `&after=${nextCursor}`
        }

        if (value === "prev") {
            state.cursorStack.pop()
            state.currentPage--

            const prevCursor = state.cursorStack[state.cursorStack.length - 1] || ""
            url += `&after=${prevCursor}`
        }
    }

    fetch(url)
        .then(res => res.text())
        .then(data => {
            container.innerHTML = data
            if (activeTabId === "artists") {
                updateArtistPage()
            }
        })
}

function updateArtistPage() {
    const el = document.getElementById("artist-page")
    if (el) {
        el.innerText = tabState.artists.currentPage + 1
    }

    const prevBtn = document.querySelector(".pagination button:first-child")

    if (tabState.artists.currentPage === 0) {
        prevBtn.style.visibility = "hidden"
    } else {
        prevBtn.style.visibility = "visible"
    }
}

function play(contextURI) {
    let deviceId = document.getElementById('device-id').value

    let reqBody = { context_uri: contextURI }
    isPlaying = true

    togglePlaybackButtonState()

    fetch(`/playback/play?device_id=${deviceId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(reqBody)
    }).then(function (res) {
        setTimeout(refreshPlaybackState, 1000)
    })
}

isPlaying = false

function togglePlay() {
    let deviceId = document.getElementById('device-id').value
    if (isPlaying) {
        fetch(`/playback/pause?device_id=${deviceId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        })
            .then(function (res) {
                setTimeout(refreshPlaybackState, 1000)
            })
    } else {
        let reqBody = parsePlaybackState()
        fetch(`/playback/play?device_id=${deviceId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(reqBody)
        }).then(function (res) {
            setTimeout(refreshPlaybackState, 1000)
        })
    }
}

function next() {
    let deviceId = document.getElementById('device-id').value
    fetch(`/playback/next?device_id=${deviceId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(function (res) {
            setTimeout(refreshPlaybackState, 1000)
        })
}

function previous() {
    let deviceId = document.getElementById('device-id').value
    fetch(`/playback/previous?device_id=${deviceId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(function (res) {
            setTimeout(refreshPlaybackState, 1000)
        })
}

function togglePlaybackButtonState(isPlaying) {
    const play = document.getElementById("play");
    const pause = document.getElementById("pause");

    if (isPlaying) {
        play.style.display = "none";
        pause.style.display = "block";
    } else {
        play.style.display = "block";
        pause.style.display = "none";
    }
}

function disconnect() {
    window.location = "/logout"
}

let currentUri = ""
let currentPlaying = null

function refreshPlaybackState(attempt = 0) {
    const MAX_ATTEMPTS = 5

    fetch('/playback/state')
        .then(res => res.json())
        .then(data => {
            if (!data || !data.item) return

            const newUri = data.item.uri
            const newPlaying = data.is_playing

            const changed =
                currentUri !== newUri ||
                currentPlaying !== newPlaying

            updateUI(data)

            currentUri = newUri
            currentPlaying = newPlaying

            if (!changed && attempt < MAX_ATTEMPTS) {
                setTimeout(() => {
                    refreshPlaybackState(attempt + 1)
                }, 800)
            }
        })
}

window.onload = function () {
    refreshPlaybackState()
}

function parsePlaybackState() {
    let current_uri = document.getElementById("pb_state_current_uri").value
    let track_uri = document.getElementById("pb_state_track_uri").value
    let position_ms = document.getElementById("pb_state_position_ms").value

    let reqBody = {
        "context_uri": current_uri,
        "position_ms": position_ms
    }

    if (!current_uri.includes("artist")) {
        reqBody.offset = { "uri": track_uri }
    }
    return reqBody
}

function updateUI(data) {
    if (!data || !data.item) {
        document.getElementById("playback-song").innerText = "No song playing"
        document.getElementById("playback-artist").innerText = "..."
        document.getElementById("pb_state_current_uri").value = ""
        document.getElementById("pb_state_track_uri").value = ""
        document.getElementById("pb_state_position_ms").value = ""
        return
    }

    if (data.item.artists && data.item.name) {
        let artistNames = data.item.artists.map(artist => artist.name).join(", ")
        console.log(data.item.artists)
        document.getElementById("playback-song").innerText = data.item.name
        document.getElementById("playback-artist").innerText = artistNames
    }
    document.getElementById("pb_state_current_uri").value = data.context.uri
    document.getElementById("pb_state_track_uri").value = data.item.uri
    document.getElementById("pb_state_position_ms").value = data.progress_ms

    isPlaying = data.is_playing
    togglePlaybackButtonState(data.is_playing)
}