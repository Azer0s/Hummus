var console = Java.type('com.sim.hummus.Console');
rep = function (string, times) {var repeatedString = "";while (times > 0) {repeatedString += string;times--;}return repeatedString;}; //repeat a string
len = function (string) { return string.length;}; //string length
add = function (a,b) { return a + b; };
sub = function (a,b) { return a - b; };
rev = function (str) { var newString = "";for (var i = str.length - 1; i >= 0; i--) {newString += str[i];}return newString;};
put = function (a) { console.write(a); };
get = function () { return console.read(); };