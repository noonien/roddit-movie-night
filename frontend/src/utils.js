
var apiKey = '5bf19fea67cd240b00f1b7ea'; // <-- Replace with your app_id from https://www.opengraph.io/

export function getOGData(url) {
  var urlEncoded = encodeURIComponent(url)
  var requestUrl = "https://opengraph.io/api/1.1/site/" + urlEncoded + '?app_id=' + apiKey

  return $.getJSON(requestUrl)
}

export function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}