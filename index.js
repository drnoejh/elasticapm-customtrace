function displayTraceID() {
    var req = new XMLHttpRequest();
    req.open('GET', document.location, false);
    req.send(null);
    var headers = parseHttpHeaders(req.getAllResponseHeaders())
    document.getElementsByTagName("p")[0].innerHTML = headers["traceid"]
}

function parseHttpHeaders(httpHeaders) {
    return httpHeaders.split("\n")
     .map(x=>x.split(/: */,2))
     .filter(x=>x[0])
     .reduce((ac, x)=>{ac[x[0]] = x[1];return ac;}, {});
}