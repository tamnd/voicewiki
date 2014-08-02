window.AudioContext = window.AudioContext || window.webkitAudioContext;
var audioContext = new AudioContext();
var audioInput = null,
    inputPoint = null,
    realAudioInput = null,
    recorder = null;

function toggleRecording( e ) {
    if (e.classList.contains("recording")) {
        // stop recording
        // audioRecorder.stop();
        stopRecording()
        e.classList.remove("recording");
        // audioRecorder.getBuffers( gotBuffers );
    } else {
        // start recording
        if (!recorder)
            return;
        e.classList.add("recording");
        recorder.clear()
        startRecording()
        // audioRecorder.clear();
        // audioRecorder.record();
    }
}
function setupDownload(blob, filename){
    var url = (window.URL || window.webkitURL).createObjectURL(blob);
    var link = document.getElementById("save");
    link.href = url;
    link.download = filename || 'output.wav';
}

function startRecording() {
    recorder && recorder.record();
    console.log("Start recording")
}

function stopRecording() {
    recorder && recorder.stop()
    recorder.exportWAV(function(blob) {
        setupDownload(blob, "record.wav");
    })
    console.log("Stop recording")
}

function startStream(stream) {
    inputPoint = audioContext.createGain();
    // Create an AudioNode from the stream.
    realAudioInput = audioContext.createMediaStreamSource(stream);
    audioInput = realAudioInput;
    audioInput.connect(inputPoint);

    recorder = new Recorder(inputPoint)
    
    zeroGain = audioContext.createGain();
    zeroGain.gain.value = 0.0;
    inputPoint.connect( zeroGain );
    zeroGain.connect( audioContext.destination );

    /*
    var input = audioContext.createMediaStreamSource(stream)
    input.connect(audioContext.destination)
    recorder = new Recorder(input)
    */
}

window.onload = function init() {
    if (!navigator.getUserMedia)
            navigator.getUserMedia = navigator.webkitGetUserMedia || navigator.mozGetUserMedia;

    navigator.getUserMedia({audio: true}, startStream, function(e) {
        console.log('No live audio input: ' + e);
    });
};