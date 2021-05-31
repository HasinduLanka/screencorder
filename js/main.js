/*
 *  Copyright (c) 2018 The WebRTC project authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a BSD-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree.
 */
'use strict';

// Polyfill in Firefox.
// See https://blog.mozilla.org/webrtc/getdisplaymedia-now-available-in-adapter-js/
if (adapter.browserDetails.browser == 'firefox') {
  adapter.browserShim.shimGetDisplayMedia(window, 'screen');
}

const video = document.querySelector('video');
var mediaRecorder;
var QualityOptions = { mimeType: 'video/webm' }
var GUMConstraints = { video: true };

var MirrorEnabled = false;
var RecordEnabled = false;
var chunksPerStack = 30; // 30 || 5
var chunkLength = 1500; // 1500 || 10000
var chunkUploadURI = "chunk/";
var startURI = "start/";

var filename = "Record-" + Date.now().toString();
var fileindex = 1;
var finalFiles = [];
var IsRecording = false;

GetURL("handshake/" + filename)

const DivSetQlt = document.getElementById('DivSetQlt');
const DivVdo = document.getElementById('DivVdo');




function StartMirror() {
  RecordEnabled = false;
  MirrorEnabled = true;

  chunksPerStack = 30;
  chunkLength = 1500;
  chunkUploadURI = "mirrorchunk/";
  startURI = "mstart/";

  navigator.mediaDevices.getDisplayMedia(GUMConstraints)
    .then(handleSuccess, handleError);
}

