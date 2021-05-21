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
var RecordEnabled = true;
var QualityOptions = { mimeType: 'video/webm' }
var GUMConstraints = { video: true };

var filename = "Record-" + Date.now().toString();
var fileindex = 1;
var finalFiles = [];
var IsRecording = false;

GetURL("handshake/" + filename)

const DivSetQlt = document.getElementById('DivSetQlt');



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


const startButton = document.getElementById('startButton');
startButton.addEventListener('click', () => {
  RecordEnabled = false;
  navigator.mediaDevices.getDisplayMedia(GUMConstraints)
    .then(handleSuccess, handleError);
});
const startRecButton = document.getElementById('startRecButton');
startRecButton.addEventListener('click', () => {
  RecordEnabled = true;
  navigator.mediaDevices.getDisplayMedia(GUMConstraints)
    .then(handleSuccess, handleError);
});


if ((navigator.mediaDevices && 'getDisplayMedia' in navigator.mediaDevices)) {
  startButton.disabled = false;
  startRecButton.disabled = false;
  DivSetQlt.hidden = false;
} else {
  errorMsg('getDisplayMedia is not supported');
}


function handleSuccess(stream) {
  startButton.disabled = true;
  startRecButton.disabled = true;
  DivSetQlt.hidden = true;
  video.srcObject = stream;

  // demonstrates how to detect that the user has stopped
  // sharing the screen via the browser UI.
  stream.getVideoTracks()[0].addEventListener('ended', () => {
    errorMsg('The user has ended sharing the screen');
    IsRecording = false;

    startButton.disabled = false;
    startRecButton.disabled = false;
    DivSetQlt.hidden = false;

  });

  if (RecordEnabled) {
    console.log(stream);
    GetURL("start/" + filename)
    // var options = { mimeType: "video/webm; codecs=vp8" };
    IsRecording = true;
    mediaRecorder = new MediaRecorder(stream, QualityOptions);

    mediaRecorder.ondataavailable = handleDataAvailable;
    mediaRecorder.start();



    setTimeout(event => {
      RestartMediaRecorder()
    }, 10000);
  }
}


function RestartMediaRecorder() {
  console.log("Chunk done");
  if (IsRecording) {
    mediaRecorder.stop();
    mediaRecorder.start();
    setTimeout(event => {
      RestartMediaRecorder()
    }, 20000);
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
    if (fileindex >= 5 || !IsRecording) {
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
  uploadBlob(`chunk/` + filename + "-" + pad(fileindex), blob)
  fileindex += 1;

  recordedChunks = []
}
async function uploadText(path, txt) {
  var blob = new Blob([txt], { type: 'text/plain' });
  uploadBlob(path, blob);
}

//http://localhost:49542
async function uploadBlob(path, blob) {
  fetch(`/` + path, { method: "POST", body: blob, mode: 'no-cors' })
    .then(response => response.text().then(value => console.log(value)))
}

function pad(num) {
  var s = "000000000" + num;
  return s.substr(s.length - 10);
}

//http://localhost:49542
async function GetURL(path) {
  fetch(`/` + path, { method: "GET", mode: 'no-cors' })
    .then(response => response.text().then(value => console.log(value)))
}



function SetQlt(q) {

  switch (q) {
    case 0:
      QualityOptions = {
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: true };
      break;

    case 1:
      QualityOptions = {
        videoBitsPerSecond: 8000000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 60 } } };
      break;

    case 2:
      QualityOptions = {
        videoBitsPerSecond: 5000000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 30 } } };
      break;

    case 3:
      QualityOptions = {
        videoBitsPerSecond: 2500000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 30, max: 30 } } };
      break;

    case 4:
      QualityOptions = {
        videoBitsPerSecond: 1000000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 24, max: 30 } } };
      break;

    case 5:
      QualityOptions = {
        videoBitsPerSecond: 600000,
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 16, max: 24 } } };
      break;
    case 100:
      QualityOptions = {
        mimeType: 'video/webm'
      }
      GUMConstraints = { video: { frameRate: { ideal: 10, max: 20 } } };
      break;
    default:
      break;
  }
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