var images=new Array();
function preload(){
    for(i = 0; i<preload.arguments.length;i++){
    images[i]=new Image();images[i].src=preload.arguments[i]
}
}          
    preload('/static/img/photo-1.jpg', '/static/img/photo-2.jpg', '/static/img/photo-3.jpg', '/static/img/photo-4.jpg')

function changeBg(){
    var bgs=['/static/img/photo-1.jpg', '/static/img/photo-2.jpg', '/static/img/photo-3.jpg', '/static/img/photo-4.jpg'];
    setInterval(function(){
        var p=bgs.shift();
        document.getElementsByClassName('main-bar')[0].style.backgroundImage='url(' + p + ')';
        bgs.push(p)},6000)
    }
changeBg();