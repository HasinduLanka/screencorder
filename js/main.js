'use strict';

// Polyfill in Firefox.
// See https://blog.mozilla.org/webrtc/getdisplaymedia-now-available-in-adapter-js/
if (adapter.browserDetails.browser == 'firefox') {
  adapter.browserShim.shimGetDisplayMedia(window, 'screen');
}

const video = document.querySelector('video');
var mediaRecorder;

var BestMimeType;
var BestVideoFormat;
var QualityOptions ;
var GUMConstraints ;

var RecordEnabled = false;
var chunksPerStack = 30; // 30 || 5
var chunkLength = 1500; // 1500 || 10000
var chunkUploadURI = "recchunk/";
var startURI = "start/";

var filename = "Record-" + Date.now().toString();
var fileindex = 1;
var FinalizedFileindex = 0;
var FinalFile = "";
var IsRecording = false;
var IsEndSuccess = false;


var recordedChunks = [];

var CreateMediaRecorder ;
var MediaRecStream;


var Host_Encoding = "c"; // c, r, fh, fb, ml, fl, sh, su

Init_main();

function Init_main(){


  const supportedMimeTypes = getSupportedMimeTypes();
  BestMimeType = supportedMimeTypes[0].mime;
  BestVideoFormat = supportedMimeTypes[0].vtype;

  console.log('All supported mime types ordered by priority : ', supportedMimeTypes);

  console.log('Best supported mime type : ', BestMimeType);
  console.log('Best video format : ', BestVideoFormat);

  document.getElementById("LblEnc").innerHTML = "Your browser will encode in " + BestMimeType + " format";

  SetQlt("0");

  GetURL("handshake/" + filename);

  const DivSetQlt = document.getElementById('DivSetQlt');
  const DivVdo = document.getElementById('DivVdo');
  
  if ((navigator.mediaDevices && 'getDisplayMedia' in navigator.mediaDevices)) {
    DivSetQlt.hidden = false;
    DivVdo.hidden = true;
  } else {
    errorMsg('getDisplayMedia is not supported');
  }

}


function StartRec() {
  RecordEnabled = true;

  chunksPerStack = 5;
  chunkLength = 5000;

  navigator.mediaDevices.getDisplayMedia(GUMConstraints)
    .then(handleSuccess, handleError);
}


function handleSuccess(stream) {

  DivSetQlt.hidden = true;
  DivVdo.hidden = false;

  // demonstrates how to detect that the user has stopped
  // sharing the screen via the browser UI.
  stream.getVideoTracks()[0].addEventListener('ended', () => {
    errorMsg('Completed');
    IsRecording = false;
    IsEndSuccess = false;

    DivVdo.hidden = true;

    setTimeout(event => {
      if (!IsEndSuccess) {
        FinalizeRecord();
      }

      DivSetQlt.hidden = false;

    }, 3000);


  });


  GetURL(startURI + filename)  
  IsRecording = true;

  MediaRecStream = stream;
  video.srcObject = MediaRecStream;


  CreateMediaRecorder = () => {
    var newmediaRecorder = new MediaRecorder(MediaRecStream, QualityOptions);

    newmediaRecorder.ondataavailable = handleDataAvailable;
    newmediaRecorder.start();

    mediaRecorder = newmediaRecorder;
  };

  CreateMediaRecorder();  

  RestartMediaRecorder()

  console.log("Started Recording");

}

function handleError(error) {
  errorMsg(`getDisplayMedia error: ${error.name}`, error);
}


function errorMsg(msg, error) {
  const errorElement = document.querySelector('#errorMsg');
  errorElement.innerHTML += `<p>${msg}</p>`;
  if (typeof error !== 'undefined') {
    console.error(error);
  }
}




function RestartMediaRecorder() {
  console.log("Chunk done");
  if (IsRecording) {
    
    var oldMediaRec = mediaRecorder;
    CreateMediaRecorder()
    oldMediaRec.stop();

    setTimeout(event => {
      RestartMediaRecorder()
    }, chunkLength);
  }

}


function handleDataAvailable(event) {
  console.log("data-available");
  if (event.data.size > 0) {

    recordedChunks.push(event.data);
    console.log(recordedChunks);
    uploadRecorded();

    if (fileindex >= chunksPerStack || !IsRecording) {

      FinalizeRecord();      
    }
  } else {
    // ...
  }
  
}

function uploadRecorded() {
  var newrecordedChunks = recordedChunks;
  recordedChunks = [];

  var blob = new Blob(newrecordedChunks, {
    type: "video/" + BestVideoFormat
  });
  uploadBlob(chunkUploadURI +Host_Encoding + "/"+ "Ch-" + filename + "-" + pad(fileindex)+ "/" + BestVideoFormat, blob)
  fileindex += 1;

}


