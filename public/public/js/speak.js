
function speakEnglish(text, cb) {
	// window.speechSynthesis.cancel();
	console.log("Start speaking English: ", text);
	var msg = new SpeechSynthesisUtterance();
	var voices = window.speechSynthesis.getVoices();
	var voiceEn = voices[0];
	msg.voice = voiceEn; // Note: some voices don't support altering params
	msg.text = text;
	// msg.lang = 'ja-JP';
	// msg.volume = 1; // 0 to 1
	// msg.rate = 10; // 0.1 to 10
	// msg.pitch = 2; //0 to 2
	// console.log(msg);
	if (typeof cb == "function") {
		msg.onend = cb
	}
	speechSynthesis.speak(msg);
}

function speakParagraph(paragraph, cb) {
	//sentences = paragraph.replace(/([.?!])\s*(?=[A-Z])/, "$1|").split("|")
	sentences = paragraph.match( /[^\.!\?]+[\.!\?]+/g );
	//*
	console.log(sentences);
	l = sentences.length
	for (i = 0; i < l-1; i++) {
		speakEnglish(sentences[i]);
	}
	speakEnglish(sentences[i], cb);
	//*/
	// _speakSentence(sentences, 0, cb);	
}

function _speakSentence(sentences, index, cb) {
	console.log(index + " - " + sentences[index]);
	if (index < sentences.length - 1) {
		speakEnglish(sentences[index], function (event) {
			_speakSentence(sentences, index + 1, cb);
		});
	} else {
		speakEnglish(sentences[index], function (event) {
			if (typeof cb == "function") {
				cb(event);
			}
		})
	}
}

function speak(nodeID, cb) {
	element = document.getElementById(nodeID);
	var text = element.textContent || element.innerText;
	if (typeof cb == "function") {
		speakEnglish(text, cb);
	} else {
		speakEnglish(text);
	}
	
	/*
	// sentences = text.replace(/([.?!])\s*(?=[A-Z])/, "$1|").split("|")
	sentences = text.match( /[^\.!\?]+[\.!\?]+/g );
	console.log(sentences)
	l = sentences.length
	for (i = 0; i < l; i++) {
		speakEnglish(sentences[i])
	}	
	*/
}

function cancelSpeak() {
	window.speechSynthesis.cancel();
}



function speakJa(nodeID) {
	element = document.getElementById(nodeID);
	var text = element.textContent || element.innerText;
	// sentences = text.replace(/([.?!])\s*(?=[A-Z])/, "$1|").split("|")
	sentences = text.match( /[^\。!\?]+[\。!\?]+/g );
	console.log(sentences)
	l = sentences.length
	for (i = 0; i < l; i++) {
		speakJapanese(sentences[i])
	}	
}