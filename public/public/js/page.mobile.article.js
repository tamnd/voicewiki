function readTitle(title) {
	speakEnglish("This is article " + title + " in Wikipedia.");
}

readTitle("Bicycle");

var autoNext = true;

var curid = 0; // current paragraph.

// nextId returns the next id of the reading paragraph.
function getId(node) {
	id = $(node).attr("id").substring(2);
	return parseInt(id);
}

function speakParaWithId(id) {
	speakPara($("#p-" + id.toString()));
}

var scrollToElement = function(node, ms){
    var speed = (ms) ? ms : 600;
    $('html,body').animate({
        scrollTop: node.offset().top - 100
    }, speed);
}


// Speak paragraph #id.
function speakPara(node) {	
	if (typeof node == "undefined") {
		return
	}
	scrollToElement(node, 600);
	$("p").css("color", "#444");
	curid = getId(node);
	play();

	console.log("Paragraph: " + curid);
	node.css("color", "red");
	// console.log($me.text());
	cancelSpeak();
	speakParagraph(node.text(), function (event) {
		pause();
		node.css("color", "#444");
		console.log("DONE!");
		console.log(autoNext);
		if (autoNext) {
			curid++;
			speakPara($("#p-" + curid.toString()));
		}
	});
}

var elPlay = $("#icon-play");

elPlay.click(function (event) {
	if (elPlay.hasClass("fa-play")) {
		play();
		speakParaWithId(curid);
	} else {
		doStop();
	}
});

function play() {
	elPlay.removeClass("fa-play");
	elPlay.addClass("fa-pause");
}

function pause() {
	elPlay.removeClass("fa-pause");
	elPlay.addClass("fa-play");
}

// ACTIONS
function doNext() {
	curid ++;
	speakParaWithId(curid);
}

function doPrev() {
	if (curid > 0) {
		curid--;
		speakParaWithId(curid);
	}
}

function doStop() {
	console.log("STOP");
	pause();
	cancelSpeak();
}

function doStart() {
	play();
	speakParaWithId(curid);
}

function doSearch(term) {
	para = "You want to search for " + term + "." +
		" Found an article. Do you want to hear this?.";
	// writelog(para);
	speakParagraph(para);
}

function initAnnyang() {
	if (annyang) {
	
	  // Let's define our first command. First the text we expect, and then the function it should call
	  var commands = {
	  	'show me *entry': doSearch,
	  	// 'list': search,
	  	//'help': help,
	  	// 'help me': help
	  	"next": doNext,
	  	"back": doPrev,
	  	"stop": doStop,
	  	"start": doStart,
	  };
	  // Add our commands to annyang
	  annyang.addCommands(commands);

	  // Start listening. You can call this here, or attach this call to an event, button, etc.
	  annyang.start();
	}
}

function goHome() {
	window.location = "/";
}

function no() {
	console.log("NO");
}

function yes() {
	console.log("YES");
}

function initGestures() {
	console.log("Init gestures");
	var gestures = new Array();
	// gestures["L"] = "46";
	
	gestures["H"] = "267012";
	// gestures["?"] = "6701232";
	gestures["P"] = "670123456";

	gestures["N"] = "616";
	gestures["Y"] = "21076234567";
	gestures["S"] = "432101234";

	$('body').gestures(gestures, function (data) {
		// document.getElementById('outputbox').innerHTML += data;
		// document.getElementById('outputbox').innerHTML = data;
		console.log("GESTURE");
		if (data !== "") {
			console.log(data);
			// speakEnglish(data);	
			switch(data) {			    
			    case "H":
			        goHome();
			        break;

			     case "S":
			        doStart();
			        break;
			    
			     case "P":
			        doStop();
			        break;

			    case "N":
			        no();
			        break;
			    case "Y":
			        yes();
			        break;
			   
			    default:			        
			}
		}		
	});
}

function initPad() {
	// setup a new canvas for drawing wait for device init
   	//  setTimeout(function(){
	// newCanvas();
	$(".content").pad("#8e44ad");
    // }, 1000);	
}

function init() {
	initAnnyang();
	initGestures();
	initPad();
}


$("p").click(function (event) {
	speakPara($(this));
});

init();

// speakParaWithId(curid);
