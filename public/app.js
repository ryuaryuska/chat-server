window.addEventListener("DOMContentLoaded", (_) => {
  var roomname = getParameterByName("room");
  var username = getParameterByName("name");
  let websocket = new WebSocket(
    "ws://localhost:8080/websocket?name=" + username
  );
  let room = document.getElementById("sender-chat");

  websocket.addEventListener("open", (e) => {
    websocket.send(JSON.stringify({ message: roomname, action: "join-room" }));

    let form = document.getElementById("input-form");
    form.addEventListener("submit", function (event) {
      event.preventDefault();
      let text = document.getElementById("input-text");
      check = text.value;

      websocket.send(
        JSON.stringify({
          action: "send-message",
          message: text.value,
          target: roomname,
          sender: {
            name: username,
          },
          // file: filename
        })
      );
      text.value = "";
    });
  });

  websocket.addEventListener("message", function (e) {
    let data = JSON.parse(e.data);

    let clasz = "talk-bubble talk-bubble-recipient ";
    if (data.sender.name != "student") {
      clasz = "talk-bubble talk-bubble-recipient ";
    }

    let chatContent = `<div class='d-flex justify-content-start' id='sender-chat'>
    <div class='${clasz} tri-right round right-in'>
        <div class='talktext text-white'>
            <p class='mb-2'>
                ${data.message}
            </p>
            <p class='text-end mx-1 fw-bold'>
                ${data.time}
            </p>
        </div>
    </div>
</div>`;
    room.innerHTML += chatContent;
    room.scrollTop = room.scrollHeight; // Auto scroll to the bottom
  });
});

function getParameterByName(name, url = window.location.href) {
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
    results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return "";
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}
