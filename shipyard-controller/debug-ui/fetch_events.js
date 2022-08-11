function fetchEvents(context, projectname, rootevent_list) {
  let list = document.getElementById("tree_container");
  list.innerHTML = "";

  fetch(`/sequence/project/${projectname}/shkeptncontext/${context}/event`,
    {
      method: "get",
    }
  )
    .then((res) => {
      return res.json();
    })
    .then((response) => {
      let rootevent_li, taskevent_ul, label;

      response.forEach((object) => {
        // rootevent has additional stage segment
        if (object.type.split(".").length === 6) {
          rootevent_li = document.createElement("li");

          label = document.createElement("span");
          label.innerHTML = object.type;

          rootevent_li.appendChild(label);

          taskevent_ul = document.createElement("ul");
          taskevent_ul.className = "nested";
        } else {
          label.className = "caret";

          let taskevent_li = document.createElement("li");

          let label2 = document.createElement("span");
          label2.className = "caret";
          label2.innerHTML = object.type;
          taskevent_li.appendChild(label2);

          let detail_ul = document.createElement("ul");
          detail_ul.className = "nested";

          for (let key in object.data) {
            let li = document.createElement("li");
            li.innerHTML = `${key}:  ${object.data[key]}`;
            detail_ul.appendChild(li);
          }

          taskevent_li.appendChild(detail_ul);
          taskevent_ul.appendChild(taskevent_li);
        }

        if (taskevent_ul !== null) {
          rootevent_li.appendChild(taskevent_ul);
          rootevent_list.appendChild(rootevent_li);
        }
      });

      let toggler = document.getElementsByClassName("caret");

      for (let toggle_item of toggler) {
        toggle_item.addEventListener("click", function () {
          this.parentElement
            .querySelector(".nested")
            .classList.toggle("active");
          this.classList.toggle("caret-down");
        });
      }
    });
}
