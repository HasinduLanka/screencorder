
const mime = 'video/webm; codecs="vp8"';

var Viewer;
var midcanvas;
var IsPlaying = false;

var buff = []; // minicache
var cfile = "new";

InitializeViewer();
Tick();

async function InitializeViewer() {
    console.log("InitializeViewer");

    Viewer = document.getElementById("Viewer");
    Viewer.addEventListener("error", async (event) => {
        IsPlaying = false;
        // console.log("Viewer Error");
        let error = event;

        // Chrome v60
        if (event.path && event.path[0]) {
            error = event.path[0].error;
        }

        // Firefox v55
        if (event.originalTarget) {
            error = error.originalTarget.error;
        }

        // Here comes the error message
        console.log(`Viewer error: ${error.message}`);

        setTimeout(() => {
            if (!IsPlaying) {
                ControlsHidden(false);
            }
        }, 3000);


        await NextBlob();
    });
    Viewer.addEventListener("ended", async event => {

        IsPlaying = false;
        console.log("Viewer end");

        setTimeout(() => {
            if (!IsPlaying) {
                ControlsHidden(false);
            }
        }, 3000);

        await NextBlob();
    });

    midcanvas = document.getElementById("MidCanvas");
    var ctx = midcanvas.getContext("2d");

    Viewer.addEventListener(
        "play",
        function () {
            ControlsHidden(true);

            var $this = this; //cache   

            (function SetCanvasHeight() {
                if (!$this.paused && !$this.ended) {
                    if ($this.videoWidth !== 0) {
                        midcanvas.width = window.innerWidth;
                        midcanvas.height = (midcanvas.width * $this.videoHeight) / $this.videoWidth;
                    }
                }
                setTimeout(SetCanvasHeight, 3000);
            })();

            (function loop() {
                if (IsPlaying && !$this.paused && !$this.ended) {
                    if ($this.videoWidth !== 0) {
                        ctx.drawImage($this, 0, 0, midcanvas.width, midcanvas.height);
                    }
                }
                setTimeout(loop, 40);
            })();

        },
        0
    );

    ControlsHidden(false);


    NextBlob();
    // Viewer.oncanplay = e => Viewer.play();
}



async function NextBlob() {
    IsPlaying = false;

    console.log("SRC");
    Viewer.src = await GetVideoSrc();
    console.log("SRC2");
    Viewer.load();
    IsPlaying = true;

}

async function GetVideoSrc() {
    while (buff.length == 0) {
        await sleep(100)
    }

    var c = buff.shift();
    return c;
}


async function AddSeg(blb) {

    // if (IsPlaying) {
    //     Viewer.play();
    //     IsPlaying = true;
    // }
    await RecieveBlob(blb);
}

async function RecieveBlob(blb) {
    if (blb.size != 0) {
        buff.push(URL.createObjectURL(blb));
        if (buff.length >= 3) {
            buff.shift();
        }
        console.log("Recieved blob " + blb.size);
    }

}

async function Tick() {
    setTimeout(Tick, 500);

    var resp = await fetch(`/mapi/reqview/` + cfile, { method: "GET" })
    var s = ("cpath : ", resp.headers.get("cpath"))
    if (s == "same") {
        console.info("Viewer : same")
    } else if (s == "wait") {
        console.info("Viewer : wait")
    } else {
        console.log("Viewer : file : " + s)
        cfile = s;
        AddSeg(await resp.blob())
    }
}


async function GetURL(path) {
    return await fetch(`/api/` + path, { method: "GET" })
        .then(response => response.blob())
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}


function ControlsHidden(hidden) {
    Viewer.hidden = hidden;
    midcanvas.hidden = !hidden;
    midcanvas.width = "98%"
}