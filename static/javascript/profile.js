
const WEBSOCKET_STATUS = {
  init: {
    requestFile: "init:request-file-path",
    register: "init:register"
  },
  start: {
    register: "register"
  }

}
const WEBSOCKET_RESOURCE = {
  profile: "profile"
}
const WEBSOCKET_RESPONSE_STATUS = {
  success: {
    filePath: "success:file-path"
  }
}


let url = "ws://" + window.location.host + "/api/v1/websocket";
let ws = new WebSocket(url);
var now = function() {
  var iso = new Date().toISOString();
  return iso.split("T")[1].split(".")[0];
};
ws.onmessage = function(msg) { //#region 
  var line = now() + " " + msg.data + "\n";
  try {
    const server_answer = JSON.parse(msg.data)
    switch (server_answer.resource) {
      case WEBSOCKET_RESOURCE.profile:

        ///////////////////////////////////
        switch (server_answer.status) {
          case WEBSOCKET_RESPONSE_STATUS.success.filePath:
            updatePathView(server_answer.data.value)
            switchButton(document.getElementById("button-request-path"));
            break;

          default:
            console.log("The server responded with an unknown status for the resource Profile")
            break;
        }


        break;

      default:
        console.log("The server responded with an unknown resource...")
        break;
    }
  } catch (err) {
    console.error(err)
  }

  console.log(line)
  // chat.innerText += line;
};
//#endregion
function websocketRefresh() {
  if (ws.readyState === ws.CLOSED || ws.readyState === ws.CLOSING) {
    ws = new WebSocket(url);
  }
}
function updatePathView(path) {
  let pathInput = document.getElementById("profilePath");
  pathInput.value = path;
}

////
function requestFilePath() {
  let button = document.getElementById("button-request-path")
  switchButton(button);
  ws.send(
    JSON.stringify({
      resource: "profile",
      status: WEBSOCKET_STATUS.init.requestFile,
      data: {}
    })
  )
}
function displayError(err) {
  //TODO
  /////////
  console.error("An error occured !")
  console.error(err)
}
function switchButton(button) {
  button.disabled = !button.disabled
}
////
function onSubmit(evt) {
  websocketRefresh()
  let list = document.getElementById("notification-list")
  evt.preventDefault()
  let form = evt.target;

  let formData = new FormData(form);
  const data = {
    profileName: formData.get("profileName"),
    gamePath: formData.get("profilePath")
  }
  const message = {
    resource: "profile", // profile/else
    status: WEBSOCKET_STATUS.init.register,
    data
  }
  console.log(message)
  ws.send(JSON.stringify(message))
  // let url = new URLSearchParams();
  // url.set()
  // fetch("/api/v1/profile", {
  //   method: "POST",

  // }).then(response => {
  //   if (response.ok) {

  //     console.log(response.json())
  //   } else {
  //     console.log(response)
  //     let err = document.createElement("div")
  //     err.className = "w-full rounded"
  //     // err.appendChild("hello")
  //     err.style.backgroundColor = "red"
  //     err.textContent = "Hello"
  //     list.appendChild(err)
  //     //document.removeChild(err)
  //   }
  // }).catch((err) => {
  //   console.log(err)
  //   console.error(err)
  // })
  // let fname = formData.get('profileName');
  // console.log(fname)
  // let lname = formData.get('lastName');
  // let gender = formData.get('gender');
  // if (fname != "" || lname != "" || gender != "") {
  //   $.ajax({
  //     url: '/api/profile',
  //     method: 'post',
  //     data: formData,
  //     processData: false,
  //     contentType: false,
  //     success: (d) => {
  //       console.log("Player Added", d);
  //       // profiles.innerHTML += d;
  //       // form.reset();
  //       location.reload()
  //     },
  //     error: (d) => {
  //       console.error(d);
  //       console.log("An error occurred. Please try again");
  //     }
  //   });
  // }

  return false;
}

// function disableButton() {

//   let pathInput = document.getElementById("profilePath");
// }