function FinalizeRecord() {


  var fflist = "";

  if (FinalFile.length != 0) {
    fflist = FinalFile + "\n";
  }
  
  FinalFile = filename;

  for (let i = FinalizedFileindex + 1; i < fileindex; i++) {
    fflist = fflist + "Ch-" + filename + "-" + pad(i) + "\n";
  }

  FinalizedFileindex = fileindex;

  var cfilename = filename;
  filename = "Record-" + Date.now().toString();
  var nextfilename;

  if (IsRecording) {
    nextfilename = filename;
    // GetURL(startURI + filename);
  }
  else
  {
    IsEndSuccess = true;
    nextfilename = "end";
  }

  uploadText("final/" + cfilename + "/" + nextfilename, fflist);

  fileindex = 1;
  FinalizedFileindex = 0;

  if (!IsRecording) {
    EndRecord();
  }


}


function EndRecord() {

  uploadText("end/" + filename, "");

  filename = "Record-" + Date.now().toString();
  fileindex = 1;
  FinalizedFileindex = 0;
  FinalFile = "";
  // finalFiles = [];
}


async function uploadText(path, txt) {
  var blob = new Blob([txt], { type: 'text/plain' });
  uploadBlob(path, blob);
}

//http://localhost:49542
async function uploadBlob(path, blob) {
  fetch(`/api/` + path, { method: "POST", body: blob, mode: 'no-cors' })
    .then(response => response.text().then(value => { console.log("API POST : " + path, value);}))
}

function pad(num) {
  var s = "000000000" + num;
  return s.substr(s.length - 10);
}

//http://localhost:49542
async function GetURL(path) {
  fetch(`/api/` + path, { method: "GET" })
    .then(response => response.text().then(value => console.log("API GET : " + path, value)))
}


function SelectQltyChanged() {
  var x = document.getElementById("SlctQlty").value;
  SetQlt(x);
  console.log("Quality selected : " + x);

}



function SetQlt(q) {

  

  var TQ = "";

  switch (q) {
    case "0":
      QualityOptions = {
        mimeType: BestMimeType
      }
      GUMConstraints = { video: true, audio: true };
      TQ = "Full quality and window framerate";
      break;

    case "1":
      QualityOptions = {
        videoBitsPerSecond: 8000000,
        audioBitsPerSecond: 128000,
        mimeType: BestMimeType
      }
      GUMConstraints = { video: { frameRate: { ideal: 60, max: 60 } }, audio: true  };
      TQ = "8 Mbps @ 60 FPS max";
      break;

    case "2":
      QualityOptions = {
        videoBitsPerSecond: 5000000,
        audioBitsPerSecond: 128000,
        mimeType: BestMimeType
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 30 } } , audio: true };
      TQ = "5 Mbps @ 30 FPS max";
      break;

    case "3":
      QualityOptions = {
        videoBitsPerSecond: 2500000,
        audioBitsPerSecond: 128000,
        mimeType: BestMimeType
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 30 } }, audio: true  };
      TQ = "2.5 Mbps @ 30 FPS max";
      break;

    case "4":
      QualityOptions = {
        videoBitsPerSecond: 1000000,
        audioBitsPerSecond: 128000,
        mimeType: BestMimeType
      }
      GUMConstraints = { video: { frameRate: { ideal: 24, max: 30 } } , audio: true };
      TQ = "1 Mbps @ 24 FPS max";
      break;

    case "5":
      QualityOptions = {
        videoBitsPerSecond: 600000,
        audioBitsPerSecond: 64000,
        mimeType: BestMimeType
      }
      GUMConstraints = { video: { frameRate: { ideal: 22, max: 24 } }, audio: true  };
      TQ = "600 kbps @ 22 FPS max";
      break;
    case "100":
      QualityOptions = {
        audioBitsPerSecond: 128000,
        mimeType: BestMimeType
      }

      GUMConstraints = { video: { frameRate: { ideal: 8, max: 10 } }, audio: true  };
      TQ = "Full quality @ 8 FPS max";
      break;
    default:
      break;
  }


  document.getElementById("LblQlty").innerHTML = TQ;

}



function SelectEncChanged() {
  var x = document.getElementById("SlctEnc").value;
  SetHostEnc(x);

}

function SetHostEnc(enc) {
  Host_Encoding = enc;
  console.log("Host Encoding : " + Host_Encoding);
  
}



function getSupportedMimeTypes() {
  const VIDEO_TYPES = [
    "webm", 
    "mp4",
  ];
  const VIDEO_CODECS = [
    "h264",
    "h.264",
    "h265",
    "h.265",
    "vp9",
    "vp9.0",
    "vp8",
    "vp8.0",
  ];

  const supportedTypes = [];
  VIDEO_TYPES.forEach((videoType) => {
    const type = `video/${videoType}`;
    VIDEO_CODECS.forEach((codec) => {
        const variations = [
        `${type};codecs=${codec}`,
        `${type};codecs:${codec}`,
        `${type};codecs=${codec.toUpperCase()}`,
        `${type};codecs:${codec.toUpperCase()}`,
      ]
      variations.forEach(variation => {
        if(MediaRecorder.isTypeSupported(variation)) 
            supportedTypes.push({mime : variation , vtype : videoType});
      })
    });
  });

  VIDEO_TYPES.forEach((videoType) => {
    const type = `video/${videoType}`;
    if(MediaRecorder.isTypeSupported(type)) 
            supportedTypes.push({mime : type , vtype : videoType});
      
  });
  
  return supportedTypes;
}
