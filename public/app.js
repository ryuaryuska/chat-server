window.addEventListener('DOMContentLoaded', (_) => {
  // let username = "ryu";
  var roomname = getParameterByName("room");
  var username = getParameterByName("name");
  let websocket = new WebSocket("ws://192.168.1.16:8080/websocket?name="+username);
  let room = document.getElementById("chat-text");
  let image = document.getElementById("image");

  websocket.addEventListener("open", (e) => {
    // websocket.send("test")
    websocket.send(JSON.stringify({message: roomname, action: "join-room", sender:{ name: "ryu"}}))
  });
  websocket.addEventListener("message", function (e) {
    let data = JSON.parse(e.data);
    console.log("request client: " + e.data)
    let chatContent = `<p>${data.sender.name}: ${data.message}  ${data.time}</p>`;
    room.innerHTML += chatContent;
    room.scrollTop = room.scrollHeight; // Auto scroll to the bottom

    if(data.file !== "") {
      let imageContent  = `<img width="700" height="500" src="./images/${data.file}"/>`
      image.innerHTML += imageContent;
      image.scrollTop = image.scrollHeight;
    }

  });

  let form = document.getElementById("input-form");
  form.addEventListener("submit", function (event) {
    event.preventDefault();
    let text = document.getElementById("input-text");
    let file = document.getElementById("file");
    
    let filename = ""

    if (file.files.item(0) != null) {

        filename = file.files.item(0).name

        let formData = new FormData();
        formData.append('myFile', file.files[0],  file.files.item(0).name);

        fetch("http://localhost:8080/upload",
          {
            body: formData,
            method: "post",
            redirect: 'follow'
          }).then((response) => response.json())
            .then((result) => console.log("result: " + result))
            .catch(error => console.log('error', error));
      }
    check = text.value


    console.log("filename: " + filename)

    if(check.substring(0, 1) === "/"){
      websocket.send(
        JSON.stringify({
          action: 'bot-message',
          message: text.value,
          target: roomname,
          sender: {
            name: "Topin"
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
          },
          file: filename
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