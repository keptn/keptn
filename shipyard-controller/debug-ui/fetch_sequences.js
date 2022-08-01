function fetchSequences() {
  const project = document.getElementById("projectSelector").value;

  fetch("/sequence/project/" + project, {
    method: "get",
  })
    .then((res) => {
      return res.json();
    })
    .then((response) => {
      let table = document.getElementById("prodTable");

      const show_finished = document.getElementById("checkboxFinished").checked;
      const show_triggered =
        document.getElementById("checkboxTriggered").checked;
      const show_timedOut = document.getElementById("checkboxTimedOut").checked;

      table.innerHTML =
        "<tr><th>shkeptncontext</th><th>Service</th><th>Project</th><th>Stage</th><th>Service</th><th>view events</ht><th>getblocking</th></tr>";

      response.sequenceExecutions.forEach((object) => {
        if (
          (object.status.state == "triggered" && show_triggered) ||
          (object.status.state == "finished" && show_finished) ||
          (object.status.state == "timedOut" && show_timedOut)
        ) {
          let tr = document.createElement("tr");
          let td_blocking = document.createElement("td");

          if (object.status.state == "triggered") {
            td_blocking.innerHTML =
              "<button onclick=\"getBlocking('" +
              object.scope.keptnContext +
              "', '" +
              object.scope.project +
              "')\">" +
              "get blocking sequences" +
              "</button>";
          }

          tr.className = object.status.state;
          tr.innerHTML =
            "<td>" +
            object.scope.keptnContext +
            "</td>" +
            "<td>" +
            object.scope.service +
            "</td>" +
            "<td>" +
            object.scope.project +
            "</td>" +
            "<td>" +
            object.scope.stage +
            "</td>" +
            "<td>" +
            object.scope.service +
            "</td>" +
            '<td><a href="viewevents.html?shkeptncontext=' +
            object.scope.keptnContext +
            "&projectname=" +
            project +
            '"><button>View Events</button></a></td>';
          tr.appendChild(td_blocking);
          table.appendChild(tr);
        }
      });
    });
}

function getBlocking(context, project) {
  fetch(
    "/sequence/project/" + project + "/shkeptncontext/" + context + "/blocking",
    {
      method: "get",
    }
  )
    .then((res) => {
      return res.json();
    })
    .then((response) => {
      out = "";
      response.forEach((object) => {
        out += object.scope.keptnContext + "\n";
      });
      alert(out);
    });
}
