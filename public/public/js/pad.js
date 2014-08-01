var ctx;

$.fn.pad = function () {
	// $("#content").height($(window).height()-90);
    // var canvas = '<canvas id="canvas" width="'+$("#content").width()+'" height="'+$("#content").height()+'"></canvas>';
    var canvas = '<canvas id="canvas" width="'+$(window).width()+'" height="'+$(window).height()+'" style="position:fixed;top:0;left:0;right:0;bottom:80px;z-index:10;"></canvas>';
	// $(this).html(canvas);
	$(this).append(canvas);

	// setup canvas
	ctx=document.getElementById("canvas").getContext("2d");
	ctx.strokeStyle = "#fff";
	// ctx.lineWidth = 5;	
	ctx.lineWidth = 20;

	// setup to trigger drawing on mouse or touch
	$("#canvas").drawTouch();
    $("#canvas").drawPointer();
	$("#canvas").drawMouse();
}

// prototype to	start drawing on touch using canvas moveTo and lineTo
$.fn.drawTouch = function() {
	var clear = function(e) {
		canvas = $("#canvas");
		// $("h1").html("TOUCH END");
		console.log("TOUCH clear canvas: ", canvas.width(), canvas.height());
		ctx.clearRect(0, 0, canvas.width(), canvas.height());
	}

	var start = function(e) {
		clear(e);

        e = e.originalEvent;
		ctx.beginPath();
		x = e.changedTouches[0].pageX;
		y = e.changedTouches[0].pageY-44;
		ctx.moveTo(x,y);


	};
	var move = function(e) {
		e.preventDefault();
        e = e.originalEvent;
		x = e.changedTouches[0].pageX;
		y = e.changedTouches[0].pageY-44;
		ctx.lineTo(x,y);
		ctx.stroke();
	};

	
	$(this).on("touchstart", start);
	$(this).on("touchmove", move);
	// $(this).on("touchend", clear);	
}; 
    
// prototype to	start drawing on pointer(microsoft ie) using canvas moveTo and lineTo
$.fn.drawPointer = function() {
	var start = function(e) {
        e = e.originalEvent;
		ctx.beginPath();
		x = e.pageX;
		y = e.pageY-44;
		ctx.moveTo(x,y);
	};
	var move = function(e) {
		e.preventDefault();
        e = e.originalEvent;
		x = e.pageX;
		y = e.pageY-44;
		ctx.lineTo(x,y);
		ctx.stroke();
    };
	$(this).on("MSPointerDown", start);
	$(this).on("MSPointerMove", move);
};        

// prototype to	start drawing on mouse using canvas moveTo and lineTo
$.fn.drawMouse = function() {
	var clicked = 0;

	var clear = function(e) {
		canvas = $("#canvas");
		console.log("MOUSE clear canvas: ", canvas.width(), canvas.height());
		ctx.clearRect(0, 0, canvas.width(), canvas.height());
	}

	var start = function(e) {
		clear();
		console.log(e);
		clicked = 1;
		ctx.beginPath();
		x = e.offsetX;
		// y = e.pageY-44;
		y = e.offsetY;
		ctx.moveTo(x,y);
	};
	var move = function(e) {
		if(clicked){
			x = e.offsetX;
			// y = e.pageY-44;
			y = e.offsetY;
			ctx.lineTo(x,y);
			ctx.stroke();
		}
	};
	var stop = function(e) {
		clicked = 0;
		// clear canvas
		
	};
	$(this).on("mousedown", start);
	$(this).on("mousemove", move);
	$(window).on("mouseup", stop);
};