function StartRec() {
  RecordEnabled = true;
  MirrorEnabled = false;

  chunksPerStack = 5;
  chunkLength = 10000;
  chunkUploadURI = "recchunk/";
  startURI = "start/";

  navigator.mediaDevices.getDisplayMedia(GUMConstraints)
    .then(handleSuccess, handleError);
}
function StartMirrorRec() {
  RecordEnabled = true;
  MirrorEnabled = true;

  chunksPerStack = 40;
  chunkLength = 1500;
  chunkUploadURI = "mirecchunk/";
  startURI = "start/";

  navigator.mediaDevices.getDisplayMedia(GUMConstraints)
    .then(handleSuccess, handleError);
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
if ((navigator.mediaDevices && 'getDisplayMedia' in navigator.mediaDevices)) {
  DivSetQlt.hidden = false;
  DivVdo.hidden = true;
} else {
  errorMsg('getDisplayMedia is not supported');
}


function handleSuccess(stream) {
  DivSetQlt.hidden = true;
  DivVdo.hidden = false;
  video.srcObject = stream;

  // demonstrates how to detect that the user has stopped
  // sharing the screen via the browser UI.
  stream.getVideoTracks()[0].addEventListener('ended', () => {
    errorMsg('Completed');
    IsRecording = false;

    DivSetQlt.hidden = false;
    DivVdo.hidden = true;

  });


  console.log(stream);
  GetURL(startURI + filename)
  // var options = { mimeType: "video/webm; codecs=vp8" };
  IsRecording = true;
  mediaRecorder = new MediaRecorder(stream, QualityOptions);

  mediaRecorder.ondataavailable = handleDataAvailable;
  mediaRecorder.start();

  RestartMediaRecorder()


}


function RestartMediaRecorder() {
  console.log("Chunk done");
  if (IsRecording) {
    mediaRecorder.stop();
    mediaRecorder.start();
    setTimeout(event => {
      RestartMediaRecorder()
    }, chunkLength);
  }

}

var recordedChunks = [];



async function downloadText(txt, filename) {
  var blob = new Blob([txt], { type: 'text/plain' });

  // this will create a link tag on the fly
  // <a href="..." download>
  var link = document.createElement('a');
  link.setAttribute('href', URL.createObjectURL(blob));
  link.setAttribute('download', filename);

  // NOTE: We need to add temporarily the link to the DOM so
  //       we can trigger a 'click' on it.
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

function handleDataAvailable(event) {
  console.log("data-available");
  if (event.data.size > 0) {
    recordedChunks.push(event.data);
    console.log(recordedChunks);
    // download();
    uploadRecorded();
    if (fileindex >= chunksPerStack || !IsRecording) {
      FinalizeRecord()
      if (IsRecording) {
        GetURL("start/" + filename)
      }
    }
  } else {
    // ...
  }

  if (!IsRecording) {
    EndRecord()
  }


}


function FinalizeRecord() {

  finalFiles.push(filename);

  var fflist = "";

  for (let i = 1; i < fileindex; i++) {
    fflist = fflist + "file '" + filename + "-" + pad(i) + ".webm'\n";
  }

  uploadText("final/" + filename, fflist)

  filename = "Record-" + Date.now().toString()
  fileindex = 1
}


function EndRecord() {

  var fflist = "";

  finalFiles.forEach(f => {
    fflist = fflist + "file '" + f + ".webm'\n";
  });

  uploadText("end/" + filename, fflist)

  filename = "Record-" + Date.now().toString()
  fileindex = 1
  finalFiles = [];
}

function download() {
  var blob = new Blob(recordedChunks, {
    type: "video/webm"
  });
  var url = URL.createObjectURL(blob);
  var a = document.createElement("a");
  document.body.appendChild(a);
  a.style = "display: none";
  a.href = url;
  a.download = filename + "-" + pad(fileindex) + ".webm";
  fileindex += 1;
  a.click();
  window.URL.revokeObjectURL(url);
  recordedChunks = []
}


function uploadRecorded() {
  var blob = new Blob(recordedChunks, {
    type: "video/webm"
  });
  uploadBlob(chunkUploadURI + filename + "-" + pad(fileindex), blob)
  fileindex += 1;

  recordedChunks = []
}
async function uploadText(path, txt) {
  var blob = new Blob([txt], { type: 'text/plain' });
  uploadBlob(path, blob);
}

//http://localhost:49542
async function uploadBlob(path, blob) {
  fetch(`/api/` + path, { method: "POST", body: blob, mode: 'no-cors' })
    .then(response => response.text().then(value => console.log(value)))
}

function pad(num) {
  var s = "000000000" + num;
  return s.substr(s.length - 10);
}

//http://localhost:49542
async function GetURL(path) {
  fetch(`/api/` + path, { method: "GET" })
    .then(response => response.text().then(value => console.log(value)))
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
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: true };
      TQ = "Full quality and window framerate";
      break;

    case "1":
      QualityOptions = {
        videoBitsPerSecond: 8000000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 60, max: 60 } } };
      TQ = "8 Mbps @ 60 FPS max";
      break;

    case "2":
      QualityOptions = {
        videoBitsPerSecond: 5000000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 30 } } };
      TQ = "5 Mbps @ 30 FPS max";
      break;

    case "3":
      QualityOptions = {
        videoBitsPerSecond: 2500000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 30 } } };
      TQ = "2.5 Mbps @ 30 FPS max";
      break;

    case "4":
      QualityOptions = {
        videoBitsPerSecond: 1000000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 24, max: 30 } } };
      TQ = "1 Mbps @ 24 FPS max";
      break;

    case "5":
      QualityOptions = {
        videoBitsPerSecond: 600000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 22, max: 24 } } };
      TQ = "600 kbps @ 22 FPS max";
      break;
    case "100":
      QualityOptions = {
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 8, max: 10 } } };
      TQ = "Full quality @ 8 FPS max";
      break;
    default:
      break;
  }


  document.getElementById("LblQlty").innerHTML = TQ;

}



function PlaybackMic() {
  navigator.getUserMedia = navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia;

  var aCtx;
  var analyser;
  var microphone;
  if (navigator.getUserMedia) {
    navigator.getUserMedia(
      { audio: true },
      function (stream) {

        aCtx = new (window.AudioContext || window.webkitAudioContext)();
        // aCtx = new AudioContext();

        window.AudioContext.
          microphone = aCtx.createMediaStreamSource(stream);
        var destination = aCtx.destination;
        microphone.connect(destination);


      },
      function () { console.log("Error 003.") }
    );
  }
}

// const downloadButton = document.getElementById('downloadButton');

// downloadButton.addEventListener('click', () => {

//   var fflist = "";
//   var fflistfilename = filename + ".fflist";

//   var cmd = "#!/bin/bash\n";
//   cmd = cmd + "ffmpeg -f concat -safe 0 -i " + fflistfilename + " -c copy " + filename + ".webm\nrm -f " + filename + "-*.webm\nrm -f " + filename + ".sh\n";


//   for (let i = 1; i < fileindex; i++) {
//     fflist = fflist + "file '" + filename + "-" + pad(i) + ".webm'\n";
//   }

//   downloadText(fflist, fflistfilename);
//   // uploadText("final/" + filename, fflist)
//   downloadText(cmd, filename + ".sh");

// });