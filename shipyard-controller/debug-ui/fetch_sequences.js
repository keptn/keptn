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
      const show_aborted = document.getElementById("checkboxAborted").checked;
      const show_active = document.getElementById("checkboxActive").checked;
      const show_blocked = document.getElementById("checkboxBlocked").checked;

      table.innerHTML =
        "<tr><th>shkeptncontext</th><th>SequenceName</th><th>Projectname</th><th>service</th></tr>";

      response.states.forEach((object) => {
        if (
          (object.state == "finished" && show_finished) ||
          (object.state == "active" && show_active) ||
          (object.state == "blocked" && show_blocked) ||
          (object.state == "aborted" && show_aborted)
        ) {
          let tr = document.createElement("tr");
          tr.className = object.state;
          tr.innerHTML =
            "<td>" +
            object.shkeptncontext +
            "</td>" +
            "<td>" +
            object.name +
            "</td>" +
            "<td>" +
            object.project +
            "</td>" +
            "<td>" +
            object.service +
            "</td>" +
            "<td><a href=\"viewevents.html?shkeptncontext=" +
            object.shkeptncontext +
            "&projectname=" +
            project +
            "\"><button>View Events</button></a></td>"
          table.appendChild(tr);
        }
      });
    });
}
