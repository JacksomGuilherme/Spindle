function openTab(tabId, e) {
    document.querySelectorAll('.tab-btn').forEach(e => {
        e.classList.remove('active')
    })

    e.classList.add('active')

    window.location.href = `/?tab=${tabId}`
}

function play(contextURI) {
    let deviceId = document.getElementById('device-id').value
    fetch(`/playback/play?device_id=${deviceId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ context_uri: contextURI })
    })
}

let isPlaying = false;

function togglePlay() {
  const play = document.getElementById("play");
  const pause = document.getElementById("pause");

  if (isPlaying) {
    play.style.display = "block";
    pause.style.display = "none";
  } else {
    play.style.display = "none";
    pause.style.display = "block";
  }

  isPlaying = !isPlaying;

//   fetch("/playback/toggle", { method: "POST" });
}

function disconnect(){
    window.location = "/logout"
}
