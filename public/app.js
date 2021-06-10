window.addEventListener('DOMContentLoaded', (_) => {
  // let username = "ryu";
  var roomname = getParameterByName("room");
  var username = getParameterByName("name");
  let websocket = new WebSocket("ws://192.168.1.16:8080/websocket?name="+username);
  let room = document.getElementById("chat-text");

  websocket.addEventListener("open", (e) => {
    websocket.send(JSON.stringify({message: roomname, action: "join-room", sender:{ name: username}}))
  });
  websocket.addEventListener("message", function (e) {
    let data = JSON.parse(e.data);
    console.log(data);
    let chatContent = `<p>${data.sender.name}: ${data.message}</p>`;
    room.innerHTML += chatContent
    room.scrollTop = room.scrollHeight; // Auto scroll to the bottom
  });

  let form = document.getElementById("input-form");
  form.addEventListener("submit", function (event) {
    event.preventDefault();
    let text = document.getElementById("input-text");
    check = text.value

    if(check.substring(0, 1) === "/"){
      websocket.send(
        JSON.stringify({
          action: 'bot-message',
          message: text.value,
          target: roomname,
          sender: {
            name: "BOT"
          }
        }));
    }else{
      websocket.send(
        JSON.stringify({
          action: 'send-message',
          message: text.value,
          target: roomname,
          sender: {
            name: username
          }
        }));
    }
    text.value = "";
  });
});


function getParameterByName(name, url = window.location.href) {
  name = name.replace(/[\[\]]/g, '\\$&');
  var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
      results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}