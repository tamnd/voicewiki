function welcome() {
	/*
	speakParagraph("Welcome to Voice Wikipedia. \
		Use your voice to tell which Wikipedia entry to search for.\
		No need to use your eyes, the entry's content will be reading to you!.\
		Let's try it. Say. Show me.")
*/
	speakParagraph("Welcome to Voice Wikipedia.")
}

function list() {
	speakParagraph("There are 10 recent wikipedia articles.\
		Do you want to hear them?.")
}

function help() {
	speakParagraph("Welcome to Voice Wikipedia. \
		You can using voice commands to control this application.\
		If you want to search for a wikipedia entry. Say. Show me.\
		If you want to hear recent wikpedia. Say. List.")
}

function search(term) {
	para = "You want to search for " + term + "." +
		" Found an article. Do you want to hear this?.";
	writelog(para);
	speakParagraph(para);
}

function writelog(str) {
	console.log(str);
	node = document.getElementById("log");
	node.innerHTML = str;
	// $(".log").text(str);
}

function initAnnyang() {
	console.log("INIT");
	console.log(annyang);
	if (annyang) {
	
	  // Let's define our first command. First the text we expect, and then the function it should call
	  var commands = {
	  	'show me *entry': search,
	  	'list': list,
	  	'help': help,
	  	'help me': help

	  	// 'help': help
	  };
	  // Add our commands to annyang
	  annyang.addCommands(commands);

	  // Start listening. You can call this here, or attach this call to an event, button, etc.
	  annyang.start();
	}
}

function initGestures() {
	var gestures = new Array();
	// gestures["L"] = "46";
	gestures["L"] = "64";
	gestures["S"] = "432101234";
	gestures["?"] = "6701232";
	gestures["R"] = "267012341";
	gestures["Z"] = "030";

	$('body').gestures(gestures, function (data) {
		// document.getElementById('outputbox').innerHTML += data;
		// document.getElementById('outputbox').innerHTML = data;
		if (data !== "") {
			writelog(data);
			// speakEnglish(data);	
			switch(data) {
			    case "S":
			        search("Taylor Swift");
			        break;
			    case "L":
			        list();
			        break;
			    case "?": 
			    	help();
			    	break;
			    default:
			        
			}
		}
		
	});
}

function init() {
	initAnnyang();
	initGestures();
}


welcome();
init();